<script lang="ts" setup>
import httpClient from '@/httpClient/httpClient'
import { ref, nextTick, onMounted } from 'vue'

interface Props {
  postUrl: string
}

const props = defineProps<Props>()

const inputRef = ref<HTMLTextAreaElement>()

onMounted(() => {
  inputRef.value!.value = ''
})

// TODO ебануть логику запросов

async function postMessage() {
  if (inputRef.value) {
    const now = new Date()
    console.log(inputRef.value?.value)
    /* const response = await httpClient.post(props.postUrl)
    console.log(response) */
    resetInput()
  }
}
function adjustHeight(action: 'add-line' | 'remove-line') {
  nextTick(() => {
    if (action == 'add-line') {
      inputRef.value!.style.height = inputRef.value!.scrollHeight + 'px'
    }
    if (action == 'remove-line') {
      inputRef.value!.style.height = parseInt(inputRef.value!.style.height) - 24 + 'px'
    }
  })
}

function resetInput() {
  inputRef.value!.value = ''
  inputRef.value!.rows = 1
}

async function handleEnterPress(e: KeyboardEvent) {
  if (e.shiftKey) {
    adjustHeight('add-line')
  } else {
    await postMessage()
  }
}

function handleBackspacePress() {
  const textarea = inputRef.value as HTMLTextAreaElement
  const cursorPos = textarea.selectionStart
  const text = textarea.value
  if (text[cursorPos - 1] == '\n') {
    adjustHeight('remove-line')
  }
}
</script>

<template>
  <div class="message-input-wrapper">
    <div class="message-input-container">
      <textarea
        rows="1"
        type="text"
        class="message-input"
        placeholder="смс-очка"
        ref="inputRef"
        @keypress.enter="handleEnterPress"
        @keydown.delete="handleBackspacePress"
      >
      </textarea>
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
