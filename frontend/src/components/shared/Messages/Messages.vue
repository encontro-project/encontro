<script lang="ts" setup>
import { groupByDate } from '@/helpers/groupByDate'
import { getFormatedDate } from '@/helpers/getFormatedDate'
import Message from './Message.vue'
import { toRef, watch } from 'vue'

import type { ChatMessage, MessagesByDate } from '@/types'
import { ref } from 'vue'

interface Props {
  messages: ChatMessage[]
}

const props = defineProps<Props>()

const groupedMessages = ref<MessagesByDate>(groupByDate(props.messages))

const messagesRef = toRef(props)

watch(
  messagesRef,
  () => {
    groupedMessages.value = groupByDate(props.messages)
  },
  { deep: true },
)
</script>

<template>
  <div class="chat-messages">
    <div
      class="date-messages"
      v-for="dateGroup in Object.values(groupedMessages).sort((a, b) => {
        return a.dateKey > b.dateKey ? -1 : 1
      })"
    >
      <div class="messages-date">
        <div class="date-border"></div>
        <p>{{ dateGroup.dateKey }}</p>
        <div class="date-border"></div>
      </div>
      <Message
        v-for="message in dateGroup.messages.sort((a, b) => {
          return a.timestamp < b.timestamp ? -1 : 1
        })"
        :sender_id="message.sender"
        :content="message.text"
        :timestamp="message.timestamp"
      >
      </Message>
    </div>
  </div>
</template>

<style scoped>
.chat-messages {
  margin: 0 auto;
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  width: 97%;
  align-self: center;
}
.messages-date {
  color: white;
  font-size: 16px;
  height: 40px;
  display: flex;
  align-items: center;
  width: 100%;
  gap: 10px;
  text-align: center;
}
.messages-date p {
  display: inline-block;
  white-space: nowrap;
  flex: 0 0 auto;
  margin-top: 0;
  margin-bottom: 2px;
}
.date-border {
  flex: 1;
  background-color: white;
  height: 1px;
  width: calc(50% - 40px);
}
</style>
