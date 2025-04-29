import { v4 as uuidv4 } from 'uuid'
import { defineStore } from 'pinia'
import type { PeerConnection } from '@/types'

//TODO (1) Переписать на composition api

//TODO (1) Использовать MediaStreams api вместо dataChannels

interface IState {
  peerConnections: {
    [key: string]: PeerConnection
  }
  audioEncoder: AudioEncoder | null
  videoEncoder: VideoEncoder | null
  localStream: MediaStream | null
  localUuid: string
  localDisplayName: string
}

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

const VIDEO_ENCODER_CFG: VideoEncoderConfig = {
  codec: 'vp8',
  width: VIDEO_STREAM_WIDTH,
  height: VIDEO_STREAM_HEIGHT,
  framerate: VIDEO_STREAM_FRAME_RATE,
  bitrate: 1_000_000,
  latencyMode: 'realtime',
}

const VIDEO_DECODER_CFG: VideoDecoderConfig = {
  codec: VIDEO_ENCODER_CFG.codec,
  codedWidth: VIDEO_ENCODER_CFG.width,
  codedHeight: VIDEO_ENCODER_CFG.height,
  hardwareAcceleration: 'no-preference',
}

const AUDIO_ENCODER_CFG: AudioEncoderConfig = {
  codec: 'opus',
  numberOfChannels: 2,
  sampleRate: 48000,
  bitrate: 128_000,
}

const AUDIO_DECODER_CFG: AudioDecoderConfig = {
  codec: AUDIO_ENCODER_CFG.codec,
  numberOfChannels: AUDIO_ENCODER_CFG.numberOfChannels,
  sampleRate: AUDIO_ENCODER_CFG.sampleRate,
}

const generateUUID = () => {
  const uuid = uuidv4()
  localStorage.setItem('uuid', JSON.stringify(uuid))
  localStorage.setItem('displayName', JSON.stringify(uuid))
  return uuid
}

export const useConnectionsStore = defineStore('connections', {
  state: (): IState => {
    return {
      peerConnections: {},
      audioEncoder: null,
      videoEncoder: null,
      localStream: null,
      localUuid: localStorage.getItem('uuid')
        ? JSON.parse(localStorage.getItem('uuid') as string)
        : generateUUID(),
      localDisplayName: localStorage.getItem('displayName')
        ? JSON.parse(localStorage.getItem('displayName') as string)
        : generateUUID(),
    }
  },
  actions: {
    async getMediaTracks() {
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
        this.localStream = await navigator.mediaDevices.getDisplayMedia(CONSTRAINTS)
        return
      } else if (navigator.mediaDevices.getUserMedia) {
        this.localStream = await navigator.mediaDevices.getUserMedia(CONSTRAINTS)
        return
      } else {
        alert('Your browser does not support getDisplayMedia or getUserMedia API')
        throw new Error('Unsupported media API')
      }
    },
    handlePeerDisconnect(peerUuid: string) {
      console.log(`Peer ${peerUuid} disconnected`)
      this.$patch((state: IState) => {
        state.peerConnections[peerUuid]?.pc.close()
        delete state.peerConnections[peerUuid]
      })
    },
    updateStream() {
      Object.values(this.peerConnections).forEach((connection) =>
        connection.pc.getSenders().forEach((sender) => {
          try {
            const videoTrack = this.localStream?.getVideoTracks()[0]
            const audioTrack = this.localStream?.getAudioTracks()[0]
            console.log(audioTrack)
            if (sender.track?.kind == 'video') {
              sender.replaceTrack(videoTrack as MediaStreamTrack)
            }
            if (sender.track?.kind == 'audio') {
              sender.replaceTrack(audioTrack as MediaStreamTrack)
            }
          } catch {}
        }),
      )
    },

    gotIceCandidate(
      event: RTCPeerConnectionIceEvent,
      peerUuid: string,
      serverConnection: WebSocket,
    ) {
      if (event.candidate) {
        serverConnection!.send(
          JSON.stringify({
            ice: event.candidate,
            uuid: this.localUuid,
            dest: peerUuid,
          }),
        )
      }
    },

    checkPeerDisconnect(peerUuid: string): void {
      const state = this.peerConnections[peerUuid]?.pc.connectionState
      if (['failed', 'closed', 'disconnected'].includes(state)) {
        this.handlePeerDisconnect(peerUuid)
      }
    },
    createdDescription(
      description: RTCSessionDescriptionInit,
      peerUuid: string,
      serverConnection: WebSocket,
    ): void {
      this.peerConnections[peerUuid].pc
        .setLocalDescription(description)
        .then(() => {
          serverConnection!.send(
            JSON.stringify({
              sdp: this.peerConnections[peerUuid].pc.localDescription,
              uuid: this.localUuid,
              dest: peerUuid,
            }),
          )
        })
        .catch((e: Error) => console.log(e))
    },

    shareScreen(peerUuid: string) {
      if (this.localStream) {
        console.log(this.localStream.getTracks())
        this.localStream.getTracks().forEach((track) => {
          this.peerConnections[peerUuid].pc.addTrack(track, this.localStream as MediaStream)
        })
      }
    },

    setupPeer(
      serverConnection: WebSocket,
      peerUuid: string,
      displayName: string,
      initCall: boolean,
    ) {
      const peerConnection: PeerConnection = {
        displayName,
        uuid: peerUuid,
        pc: new RTCPeerConnection(PEER_CONNECTION_CFG),
      }

      peerConnection.pc.onnegotiationneeded = () => {
        peerConnection.pc
          .createOffer({ offerToReceiveAudio: true, offerToReceiveVideo: true })
          .then((description) => this.createdDescription(description, peerUuid, serverConnection))
          .catch((e) => console.log(e))
      }

      peerConnection.pc.onconnectionstatechange = (e) => {
        this.checkPeerDisconnect(peerUuid)
        const state = this.peerConnections[peerUuid]?.pc.connectionState
        if (state == 'connected') {
          serverConnection.send(
            JSON.stringify({
              uuid: this.localUuid,
              type: 'get-tracks',
            }),
          )
        }
      }

      peerConnection.pc.onicecandidate = (event) =>
        this.gotIceCandidate(event, peerUuid, serverConnection)

      if (initCall) {
        peerConnection.pc
          .createOffer({ offerToReceiveAudio: true, offerToReceiveVideo: true })
          .then((description) => this.createdDescription(description, peerUuid, serverConnection))
          .catch((e) => console.log(e))
      }

      this.peerConnections[peerUuid] = peerConnection
    },
    handleSdpSignal(signal: any, peerUuid: string, serverConnection: WebSocket): void {
      this.peerConnections[peerUuid].pc
        .setRemoteDescription(new RTCSessionDescription(signal.sdp))
        .then(() => {
          if (signal.sdp.type === 'offer') {
            this.peerConnections[peerUuid].pc
              .createAnswer()
              .then((description) =>
                this.createdDescription(description, peerUuid, serverConnection),
              )
              .catch((e) => console.log(e))
          }
        })
        .catch((e) => console.log(e))
    },
    handleIceCandidate(signal: any, peerUuid: string): void {
      this.peerConnections[peerUuid].pc
        .addIceCandidate(new RTCIceCandidate(signal.ice))
        .catch((e) => console.log(e))
    },
  },

  getters: {},
})

type ConnectionsStore = ReturnType<typeof useConnectionsStore>
