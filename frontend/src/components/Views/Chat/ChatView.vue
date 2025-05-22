<script lang="ts" setup>
import httpClient from '@/httpClient/httpClient'
import type { ChatInfo, MessagesByDate } from '@/types'
import MessageInput from '@/components/shared/MessageInput.vue'
import { onBeforeMount, onMounted, ref, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { groupByDate } from '@/helpers/groupByDate'
import { getFormatedTime } from '@/helpers/getFormatedTime'
import Messages from '@/components/shared/Messages/Messages.vue'

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
      <Messages :messages="serverInfoRef.messages"></Messages>
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

::-webkit-scrollbar {
  width: 6px;
  background-color: black;
}
::-webkit-scrollbar-thumb {
  background-color: white;
  width: 10px;
  border-radius: 5px;
}
.chat-wrapper {
  width: 100%;
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
</style>
