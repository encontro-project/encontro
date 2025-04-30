<script setup lang="ts">
import Videos from './components/Videos.vue'
</script>

<template>
  <div class="videos-wrapper">
    <Videos></Videos>
  </div>
</template>

<style scoped>
.videos-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 20px;
}
</style>
