<script lang="ts" setup>
import httpClient from '@/httpClient/httpClient'
import { ref, nextTick, onMounted } from 'vue'
import Textarea from 'primevue/textarea'

interface Props {
  postUrl: string
}

const props = defineProps<Props>()

const inputRef = ref<InstanceType<typeof Textarea>>()

const inputValue = ref<string>('')

// TODO ебануть логику запросов

async function submitMessage() {
  if (inputValue.value && inputValue.value.trim()) {
    console.log(inputValue.value)
    /* const response = await httpClient.post(props.postUrl)
    console.log(response) */
    resetInput()
  }
}

function resetInput() {
  inputValue.value = ''
}

async function handleEnterPress(e: KeyboardEvent) {
  if (!e.shiftKey) {
    e.preventDefault()
    await submitMessage()
  }
}
</script>

<template>
  <div class="message-input-wrapper">
    <div class="message-input-container">
      <Textarea
        rows="1"
        class="message-input"
        placeholder="смс-очка"
        ref="inputRef"
        @keypress.enter="handleEnterPress"
        v-model="inputValue"
        :auto-resize="true"
      >
      </Textarea>
    </div>
  </div>
</template>

<style scoped>
.message-input-wrapper {
  display: flex;
  justify-content: center;
  width: 100%;
  background-color: black;
  min-height: 40px;
  position: sticky;
  bottom: 20px;
}
.message-input-container {
  margin-top: 40px;
  width: 95%;
  z-index: 1;
  min-height: 24px;
  background-color: white;
  bottom: 25px;
  border-radius: 6px;
}
.message-input {
  line-height: 24px;
  font-family: Roboto, system-ui;
  resize: none;
  width: 90%;
  min-height: 24px;
  outline: none;
  border: none;
  overflow-y: hidden;
  padding-top: 10px;
  padding-bottom: 10px;
  border-radius: 5px;
  padding-left: 15px;
  font-size: 20px;
}
</style>
