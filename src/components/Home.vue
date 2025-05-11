<script lang="ts" setup>
import VoiceChannel from './Views/VoiceChannel/VoiceChannel.vue'
</script>

<template>
  <div class="videos-wrapper"><VoiceChannel></VoiceChannel></div>
</template>

<style scoped>
.videos-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  background-color: #2a2a2a;
  gap: 20px;
}
</style>
