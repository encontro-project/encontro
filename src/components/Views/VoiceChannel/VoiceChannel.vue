<script setup lang="ts">
import { useConnectionsStore } from '@/stores/rtcConnections'
import { storeToRefs } from 'pinia'
import { onUnmounted, ref, watch } from 'vue'
import User from './User.vue'
import { useRoomWsStore } from '@/stores/wsConnection'
import MicrophoneIcon from "../../Icons/MicrophoneIcon.vue"
import MicrophoneOffIcon from '../../Icons/MicrophoneOffIcon.vue'

const rtcConnectionsStore = useConnectionsStore()

const roomWsConnectionStore = useRoomWsStore()


const {
  setupPeer,
  handlePeerDisconnect,
  handleIceCandidate,
  handleSdpSignal,
  updateStream,
  shareScreen,
  stopStream,
  unsubscribeFromStream,
  leaveCall,
  toggleMicrophone,
  shareMicrophone,
  setTrackMetadata,
} = rtcConnectionsStore



const {closeRoomWsConnection, getMediaTracks} = roomWsConnectionStore

const localVideo = ref<null | HTMLVideoElement>(null)

const {roomWs, isWsConnected, localUuid, localDisplayName, localStream} = storeToRefs(roomWsConnectionStore)

const {  peerConnections, isMicrophoneOn } = storeToRefs(rtcConnectionsStore)

onUnmounted(() => {
  roomWs.value?.send(
    JSON.stringify({
      type: 'peer-disconnect',
      uuid: localUuid.value,
    }),
  )
  closeRoomWsConnection()
  
})


const gotMessageFromServer = async (message: MessageEvent) => {
    const signal = JSON.parse(message.data)
    const peerUuid = signal.uuid


    if (!roomWs.value) {
      return
    }

    /* console.log(signal, peerConnections.value, peerUuid) */
    if (signal.type === 'peer-disconnect') {
      handlePeerDisconnect(peerUuid)
        
      return
    }
    if (signal.type == 'get-tracks' && peerUuid != localUuid.value) {
      if (!peerConnections.value[peerUuid].ssVideoSender) {
        shareScreen(peerUuid)
      }
    }
    if (signal.type == "get-microphone" && peerUuid != localUuid.value) {
      if (!peerConnections.value[peerUuid].microphoneStream) {
        shareMicrophone(peerUuid)
      }
    } 
    if (signal.type == 'stop-stream' && peerUuid != localUuid.value) {
      if (peerConnections.value[peerUuid].ssVideoStream) {
        unsubscribeFromStream(peerUuid)
      }
    }
    if (peerUuid === localUuid.value || (signal.dest !== localUuid.value && signal.dest !== 'all'))
    return
  
  if (signal.metadata) {
    setTrackMetadata(peerUuid, signal.metadata)
  }
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
    console.log("dest: ", signal.dest, "uuid: ", peerUuid)
    await handleSdpSignal(signal, peerUuid)
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
<div class="users-container">
  <button @click="handleStreamChange">stream change</button>
  <button @click="handleStreamStart">share screen</button>
  <button @click="handleStreamStop">stop stream</button>
  <button @click="handleLeaveCall">Quit Call</button>
  <div class="videos-container">
    <div class="user-container" v-if="isWsConnected">
      <div class="video-container" v-if="localStream">
        <video ref="localVideo" autoplay muted class="rtc-stream"></video>
      </div>
      <div  class="video-template" v-else>
        <img class="user-pic" src="https://e7.pngegg.com/pngimages/719/959/png-clipart-celebes-crested-macaque-monkey-selfie-grapher-people-for-the-ethical-treatment-of-animals-funny-mammal-animals-thumbnail.png"></img>
      </div>
      <div class="user-info">
        <div class="microphone-status">
          <MicrophoneIcon class="mic-icon" v-if="isMicrophoneOn" v-on:click="toggleMicrophone"></MicrophoneIcon>
          <MicrophoneOffIcon class="mic-icon" v-else-if="!isMicrophoneOn" v-on:click="toggleMicrophone"></MicrophoneOffIcon>
        </div>
        <div class="display-name">
          {{ localDisplayName }}
        </div>
      </div>
    </div>
    <User
      v-for="(peerConnection, i) in peerConnections"
      :key="i"
      :peer-connection="peerConnection"
    ></User>
  </div>
  </div>
</template>

<style>
.videos-container {
  display: grid;
  max-width: 100vw;
  grid-template-columns: repeat(3, 600px);
  gap: 20px;
  
}

.users-container {
  display: block;
}

.user-container {
  position: relative;
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
.user-info {
  position: absolute;
  left: 5px;
  background-color: rgba(0, 0, 0, 0.6);
  color: white;
  display: flex;
  gap: 5px;
  padding: 3px 3px 3px 3px;
  align-items: center;
  bottom: 5px;
  border-radius: 6px;
}
.mic-icon {
  width: 24px;
  height: 24px;
  stroke: white;
}
.display-name {
  font-size: 18px;
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
