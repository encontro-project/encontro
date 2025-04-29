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
  setupPeer,
  handlePeerDisconnect,
  handleIceCandidate,
  handleSdpSignal,
  updateStream,
  shareScreen
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

  console.log(signal.type, peerConnections.value, peerUuid)
  if (signal.type === 'peer-disconnect') {
    handlePeerDisconnect(peerUuid)
    return
  }
  if (signal.type == "get-tracks") {
    shareScreen(peerUuid)
  }
  if (peerUuid === localUuid.value || (signal.dest !== localUuid.value && signal.dest !== 'all'))
  return


  if (signal.displayName && signal.dest === 'all') {
    setupPeer(wsRef.value, peerUuid, signal.displayName, false)
    wsRef.value?.send(
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
  initWebSocket()
}

const handleStreamChange = async () => {
  await getTracks()
  updateStream()
}

const handleStreamStart = async () => {
  await getTracks()
  Object.values(peerConnections.value).forEach((peerConnection) => {
    shareScreen(peerConnection.uuid)
  })
}

watch(
  wsRef,
  (newSocket, oldSocket) => {
    if (oldSocket) {
      oldSocket.close()
      console.log('cleanup')
    }
    if (newSocket) {
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
  <button @click="handleConnectionStart">connect</button>
  <button @click="handleStreamChange">stream change</button>
  <button @click="handleStreamStart">share screen</button>
  <div class="videos-container">
    <div class="video-container" v-if="localStream?.active">
      <video ref="localVideo" autoplay class="rtc-stream"></video>
    </div>
    <div v-else-if="isConnected" class="video-template">
      <img class="user-pic" src="https://e7.pngegg.com/pngimages/719/959/png-clipart-celebes-crested-macaque-monkey-selfie-grapher-people-for-the-ethical-treatment-of-animals-funny-mammal-animals-thumbnail.png"></img>
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
.video-template {
  background-color: aqua;
  width: 600px;
  height: 400px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
.user-pic {
  width: 100px;
  height: 100px;
  border-radius: 50%;
}
.rtc-stream {
  width: 600px;
  height: 400px;
}
</style>
