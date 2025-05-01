<script setup lang="ts">
import { useConnectionsStore } from '@/stores/rtcConnections'
import { storeToRefs } from 'pinia'
import { onUnmounted, ref, watch, watchEffect } from 'vue'
import TestVideo from './TestVideo.vue'
import { useRoomWsStore } from '@/stores/wsConnection'

const rtcConnectionsStore = useConnectionsStore()

const roomWsConnectionStore = useRoomWsStore()


const {
  getMediaTracks,
  setupPeer,
  handlePeerDisconnect,
  handleIceCandidate,
  handleSdpSignal,
  updateStream,
  shareScreen,
  stopStream,
  unsubscribeFromStream,
  leaveCall
} = rtcConnectionsStore



const {initWebSocket, closeRoomWsConnection} = roomWsConnectionStore

const localVideo = ref<null | HTMLVideoElement>(null)

const {roomWs, isWsConnected} = storeToRefs(roomWsConnectionStore)

const { localStream, localUuid, peerConnections, localDisplayName } = storeToRefs(rtcConnectionsStore)

onUnmounted(() => {

  roomWs.value?.send(
    JSON.stringify({
      type: 'peer-disconnect',
      uuid: localUuid.value,
    }),
  )
  closeRoomWsConnection()
  
})


const gotMessageFromServer = (message: MessageEvent) => {
    const signal = JSON.parse(message.data)
    const peerUuid = signal.uuid

    console.log('zalupa')

    if (!roomWs.value) {
      return
    }

    console.log(signal.type, peerConnections.value, peerUuid)
    if (signal.type === 'peer-disconnect') {
      handlePeerDisconnect(peerUuid)
      return
    }
    if (signal.type == 'get-tracks' && peerUuid != localUuid.value) {
      if (!peerConnections.value[peerUuid].ssVideoStream) {
        shareScreen(peerUuid)
      }
    }
    if (signal.type == 'stop-stream' && peerUuid != localUuid.value) {
      if (peerConnections.value[peerUuid].ssVideoStream) {
        unsubscribeFromStream(peerUuid)
      }
    }
    if (peerUuid === localUuid.value || (signal.dest !== localUuid.value && signal.dest !== 'all'))
      return

    if (signal.displayName && signal.dest === 'all') {
      setupPeer(peerUuid, signal.displayName, false)
      roomWs.value?.send(
        JSON.stringify({
          displayName: localDisplayName.value,
          uuid: localUuid.value,
          dest: peerUuid,
        }),
      )
    } else if (signal.displayName && signal.dest === localUuid.value) {
      setupPeer(peerUuid, signal.displayName, true)
    } else if (signal.sdp) {
      handleSdpSignal(signal, peerUuid)
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


const handleConnectionStart = async (room: string) => {
  initWebSocket(room)
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

const handleStreamStop = () => {
  stopStream()
  roomWs.value?.send(
    JSON.stringify({
      uuid: localUuid.value,
      type: 'stop-stream',
      }),
      )
  
}

const handleLeaveCall = () => {
  if (roomWs.value){
    roomWs.value.send(
      JSON.stringify({
        type: 'peer-disconnect',
        uuid: localUuid.value,
      }),
    )
    leaveCall()
    closeRoomWsConnection()
  }
}

watch(
  roomWs,
  (newSocket, oldSocket) => {
    if (oldSocket) {
      leaveCall()
      oldSocket.send(
      JSON.stringify({
        type: 'peer-disconnect',
        uuid: localUuid.value,
      }),
    )
    setTimeout(() => {

      oldSocket.close()
    }, 100)
    
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

  <div class="roomButtons">
    <button @click="handleConnectionStart('room1')">room1</button>
    <button @click="handleConnectionStart('room2')">room2</button>
  </div>
  <button @click="handleStreamChange">stream change</button>
  <button @click="handleStreamStart">share screen</button>
  <button @click="handleStreamStop">stop stream</button>
  <button @click="handleLeaveCall">Quit Call</button>
  <div class="videos-container">
    <div class="video-container" v-if="localStream">
      <video ref="localVideo" autoplay muted class="rtc-stream"></video>
    </div>
    <div v-else-if="isWsConnected" class="video-template">
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
  max-width: 100vw;
  grid-template-columns: repeat(3, 600px);
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
@media (max-width: 1859px) {
  .videos-container {
    grid-template-columns: repeat(2, 600px);
  }
}
@media (max-width: 1239px) {
  .videos-container{
    grid-template-columns: repeat(1, 600px);
  }
}
</style>
