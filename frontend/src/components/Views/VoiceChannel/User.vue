<script lang="ts" setup>
import type { PeerConnection } from '@/types'
import { ref, toRef, watch, watchEffect } from 'vue'
import MicrophoneIcon from "../../Icons/MicrophoneIcon.vue"
import MicrophoneOffIcon from '../../Icons/MicrophoneOffIcon.vue'

interface Props {
  peerConnection: PeerConnection
}
const props = defineProps<Props>()


const videoRef = ref<HTMLVideoElement | null>(null)

const videoStreamRef = toRef(props.peerConnection, "ssVideoStream")

watch([videoRef, videoStreamRef], ([newVideo, newStream], [oldVideo, oldStream]) => {
  if (oldStream) {
  }
  if (newStream && videoRef.value) {
    videoRef.value.srcObject = newStream
  }
})

</script>

<template>
  <div class="user-container">
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

<style scoped></style>
