import { defineStore, storeToRefs } from 'pinia'
import { ref, computed } from 'vue'

import type { PeerConnection } from '@/types'
import { useRoomWsStore } from './wsConnection'

//TODO (4) Сделать нормальные cleanup функции

//TODO (4) Добавить захват микрофона/вебки

//BUGFIX (5) После реконнекта через перезагрузку страницы не отправляются/не получаются треки

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
  const isMicrophoneOn = ref<boolean>(false)
  const microphoneStream = ref<MediaStream | null>(null)

  const rtcConnectionQueue = ref<any>([])

  async function getMicrophoneTrack() {
    microphoneStream.value = await navigator.mediaDevices.getUserMedia({
      audio: {
        sampleRate: 48000,
        noiseSuppression: true,
        echoCancellation: true,
      },
      peerIdentity: localUuid.value,
    })
    microphoneStream.value.getAudioTracks()[0].contentHint = 'speech'
  }

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
      roomWs.value!.send(
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

  async function createdDescription(description: RTCSessionDescriptionInit, peerUuid: string) {
    peerConnections.value[peerUuid].pc
      .setLocalDescription(description)
      .then(() => {
        roomWs.value!.send(
          JSON.stringify({
            sdp: peerConnections.value[peerUuid].pc.localDescription,
            uuid: localUuid.value,
            dest: peerUuid,
          }),
        )
      })
      .catch((e: Error) => {
        console.log(e)
      })
  }

  function toggleMicrophone() {
    isMicrophoneOn.value ? (isMicrophoneOn.value = false) : (isMicrophoneOn.value = true)
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

  function shareMicrophone(peerUuid: string) {
    roomWs.value!.send(
      JSON.stringify({
        uuid: localUuid.value,
        metadata: 'microphone',
        dest: peerUuid,
      }),
    )
    if (microphoneStream.value) {
      const senders = peerConnections.value[peerUuid].pc.getSenders()
      const microhoneSender = senders.find((sender) => sender.track?.contentHint == 'speech')
      if (microhoneSender) {
        microhoneSender.replaceTrack(microphoneStream.value.getAudioTracks()[0])
      } else {
        microphoneStream.value?.getTracks().forEach((track) => {
          peerConnections.value[peerUuid].microphoneSender = peerConnections.value[
            peerUuid
          ].pc.addTrack(track, microphoneStream.value as MediaStream)
        })
      }
    }
  }

  function shareScreen(peerUuid: string) {
    if (localStream.value) {
      roomWs.value!.send(
        JSON.stringify({
          uuid: localUuid.value,
          metadata: 'screen-share',
          dest: peerUuid,
        }),
      )
      localStream.value.getTracks().forEach((track) => {
        if (track.kind === 'video') {
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

  function setTrackMetadata(peerUuid: string, metadata: 'microphone' | 'screen-share' | 'webcam') {
    peerConnections.value[peerUuid].trackMetadata = metadata
  }

  function setupPeer(peerUuid: string, displayName: string, initCall: boolean) {
    const peerConnection: PeerConnection = {
      displayName,
      uuid: peerUuid,
      ssVideoStream: null,
      microphoneStream: null,
      pc: new RTCPeerConnection(PEER_CONNECTION_CFG),
    }

    peerConnection.pc.ontrack = (e) => {
      if (!e.streams || !(e.streams.length > 0)) {
        return
      }
      console.log(peerConnections.value[peerUuid].trackMetadata)
      if (peerConnections.value[peerUuid].trackMetadata == 'microphone') {
        peerConnections.value[peerUuid].microphoneStream = e.streams[0]
      } else if (peerConnections.value[peerUuid].trackMetadata == 'screen-share') {
        peerConnections.value[peerUuid].ssVideoStream = e.streams[0]
      }
    }

    peerConnection.pc.onnegotiationneeded = async (e) => {
      const leaderId = [localUuid.value, peerUuid].sort()[0]

      const offer = await peerConnections.value[peerUuid].pc.createOffer({
        offerToReceiveAudio: true,
        offerToReceiveVideo: true,
      })
      if (leaderId == localUuid.value) {
        await createdDescription(offer, peerUuid)
      }
    }

    peerConnection.pc.onconnectionstatechange = (e) => {
      checkPeerDisconnect(peerUuid)
      const state = peerConnections.value[peerUuid]?.pc.connectionState
      if (state === 'connected') {
        roomWs.value!.send(
          JSON.stringify({
            uuid: localUuid.value,
            type: 'get-microphone',
          }),
        )

        /*         peerConnections.value[peerUuid].microphoneStream = new MediaStream([
          peerConnections.value[peerUuid].pc.getTransceivers()[0].receiver.track,
        ])
 */
        /*  peerConnections.value[peerUuid].microphoneStream =  */

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
        .then(async (description) => await createdDescription(description, peerUuid))
        .catch((e) => console.log(e))
    }

    peerConnections.value[peerUuid] = peerConnection
  }

  async function handleSdpSignal(signal: any, peerUuid: string) {
    peerConnections.value[peerUuid].pc
      .setRemoteDescription(new RTCSessionDescription(signal.sdp))
      .then(() => {
        if (signal.sdp.type === 'offer') {
          peerConnections.value[peerUuid].pc
            .createAnswer()
            .then(async (description) => await createdDescription(description, peerUuid))
            .catch((e) => console.log(e))
        } else {
          console.log(signal)
        }
      })
  }

  function handleIceCandidate(signal: any, peerUuid: string): void {
    console.log(peerConnections.value[peerUuid].pc.signalingState)
    peerConnections.value[peerUuid].pc
      .addIceCandidate(new RTCIceCandidate(signal.ice))
      .catch((e) => console.log(e))
  }

  return {
    peerConnections,
    localStream,
    localUuid,
    localDisplayName,
    isMicrophoneOn,
    microphoneStream,
    rtcConnectionQueue,
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
    toggleMicrophone,
    getMicrophoneTrack,
    shareMicrophone,
    setTrackMetadata,
  }
})
