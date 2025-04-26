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
        state.peerConnections[peerUuid].pc.close()
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
    async initializeCodecs() {
      if (!('MediaStreamTrackProcessor' in window)) {
        console.log('MediaStreamTrackProcessor not found')
        return
      }

      console.log(MediaStreamTrackGenerator)
      const videoTrack = this.localStream?.getVideoTracks()[0]
      const audioTrack = this.localStream?.getAudioTracks()[0]

      const videoProcessor = new MediaStreamTrackProcessor({
        track: videoTrack as MediaStreamVideoTrack,
      })
      const videoReader = videoProcessor.readable.getReader()

      this.videoEncoder = new VideoEncoder({
        output: (encodedChunk: EncodedVideoChunk) => {
          const data = new ArrayBuffer(encodedChunk.byteLength)
          const view = new Uint8Array(data)
          encodedChunk.copyTo(view)

          Object.values(this.peerConnections).forEach((peer) => {
            if (peer.dataChannel?.readyState === 'open') {
              peer.dataChannel.send(
                JSON.stringify({
                  type: encodedChunk.type,
                  timestamp: encodedChunk.timestamp,
                  duration: encodedChunk.duration,
                }),
              )
              peer.dataChannel.send(data)
            }
          })
        },
        error: (e) => console.error(e),
      })

      this.videoEncoder.configure(VIDEO_ENCODER_CFG)

      const processVideoFrames = async () => {
        try {
          while (true) {
            const { value: frame, done } = await videoReader.read()
            if (done) break
            if (this.videoEncoder?.state === 'configured') {
              this.videoEncoder?.encode(frame)
            }
            frame.close()
          }
          await this.videoEncoder?.flush()
          this.videoEncoder?.close()
        } catch (e) {
          console.error('Error processing frames:', e)
        }
      }

      processVideoFrames()

      if (audioTrack && 'AudioEncoder' in window) {
        console.log(audioTrack.getSettings())
        const audioProcessor = new MediaStreamTrackProcessor({ track: audioTrack })
        const audioReader = audioProcessor.readable.getReader()

        this.audioEncoder = new AudioEncoder({
          output: (encodedChunk: EncodedAudioChunk) => {
            const data = new ArrayBuffer(encodedChunk.byteLength)
            const view = new Uint8Array(data)
            encodedChunk.copyTo(view)

            Object.values(this.peerConnections).forEach((peer) => {
              if (peer.dataChannel?.readyState === 'open') {
                peer.dataChannel.send(
                  JSON.stringify({
                    type: 'audio',
                    timestamp: encodedChunk.timestamp,
                    duration: encodedChunk.duration,
                  }),
                )
                peer.dataChannel.send(data)
              }
            })
          },
          error: (e) => console.error('Audio Encoder Error:', e),
        })

        this.audioEncoder.configure(AUDIO_ENCODER_CFG)

        const processAudioFrames = async () => {
          try {
            while (true) {
              const { value: frame, done } = await audioReader.read()
              if (done) break
              if (this.audioEncoder?.state === 'configured' && frame) {
                try {
                  this.audioEncoder.encode(frame)
                } catch (e) {
                  console.error('Error encoding audio frame:', e)
                }
              }
              frame.close()
            }
          } catch (e) {
            console.error('Error processing audio frames:', e)
          }
        }

        processAudioFrames()
      }
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
      const state = this.peerConnections[peerUuid]?.pc.iceConnectionState
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

    setupPeer(
      serverConnection: WebSocket,
      peerUuid: string,
      displayName: string,
      initCall = false,
    ) {
      const peerConnection: PeerConnection = {
        displayName,
        uuid: this.localUuid,
        pc: new RTCPeerConnection(PEER_CONNECTION_CFG),
      }

      const onDataChannelOpen = async () => {
        if (!('VideoDecoder' in window)) {
          console.log('no videoDecoder')
          return
        }

        console.log('videoDecoder')
        peerConnection.videoDecoder = new VideoDecoder({
          output: (frame: VideoFrame) => {
            if (!videoStreamGenerator) {
              videoStreamGenerator = new MediaStreamTrackGenerator({
                kind: 'video',
              })
              videoWritable = videoStreamGenerator.writable.getWriter()
              this.peerConnections[peerUuid].videoStream = new MediaStream([videoStreamGenerator])
            }

            videoWritable!.write(frame)
          },
          error: (e) => console.error(e),
        })

        peerConnection.videoDecoder.configure(VIDEO_DECODER_CFG)

        if ('AudioDecoder' in window) {
          console.log('AudioDecoder')
          peerConnection.audioDecoder = new AudioDecoder({
            output: async (audioData: AudioData) => {
              if (!audioStreamGenerator) {
                audioStreamGenerator = new MediaStreamTrackGenerator({
                  kind: 'audio',
                })
                audioWritable = audioStreamGenerator.writable.getWriter()
              }
              this.peerConnections[peerUuid].audioStream = new MediaStream([audioStreamGenerator])
              audioWritable!.write(audioData)
            },
            error: (e) => console.error('Audio decoder error:', e),
          })

          peerConnection.audioDecoder.configure(AUDIO_DECODER_CFG)
        }

        this.initializeCodecs()
      }

      const onDataChannelMessage = async (event: MessageEvent) => {
        if (typeof event.data === 'string') {
          metadata = JSON.parse(event.data)
        } else if (metadata.type === 'delta' || metadata.type == 'key') {
          if (metadata.type == 'key') {
          }
          const chunk = new EncodedVideoChunk({
            type: 'key',
            timestamp: metadata.timestamp,
            duration: metadata.duration,
            data: event.data,
          })

          if (peerConnection.videoDecoder?.state === 'configured') {
            try {
              peerConnection.videoDecoder.decode(chunk)
            } catch (error) {
              console.log('decoding error')
            }
          }
        } else if (event.data.type === 'audio') {
          const chunk = new EncodedAudioChunk({
            type: 'key',
            timestamp: metadata.timestamp,
            duration: metadata.duration,
            data: event.data,
          })

          if (peerConnection.audioDecoder?.state === 'configured') {
            peerConnection.audioDecoder.decode(chunk)
          }
        }
      }

      const dataChannel = peerConnection.pc.createDataChannel('media-channel')
      peerConnection.dataChannel = dataChannel

      let metadata: any

      let videoWritable: WritableStreamDefaultWriter<VideoFrame> | undefined
      let videoStreamGenerator: MediaStreamTrackGenerator<VideoFrame> | undefined

      let audioStreamGenerator: MediaStreamTrackGenerator<AudioData> | undefined
      let audioWritable: WritableStreamDefaultWriter<AudioData> | undefined

      peerConnection.dataChannel.onopen = onDataChannelOpen

      peerConnection.dataChannel.onmessage = onDataChannelMessage
      peerConnection.pc.onicecandidate = (event) =>
        this.gotIceCandidate(event, peerUuid, serverConnection)
      peerConnection.pc.oniceconnectionstatechange = () => this.checkPeerDisconnect(peerUuid)
      if (this.localStream) {
        console.log(this.localStream.getTracks())
        this.localStream.getTracks().forEach((track) => {
          peerConnection.pc.addTrack(track, this.localStream as MediaStream)
        })
      }

      if (initCall) {
        peerConnection.pc.ondatachannel = (event) => {
          peerConnection.dataChannel = event.channel
          peerConnection.dataChannel.onmessage = onDataChannelMessage
          peerConnection.dataChannel.onopen = onDataChannelOpen
        }

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
