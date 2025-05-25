<script lang="ts" setup>
import type { PeerConnection } from '@/types'
import { nextTick, ref, toRef, watch, watchEffect } from 'vue'
import MicrophoneIcon from "../../Icons/MicrophoneIcon.vue"
import MicrophoneOffIcon from '../../Icons/MicrophoneOffIcon.vue'
import { useConnectionsStore } from '@/stores/rtcConnections'
import Slider from 'primevue/slider'
import { useContextMenuStore } from '@/stores/contextMenu'
import { storeToRefs } from 'pinia'


interface Props {
  peerConnection: PeerConnection
}
const props = defineProps<Props>()

const rtcConnectionsStore = useConnectionsStore()

const contextMenuStore = useContextMenuStore()

const {isMenuActive} = storeToRefs(contextMenuStore)

const {openMenu} = contextMenuStore 

const {changeUserVolume} = rtcConnectionsStore

const videoRef = ref<HTMLVideoElement>()

const userContainer = ref<HTMLDivElement>()

const sliderValue = toRef(props.peerConnection, "userVolume")

const contextMenuRef = ref<HTMLDivElement>()

const videoStreamRef = toRef(props.peerConnection, "ssVideoStream")

watch([videoRef, videoStreamRef], ([newVideo, newStream], [oldVideo, oldStream]) => {
  if (oldStream) {
  }
  if (newStream && videoRef.value) {
    videoRef.value.srcObject = newStream
  }
})

const toggleDisplayContextMenu = async (e: MouseEvent) => {
  if (e.button == 2) {
    openMenu()
    await nextTick() // Ждем следующий тик чтобы контекст меню реф инициализировался, кто тронет, того маму ебал 
    const rect = userContainer.value!.getBoundingClientRect();
    contextMenuRef.value!.style.top = e.clientY - rect.top  + "px"
    contextMenuRef.value!.style.left = e.clientX - rect.left  + "px"
  }
}

const handleUserVolumeChange = () => {
  changeUserVolume(props.peerConnection.uuid, sliderValue.value)
}

//UI библиотека для слайдеров


</script>

<template>
  <div class="user-container" @mousedown.stop="toggleDisplayContextMenu" @contextmenu.prevent="" ref="userContainer">
    <div class="context-menu-container" ref="contextMenuRef" v-if="isMenuActive" @click.stop="">
      <div class="context-menu">
        <div class="menu-display-name">
          <div class="menu-item">
            <p>Профиль</p>
          </div>
          <p>Громкость пользователя</p>
        </div>
        <div class="slider-container">
          <Slider :min="0" :max="1" :step="0.01" class="volume-input" v-model="sliderValue" @value-change="handleUserVolumeChange()"></Slider>
        </div>
      </div>
    </div>
    <div class="video-container" v-if="videoStreamRef">
      <video ref="videoRef" autoplay class="rtc-stream"></video>
    </div>
    <div v-else class="video-template">
      <img class="user-pic" src="https://masterpiecer-images.s3.yandex.net/5fcb1cda5223d2d:upscaled"></img>
    </div>
    <div class="user-info">
      <div class="microphone-status">
        <MicrophoneIcon class="mic-icon" v-if="!peerConnection.isMuted"></MicrophoneIcon>
        <MicrophoneOffIcon class="mic-icon" v-else></MicrophoneOffIcon>
      </div>
      <div class="display-name">
        {{ props.peerConnection.displayName }}
      </div>
    </div>
  </div>
</template>

<style scoped>

.volume-input {
  width: 90%;
}

.slider-container {
  position: relative;
  color: white;
  margin-top: 20px;
  width: 100%;
  height: 16px;
  display: flex;
  justify-content: center;
  flex-direction: column;
  align-items: center;
  border-radius: 5px;
}

:deep(.volume-input.p-slider) {
  height: 8px; 
  background-color: white;
}

/* Active progress portion */
:deep(.volume-input .p-slider-range) {
  height: 8px; /* Must match container height */
  background: #ff004c; /* Progress color */
}

/* Handle styling - removes vertical line */
:deep(.volume-input .p-slider-handle) {
  width: 16px;
  height: 16px;
  background: white;
  margin-top: -8px; /* Centers vertically */
  border-radius: 50%; /* Makes it circular */
}

/* Removes the vertical line inside handle */
:deep(.volume-input .p-slider-handle::before) {
  content: none !important;
  display: none !important;
}
:deep(.volume-input .p-slider-handle::after) {
  display: none;
}


.context-menu-container {
  border-radius: 5px;
  border: 2px solid white;
  display: flex;
  position: absolute;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 20px;
  background-color: black;
  width: 250px;
  height: 400px;
  z-index: 2;
}

.context-menu {
  width: 90%;
  height: 90%;
}

.menu-display-name  {
  font-weight: 500;
  color: white;
  text-align: start;
  display: flex;
  gap: 20px;
  flex-direction: column;
}

.menu-item {
  cursor: pointer;
  height: 35px;
  width: 100%;
  display: flex;
  justify-content: start;
  align-items: center;
  background-color: black;
  border-radius: 5px;
}


.menu-item:hover {
  background-color: white;
  color: black;
}

.menu-display-name p {
  margin-top: 0;
  margin-bottom: 0;
  padding-left: 10px;
}
</style>
