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
  toggleMicrophone,
} = rtcConnectionsStore


const localVideo = ref<null | HTMLVideoElement>(null)

const {isWsConnected, localDisplayName, localStream} = storeToRefs(roomWsConnectionStore)

const {  peerConnections, isMicrophoneOn } = storeToRefs(rtcConnectionsStore)








watch([localStream, localVideo], ([newStream, newVideo], [oldStream, oldVideo]) => {
  if (newStream) {
    if (newVideo) {
    newVideo.srcObject = localStream.value
  }
  }
})


</script>

<template>
<div class="users-container">
  <div class="videos-container">
    <div class="user-container" v-if="isWsConnected == 'connected'">
      <div class="video-container" v-if="localStream">
        <video ref="localVideo" autoplay muted class="rtc-stream"></video>
      </div>
      <div  class="video-template" v-else>
        <img class="user-pic" src="https://masterpiecer-images.s3.yandex.net/5fcb1cda5223d2d:upscaled"></img>
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
  position: relative;
  margin-top: 20px;
  display: grid;
  max-width: 80%;
  grid-template-columns: repeat(2, 600px);
  gap: 20px;
  
}

.users-container {
  display: flex;
  align-items: center;
  flex-direction: column;
  
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
  background-color: white;
  width: 600px;
  height: 400px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius:  5px;
}
.user-pic {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  border: 3px solid black;
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
@media (max-width: 1239px) {
  .videos-container{
    grid-template-columns: repeat(1, 600px);
  }
}
</style>
