<script lang="ts" setup>
import { getFormatedTime } from '@/helpers/getFormatedTime'

interface Props {
  sender_id: string
  content: string
  timestamp: string
}
defineProps<Props>()
</script>

<template>
  <div class="chat-message">
    <img src="https://masterpiecer-images.s3.yandex.net/5fcb1cda5223d2d:upscaled" alt="" />
    <div class="message-container">
      <div class="message-header">
        <p class="message-username">Обезьянка</p>
        <p class="message-timestamp">
          {{ getFormatedTime(timestamp) }}
        </p>
      </div>
      <p>
        {{ content }}
      </p>
    </div>
  </div>
</template>

<style scoped>
.chat-message {
  margin-top: 20px;
  white-space: pre-wrap;
  display: flex;
  gap: 15px;
  color: white;

  font-size: 20px;
}

.chat-message img {
  width: 48px;
  height: 48px;
  border-radius: 50%;
}

.message-timestamp {
  color: #6c6c6c;
}
.message-header {
  display: flex;
  gap: 10px;
}
.message-header p {
  margin-top: 0;
  margin-bottom: 0;
}
.message-container p {
  margin-top: 10px;
}
</style>
