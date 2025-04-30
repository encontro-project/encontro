<script lang="ts" setup>
import type { PeerConnection } from '@/types'
import { onMounted, ref, toRef, watch, watchEffect } from 'vue'

interface Props {
  peerConnection: PeerConnection
}
const props = defineProps<Props>()


const videoRef = ref<HTMLVideoElement | null>(null)

const videoStreamRef = toRef(props.peerConnection, "ssVideoStream")

onMounted(() => {
 /*  props.peerConnection.pc.ontrack = (e) => {
    console.log(e)
    if (e.streams && e.streams.length > 0) {
      // Always replace the entire stream
      videoStreamRef.value = e.streams[0]
      console.log(e.streams[0].getAudioTracks())
    }
  } */
})

watch([videoRef, videoStreamRef], ([newVideo, newStream], [oldVideo, oldStream]) => {
  if (oldStream) {
  }
  if (newStream && videoRef.value) {
    console.log(videoStreamRef.value) 
    videoRef.value.srcObject = newStream
  }
})

watchEffect(() => {
  console.log('effect')
})
</script>

<template>
  <div class="video-container" v-if="videoStreamRef">
    <video ref="videoRef" autoplay class="rtc-stream"></video>
    <p>{{ peerConnection.uuid }}</p>
  </div>
  <div v-else class="video-template">
      <img class="user-pic" src="https://e7.pngegg.com/pngimages/719/959/png-clipart-celebes-crested-macaque-monkey-selfie-grapher-people-for-the-ethical-treatment-of-animals-funny-mammal-animals-thumbnail.png"></img>
    </div>
</template>

<style scoped></style>
