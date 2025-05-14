<script lang="ts" setup>
import type { ChatInfo } from '@/types'
import { onBeforeMount, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

const serverInfoRef = ref<ChatInfo>({ title: '', messages: [] })

const route = useRoute()

onBeforeMount(async () => {
  const data = await fetch(`http://localhost:3000/chat-info/${route.params.channelViewId}`)
  serverInfoRef.value = await data.json()
  console.log(serverInfoRef.value)
})

onMounted(() => {})
</script>

<template>
  <div class="chat-wrapper">
    <div class="chat-header">
      <div class="chat-title">
        <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <g id="SVGRepo_bgCarrier" stroke-width="0"></g>
          <g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
          <g id="SVGRepo_iconCarrier">
            <path
              d="M10 4L7 20M17 4L14 20M5 8H20M4 16H19"
              stroke="#6c6c6c"
              stroke-width="2"
              stroke-linecap="round"
            ></path>
          </g>
        </svg>
        <h1>{{ serverInfoRef.title }}</h1>
      </div>
    </div>
    <div class="chat-messages">
      <div
        class="chat-message"
        v-for="message in serverInfoRef.messages.sort((a, b) => {
          return a.timestamp - b.timestamp
        })"
      >
        <img
          src="https://e7.pngegg.com/pngimages/719/959/png-clipart-celebes-crested-macaque-monkey-selfie-grapher-people-for-the-ethical-treatment-of-animals-funny-mammal-animals-thumbnail.png"
          alt=""
        />
        <div class="message-container">
          <p class="message-timestamp">{{ new Date(message.timestamp) }}</p>
          <p>
            {{ message.text }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.chat-wrapper {
  position: relative;
  display: flex;
  flex-direction: column;
  width: 100%;
  overflow-y: scroll;
}
.chat-header {
  position: fixed;
  display: flex;
  height: 56px;
  align-items: center;
  width: 100%;
  border-bottom: 1px rgba(255, 255, 255, 0.3) solid;
}
.chat-header svg {
  fill: white;
  width: 28px;
  stroke: white;
}
.chat-header svg g path {
  stroke: white;
}
.chat-title {
  margin-left: 25px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.chat-title h1 {
  font-weight: 400;
  margin: 0;
  margin-bottom: 3px;
  font-size: 20px;
  color: white;
}
.chat-messages {
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  width: 97%;
  height: 100%;
  align-self: center;
}
.chat-messages .chat-message {
  margin-top: 20px;
  display: flex;
  gap: 10px;
  color: white;
  font-size: 20px;
  border-bottom: 1px rgba(255, 255, 255, 0.3) solid;
}

.chat-message img {
  width: 48px;
  height: 48px;
  border-radius: 50%;
}
.message-timestamp {
  margin-top: 0;
  color: #6c6c6c;
}
</style>
