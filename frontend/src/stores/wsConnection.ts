import { defineStore, storeToRefs } from 'pinia'
import { ref } from 'vue'
import { useConnectionsStore } from './rtcConnections'
import { v4 as uuidv4 } from 'uuid'

const VIDEO_STREAM_WIDTH: number = 1920
const VIDEO_STREAM_HEIGHT: number = 1080
const VIDEO_STREAM_FRAME_RATE: number = 60
const CONSTRAINTS: MediaStreamConstraints = {
  audio: true,
  video: {
    width: { ideal: VIDEO_STREAM_WIDTH, max: VIDEO_STREAM_WIDTH },
    height: { ideal: VIDEO_STREAM_HEIGHT, max: VIDEO_STREAM_HEIGHT },
    frameRate: { ideal: VIDEO_STREAM_FRAME_RATE, max: VIDEO_STREAM_FRAME_RATE },
  },
}

const generateUUID = () => {
  const uuid = uuidv4()
  localStorage.setItem('uuid', JSON.stringify(uuid))
  localStorage.setItem('displayName', JSON.stringify(uuid))
  return uuid
}

export const useRoomWsStore = defineStore('roomWsStore', () => {
  const microphoneStream = ref<MediaStream | null>(null)
  const isWsConnected = ref<boolean>(false)
  const roomWs = ref<WebSocket | null>(null)
  const currentRoomUrl = ref<string>('')
  const localStream = ref<MediaStream | null>(null)
  const localUuid = ref<string>(
    localStorage.getItem('uuid')
      ? JSON.parse(localStorage.getItem('uuid') as string)
      : generateUUID(),
  )

  const localDisplayName = ref<string>(
    localStorage.getItem('displayName')
      ? JSON.parse(localStorage.getItem('displayName') as string)
      : generateUUID(),
  )

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

  const initWebSocket = async (room: string) => {
    await getMicrophoneTrack()

    currentRoomUrl.value = room
    roomWs.value = new WebSocket(`wss://localhost:8443/ws/${room}`)

    roomWs.value.onopen = () => {
      roomWs.value!.send(
        JSON.stringify({
          displayName: localDisplayName.value,
          uuid: localUuid.value,
          dest: 'all',
        }),
      )

      isWsConnected.value = true
    }

    roomWs.value.onclose = () => {
      isWsConnected.value = false
    }
  }

  const closeRoomWsConnection = () => {
    roomWs.value?.close()
    roomWs.value = null
  }
  return {
    roomWs,
    isWsConnected,
    localUuid,
    localDisplayName,
    microphoneStream,
    localStream,
    currentRoomUrl,
    getMicrophoneTrack,
    initWebSocket,
    closeRoomWsConnection,
    getMediaTracks,
  }
})
