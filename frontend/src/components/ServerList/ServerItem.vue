<script lang="ts" setup>
import { onMounted } from 'vue'

interface Props {
  serverTitle: string
  serverUrl: string
  serverImg?: string
}
const props = defineProps<Props>()
</script>

<template>
  <router-link :to="`/channels/${props.serverUrl}`">
    <div class="server-item">
      <h3>{{ props.serverTitle }}</h3>
    </div>
  </router-link>
</template>

<style scoped>
.server-item {
  cursor: pointer;
  margin-top: 15px;
  width: 40px;
  height: 40px;
  background-color: white;
  border-radius: 5px;
  display: flex;
  align-items: center;
  justify-content: center;
}
a {
  text-decoration: none;
  color: black;
  margin-bottom: 9px;
}
</style>
