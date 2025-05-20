<script lang="ts" setup>
import type { PeerConnection } from '@/types'
import { toRef, ref, watch } from 'vue'
interface Props {
  peerConnection: PeerConnection
}

const props = defineProps<Props>()

const audioStreamRef = toRef(props.peerConnection, 'microphoneStream') // реф для звуковой дорожки

const audioRef = ref<HTMLAudioElement | null>(null) // реф для элемента

watch([audioRef, audioStreamRef], ([newAudio, newStream], [oldAudio, oldStream]) => {
  if (oldStream) {
  }
  if (newStream && audioRef.value) {
    audioRef.value.srcObject = newStream
  }
})
</script>

<template>
  <audio ref="audioRef" autoplay v-if="props.peerConnection.microphoneStream"></audio>
</template>
