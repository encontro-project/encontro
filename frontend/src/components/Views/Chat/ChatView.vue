<script lang="ts" setup>
import httpClient from '@/httpClient/httpClient'
import type { ChatInfo } from '@/types'
import MessageInput from '@/components/shared/MessageInput.vue'
import { onBeforeMount, onMounted, ref, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const serverInfoRef = ref<ChatInfo>({ title: '', messages: [] })

const messagesRef = ref<HTMLDivElement>()
onBeforeMount(async () => {
  const data = (await httpClient.get(
    `http://localhost:3000/chat-info/${route.params.channelViewId}`,
  )) as ChatInfo

  serverInfoRef.value = data
  console.log(serverInfoRef.value)
})

function scrollToBottom() {
  nextTick(() => {
    if (
      messagesRef.value &&
      messagesRef.value!.scrollHeight - messagesRef.value!.scrollTop < 1000
    ) {
      console.log(messagesRef.value!.scrollTop, messagesRef.value!.scrollHeight)
      messagesRef.value!.scrollTop = messagesRef.value!.scrollHeight
    }
  })
}

//Некст коммит это поменяется на уже готовые данные
watch(
  serverInfoRef,
  async () => {
    scrollToBottom()
  },
  { deep: true },
)

onMounted(() => {})
</script>

<template>
  <div class="chat-wrapper">
    <div class="chat-container" ref="messagesRef">
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
    <MessageInput post-url=""></MessageInput>
  </div>
</template>

<style scoped>
.chat-container {
  position: relative;
  width: 100%;
  overflow-y: auto;
  min-height: calc(100vh - 130px);
  max-height: calc(100vh - 130px);
}
.chat-wrapper {
  max-height: calc(100vh - 40px);
  width: 100%;
  height: calc(100vh - 40px);
  position: relative;
}
.chat-header {
  position: sticky;
  display: flex;
  top: 0;
  height: 56px;
  z-index: 1;
  background-color: black;
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
  margin: 0 auto;
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  width: 97%;
  align-self: center;
}
.chat-messages .chat-message {
  margin-top: 20px;
  white-space: pre-wrap;
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
