<script setup lang="ts">
import { useConnectionsStore } from '@/stores/connections'
import { storeToRefs } from 'pinia'
import { onUnmounted, ref, watch, watchEffect } from 'vue'
import TestVideo from './TestVideo.vue'

const wsRef = ref<null | WebSocket>(null)

const isConnected = ref(false)

const store = useConnectionsStore()

const {
  getMediaTracks,
  initializeCodecs,
  setupPeer,
  handlePeerDisconnect,
  handleIceCandidate,
  handleSdpSignal,
  updateStream,
} = store

const localVideo = ref<null | HTMLVideoElement>(null)

const { localStream, localUuid, peerConnections, localDisplayName } = storeToRefs(store)

onUnmounted(() => {
  wsRef.value?.send(
    JSON.stringify({
      type: 'peer-disconnect',
      uuid: localUuid,
    }),
  )
})

const gotMessageFromServer = (message: MessageEvent) => {
  const signal = JSON.parse(message.data)
  const peerUuid = signal.uuid

  if (!wsRef.value) {
    return
  }

  if (peerUuid === localUuid.value || (signal.dest !== localUuid.value && signal.dest !== 'all'))
    return
  console.log('msg', peerConnections.value, peerUuid)

  if (signal.type === 'peer-disconnect') {
    handlePeerDisconnect(peerUuid)
    return
  }

  if (signal.displayName && signal.dest === 'all') {
    setupPeer(wsRef.value, peerUuid, signal.displayName)
    wsRef.value!.send(
      JSON.stringify({
        displayName: localDisplayName.value,
        uuid: localUuid.value,
        dest: peerUuid,
      }),
    )
  } else if (signal.displayName && signal.dest === localUuid.value) {
    setupPeer(wsRef.value, peerUuid, signal.displayName, true)
  } else if (signal.sdp) {
    handleSdpSignal(signal, peerUuid, wsRef.value)
  } else if (signal.ice) {
    handleIceCandidate(signal, peerUuid)
  }
}

const getTracks = async () => {
  const WS_PORT = 8443
  await getMediaTracks()
  if (localVideo.value) {
    localVideo.value.srcObject = localStream.value
  }
  console.log(localStream)
}

const initWebSocket = () => {
  wsRef.value = new WebSocket('wss://localhost:8443/ws/room')

  wsRef.value.onopen = () => {
    wsRef.value!.send(
      JSON.stringify({
        displayName: localDisplayName,
        uuid: localUuid.value,
        dest: 'all',
      }),
    )
    isConnected.value = true
  }

  wsRef.value.onclose = () => {
    isConnected.value = false
  }
}

const handleConnectionStart = async () => {
  await getTracks()
  initWebSocket()
}

const handleStreamChange = async () => {
  await getTracks()
  updateStream()
}

watch(
  wsRef,
  (newSocket, oldSocket) => {
    if (oldSocket) {
      oldSocket.close()
      console.log('cleanup')
    }
    if (newSocket) {
      console.log('niggers')
      newSocket.onmessage = (e) => {
        gotMessageFromServer(e)
      }
      console.log(newSocket)
    }
  },
  { deep: false },
)
</script>

<template>
  <button @click="handleConnectionStart">get tracks</button>
  <button @click="initializeCodecs">Пенис</button>
  <button @click="handleStreamChange">stream change</button>
  <div class="videos-container">
    <div class="video-container" v-if="localStream?.active">
      <video ref="localVideo" autoplay class="rtc-stream"></video>
    </div>
    <TestVideo
      v-for="(peerConnection, i) in peerConnections"
      :key="i"
      :peer-connection="peerConnection"
    ></TestVideo>
  </div>
</template>

<style>
.videos-container {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}
.video-container {
  background-color: black;
  width: 600px;
  height: 400px;
}
.rtc-stream {
  width: 600px;
  height: 400px;
}
</style>
