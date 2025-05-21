<script lang="ts" setup>
import { useConnectionsStore } from '@/stores/rtcConnections'
import { useRoomWsStore } from '@/stores/wsConnection'
import { storeToRefs } from 'pinia'
import HangPhoneIcon from './Icons/HangPhoneIcon.vue'

const rtcConnectionsStore = useConnectionsStore()

const roomWsConnectionStore = useRoomWsStore()

const { getMediaTracks, closeRoomWsConnection } = roomWsConnectionStore

const { leaveCall, updateStream, shareScreen, stopStream } = rtcConnectionsStore

const { roomWs, localUuid, isWsConnected } = storeToRefs(roomWsConnectionStore)
const { peerConnections } = storeToRefs(rtcConnectionsStore)

const getTracks = async () => {
  await getMediaTracks()
}

const handleStreamStart = async () => {
  await getTracks()
  Object.values(peerConnections.value).forEach((peerConnection) => {
    shareScreen(peerConnection.uuid)
  })
}

const handleStreamChange = async () => {
  await getTracks()
  updateStream()
}

const handleStreamStop = () => {
  stopStream()
  roomWs.value?.send(
    JSON.stringify({
      uuid: localUuid.value,
      type: 'stop-stream',
    }),
  )
}

const handleLeaveCall = () => {
  if (roomWs.value) {
    roomWs.value.send(
      JSON.stringify({
        type: 'peer-disconnect',
        uuid: localUuid.value,
      }),
    )
    leaveCall()
    closeRoomWsConnection()
  }
}
</script>

<template>
  <div class="voice-chat-menu-wrapper" v-if="isWsConnected">
    <div class="voice-chat-menu-container">
      <div class="row1">
        <div class="button-container">
          <HangPhoneIcon @click="handleLeaveCall"></HangPhoneIcon>
        </div>
      </div>
      <div class="row2">
        <button @click="handleStreamStart">share screen</button>
        <button @click="handleStreamChange">stream change</button>
        <button @click="handleStreamStop">stop stream</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.voice-chat-menu-wrapper {
  border-radius: 5px;
  background-color: white;
  width: 300px;
  position: absolute;
  display: flex;
  justify-content: center;
  left: 87px;
  bottom: 50px;
  color: white;
}
.voice-chat-menu-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 20px 20px 20px 20px;
}
.button-container {
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: black;
  width: 36px;
  height: 36px;
  border-radius: 5px;
}
.row1,
.row2 {
  display: flex;
}
.row1 {
  width: 100%;
  justify-content: flex-end;
}
</style>
