<script lang="ts" setup>
import ServerItem from './ServerItem.vue'

const servers = [
  { serverTitle: 'test', serverUrl: '1' },
  { serverTitle: 'test2', serverUrl: '2' },
]
</script>

<template>
  <div class="servers-list">
    <ServerItem
      v-for="i in servers"
      :server-title="i.serverTitle"
      :server-url="i.serverUrl"
    ></ServerItem>
  </div>
</template>

<style scoped>
.servers-list {
  display: flex;
  flex-direction: column;
  align-items: center;
  background-color: black;
  border-right: rgba(255, 255, 255, 0.3) 1px solid;
  min-width: 72px;
  min-height: 100vh;
}
</style>
