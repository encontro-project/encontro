import { defineStore, storeToRefs } from 'pinia'
import { ref, computed } from 'vue'

import type { PeerConnection } from '@/types'
import { useRoomWsStore } from './wsConnection'

// TODO (4) Сделать нормальные cleanup функции

// TODO (6) Добавить захват вебки

// TODO (6) Добавить выбор устройств через менюшку

// TODO (7) Перенести getUserMedia логику в roomWsStore

const PEER_CONNECTION_CFG: RTCConfiguration = {
  iceServers: [{ urls: 'stun:stun.l.google.com:19302' }, { urls: 'stun:stun1.l.google.com:19302' }],
}

export const useConnectionsStore = defineStore('connections', () => {
  const roomWsStore = useRoomWsStore()

  const { roomWs, localUuid, localDisplayName, microphoneStream, localStream } =
    storeToRefs(roomWsStore)

  const peerConnections = ref<Record<string, PeerConnection>>({})
  const isMicrophoneOn = ref<boolean>(false)

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

  async function processNext(peerUuid: string, timeoutId?: number) {
    const task = peerConnections.value[peerUuid].taskQueue.shift()
    if (!timeoutId) {
      timeoutId = setTimeout(() => {
        throw new Error('exceeded time limit for task')
      }, 6000)
    }
    if (task) {
      console.log(task)
      try {
        await task()
      } catch (e) {
        console.log(e)
      }
    }
    clearTimeout(timeoutId)

    if (peerConnections.value[peerUuid].taskQueue.length > 0) {
      console.log(timeoutId)
      await processNext(peerUuid, timeoutId)
    }
  }

  async function enqueueTask(peerUuid: string, task: (...args: any) => Promise<void>) {
    peerConnections.value[peerUuid].taskQueue.push(task)

    console.log(peerConnections.value[peerUuid].taskQueue)

    await processNext(peerUuid)
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
    enqueueTask(peerUuid, async () => {
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
    })

    /*    peerConnections.value[peerUuid].pc
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
      }) */
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
        console.log(peerConnections.value[peerUuid].microphoneSender)
      }
    }
  }

  async function shareScreen(peerUuid: string) {
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

      //форсим негошиейшен типа я лидер пиздец а ты пидор ебаный ПДВЗАЗХВЛПХЩЗРВЩЛЗХАЗЩЛВПЗЩОШЛАЩЗЛХВПЗОЩШЩЗ
      enqueueTask(peerUuid, async () => {
        const offer = await peerConnections.value[peerUuid].pc.createOffer({
          offerToReceiveAudio: true,
          offerToReceiveVideo: true,
        })
        await createdDescription(offer, peerUuid)
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
      taskQueue: [],
      pc: new RTCPeerConnection(PEER_CONNECTION_CFG),
    }

    peerConnection.pc.ontrack = (e) => {
      if (!e.streams || !(e.streams.length > 0)) {
        return
      }
      console.log(peerConnections.value[peerUuid].trackMetadata)
      if (e.streams[0].getAudioTracks().length > 0) {
        peerConnections.value[peerUuid].microphoneStream = e.streams[0]
      } /* (peerConnections.value[peerUuid].trackMetadata == 'screen-share') */ else {
        peerConnections.value[peerUuid].ssVideoStream = e.streams[0]
      }
    }

    peerConnection.pc.onnegotiationneeded = async (e) => {
      /*       try {
        const leaderId = [localUuid.value, peerUuid].sort()[0]

        const offer = await peerConnections.value[peerUuid].pc.createOffer({
          offerToReceiveAudio: true,
          offerToReceiveVideo: true,
        })
        if (leaderId == localUuid.value) {
          await createdDescription(offer, peerUuid)
        }
      } catch (error) {
        console.log(error)
      } */

      enqueueTask(peerUuid, async () => {
        try {
          const leaderId = [localUuid.value, peerUuid].sort()[0]
          console.log(peerUuid, localUuid.value)
          const offer = await peerConnections.value[peerUuid].pc.createOffer({
            offerToReceiveAudio: true,
            offerToReceiveVideo: true,
          })
          console.log(peerConnections.value[peerUuid].trackMetadata)
          if (leaderId == localUuid.value) {
            console.log('fdfhgfgj')
            await createdDescription(offer, peerUuid)
          }
        } catch (error) {
          console.log(error)
        }
      })
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

        roomWs.value!.send(
          JSON.stringify({
            uuid: localUuid.value,
            type: 'get-tracks',
          }),
        )
      }
    }

    peerConnection.pc.onicecandidate = (event) => gotIceCandidate(event, peerUuid)

    peerConnections.value[peerUuid] = peerConnection
    if (initCall) {
      /*     peerConnection.pc
        .createOffer({ offerToReceiveAudio: true, offerToReceiveVideo: true })
        .then(async (description) => await createdDescription(description, peerUuid))
        .catch((e) => console.log(e)) */
      enqueueTask(peerUuid, async () =>
        peerConnection.pc
          .createOffer({ offerToReceiveAudio: true, offerToReceiveVideo: true })
          .then(async (description) => await createdDescription(description, peerUuid))
          .catch((e) => console.log(e)),
      )
    }
  }

  async function handleSdpSignal(signal: any, peerUuid: string) {
    /* enqueueTask() */

    enqueueTask(peerUuid, async () => {
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
    })
    /*     peerConnections.value[peerUuid].pc
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
      }) */
  }

  function handleIceCandidate(signal: any, peerUuid: string): void {
    console.log(peerConnections.value[peerUuid].pc.signalingState)
    enqueueTask(peerUuid, async () => {
      peerConnections.value[peerUuid].pc
        .addIceCandidate(new RTCIceCandidate(signal.ice))
        .catch((e) => console.log(e))
    })
    /*     peerConnections.value[peerUuid].pc
      .addIceCandidate(new RTCIceCandidate(signal.ice))
      .catch((e) => console.log(e)) */
  }

  return {
    peerConnections,
    localUuid,
    localDisplayName,
    isMicrophoneOn,
    microphoneStream,
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
    shareMicrophone,
    setTrackMetadata,
  }
})
