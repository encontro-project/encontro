import { defineStore, storeToRefs } from 'pinia'
import { ref, computed } from 'vue'

import type { PeerConnection } from '@/types'
import { useRoomWsStore } from './wsConnection'

//TODO (4) Сделать нормальные cleanup функции

//TODO (4) Добавить захват микрофона/вебки

const VIDEO_STREAM_WIDTH: number = 1920
const VIDEO_STREAM_HEIGHT: number = 1080
const VIDEO_STREAM_FRAME_RATE: number = 60

const PEER_CONNECTION_CFG: RTCConfiguration = {
  iceServers: [{ urls: 'stun:stun.l.google.com:19302' }, { urls: 'stun:stun1.l.google.com:19302' }],
}

const CONSTRAINTS: MediaStreamConstraints = {
  audio: true,
  video: {
    width: { ideal: VIDEO_STREAM_WIDTH, max: VIDEO_STREAM_WIDTH },
    height: { ideal: VIDEO_STREAM_HEIGHT, max: VIDEO_STREAM_HEIGHT },
    frameRate: { ideal: VIDEO_STREAM_FRAME_RATE, max: VIDEO_STREAM_FRAME_RATE },
  },
}

export const useConnectionsStore = defineStore('connections', () => {
  const roomWsStore = useRoomWsStore()

  const { roomWs, localUuid, localDisplayName } = storeToRefs(roomWsStore)

  const peerConnections = ref<Record<string, PeerConnection>>({})
  const localStream = ref<MediaStream | null>(null)
  const currentRoomUrl = ref<string>('')

  async function getMediaTracks() {
    const devices = await navigator.mediaDevices.enumerateDevices()
    console.log(
      'Audio devices:',
      devices.filter((d) => d.kind === 'audioinput'),
    )
    console.log(
      'Video devices:',
      devices.filter((d) => d.kind === 'videoinput'),
    )

    if (navigator.mediaDevices.getDisplayMedia) {
      localStream.value = await navigator.mediaDevices.getDisplayMedia(CONSTRAINTS)
    } else if (navigator.mediaDevices.getUserMedia) {
      localStream.value = await navigator.mediaDevices.getUserMedia(CONSTRAINTS)
    } else {
      alert('Your browser does not support getDisplayMedia or getUserMedia API')
      throw new Error('Unsupported media API')
    }
  }

  function handlePeerDisconnect(peerUuid: string) {
    console.log(`Peer ${peerUuid} disconnected`)
    peerConnections.value[peerUuid]?.pc.close()
    delete peerConnections.value[peerUuid]
  }

  function updateStream() {
    Object.values(peerConnections.value).forEach((connection) =>
      connection.pc.getSenders().forEach((sender) => {
        try {
          const videoTrack = localStream.value?.getVideoTracks()[0]
          const audioTrack = localStream.value?.getAudioTracks()[0]
          console.log(audioTrack)
          if (sender.track?.kind === 'video') {
            connection.ssVideoSender = sender
            sender.replaceTrack(videoTrack as MediaStreamTrack)
          }
          if (sender.track?.kind === 'audio') {
            connection.ssAudioSender = sender
            sender.replaceTrack(audioTrack as MediaStreamTrack)
          }
        } catch {}
      }),
    )
  }

  function gotIceCandidate(event: RTCPeerConnectionIceEvent, peerUuid: string) {
    if (event.candidate) {
      roomWs.value!!.send(
        JSON.stringify({
          ice: event.candidate,
          uuid: localUuid.value,
          dest: peerUuid,
        }),
      )
    }
  }

  function checkPeerDisconnect(peerUuid: string): void {
    const state = peerConnections.value[peerUuid]?.pc.connectionState
    if (['failed', 'closed', 'disconnected'].includes(state)) {
      handlePeerDisconnect(peerUuid)
    }
  }

  function createdDescription(description: RTCSessionDescriptionInit, peerUuid: string): void {
    peerConnections.value[peerUuid].pc
      .setLocalDescription(description)
      .then(() => {
        roomWs.value!!.send(
          JSON.stringify({
            sdp: peerConnections.value[peerUuid].pc.localDescription,
            uuid: localUuid.value,
            dest: peerUuid,
          }),
        )
      })
      .catch((e: Error) => console.log(e))
  }

  function stopStream() {
    localStream.value?.getTracks().forEach((track) => {
      track.stop()
    })
    localStream.value = null
  }

  function unsubscribeFromStream(peerUuid: string) {
    const peerConnection = peerConnections.value[peerUuid]
    peerConnection.pc.getReceivers().forEach((rec) => {
      rec.track.enabled = false
    })
    peerConnection.ssVideoStream = null
  }

  function leaveCall() {
    stopStream()
    Object.values(peerConnections.value).forEach((peerConnection) => {
      peerConnection.pc.close()
    })
    peerConnections.value = {}
  }

  function updateTrack() {
    Object.values(peerConnections.value).forEach((connection) =>
      connection.pc.getSenders().forEach((sender) => {
        try {
          const videoTrack = localStream.value?.getVideoTracks()[0]
          const audioTrack = localStream.value?.getAudioTracks()[0]
          console.log(audioTrack)
          if (sender.track?.kind === 'video') {
            sender.replaceTrack(videoTrack as MediaStreamTrack)
          }
          if (sender.track?.kind === 'audio') {
            sender.replaceTrack(audioTrack as MediaStreamTrack)
          }
        } catch {}
      }),
    )
  }

  function shareScreen(peerUuid: string) {
    if (localStream.value) {
      console.log(localStream.value.getTracks())
      localStream.value.getTracks().forEach((track) => {
        if (track.kind === 'video') {
          peerConnections.value[peerUuid].ssAudioTrack = track
          peerConnections.value[peerUuid].ssVideoSender = peerConnections.value[
            peerUuid
          ].pc.addTrack(track, localStream.value as MediaStream)
        }
        if (track.kind === 'audio') {
          peerConnections.value[peerUuid].ssAudioSender = peerConnections.value[
            peerUuid
          ].pc.addTrack(track, localStream.value as MediaStream)
        }
      })
    }
  }

  function setupPeer(peerUuid: string, displayName: string, initCall: boolean) {
    const peerConnection: PeerConnection = {
      displayName,
      uuid: peerUuid,
      ssVideoStream: null,
      pc: new RTCPeerConnection(PEER_CONNECTION_CFG),
    }

    peerConnection.pc.ontrack = (e) => {
      if (e.streams && e.streams.length > 0) {
        peerConnections.value[peerUuid].ssVideoStream = e.streams[0]
      }
    }

    peerConnection.pc.onnegotiationneeded = () => {
      peerConnection.pc
        .createOffer({ offerToReceiveAudio: true, offerToReceiveVideo: true })
        .then((description) => createdDescription(description, peerUuid))
        .catch((e) => console.log(e))
    }

    peerConnection.pc.onconnectionstatechange = (e) => {
      checkPeerDisconnect(peerUuid)
      const state = peerConnections.value[peerUuid]?.pc.connectionState
      if (state === 'connected') {
        roomWs.value!.send(
          JSON.stringify({
            uuid: localUuid.value,
            type: 'get-tracks',
          }),
        )
      }
    }

    peerConnection.pc.onicecandidate = (event) => gotIceCandidate(event, peerUuid)

    if (initCall) {
      peerConnection.pc
        .createOffer({ offerToReceiveAudio: true, offerToReceiveVideo: true })
        .then((description) => createdDescription(description, peerUuid))
        .catch((e) => console.log(e))
    }

    peerConnections.value[peerUuid] = peerConnection
  }

  function handleSdpSignal(signal: any, peerUuid: string): void {
    peerConnections.value[peerUuid].pc
      .setRemoteDescription(new RTCSessionDescription(signal.sdp))
      .then(() => {
        if (signal.sdp.type === 'offer') {
          peerConnections.value[peerUuid].pc
            .createAnswer()
            .then((description) => createdDescription(description, peerUuid))
            .catch((e) => console.log(e))
        }
      })
      .catch((e) => console.log(e))
  }

  function handleIceCandidate(signal: any, peerUuid: string): void {
    peerConnections.value[peerUuid].pc
      .addIceCandidate(new RTCIceCandidate(signal.ice))
      .catch((e) => console.log(e))
  }

  return {
    peerConnections,
    localStream,
    localUuid,
    localDisplayName,
    getMediaTracks,
    handlePeerDisconnect,
    updateStream,
    gotIceCandidate,
    checkPeerDisconnect,
    createdDescription,
    stopStream,
    unsubscribeFromStream,
    leaveCall,
    updateTrack,
    shareScreen,
    setupPeer,
    handleSdpSignal,
    handleIceCandidate,
  }
})
