<script lang="ts" setup>
import { useConnectionsStore } from '@/stores/rtcConnections'
import { useRoomWsStore } from '@/stores/wsConnection'
import { storeToRefs } from 'pinia'
import HangPhoneIcon from './Icons/HangPhoneIcon.vue'
import MicrophoneIcon from '@/components/Icons/MicrophoneIcon.vue'
import MicrophoneOffIcon from '@/components/Icons/MicrophoneOffIcon.vue'

const rtcConnectionsStore = useConnectionsStore()

const roomWsConnectionStore = useRoomWsStore()

const { getMediaTracks, closeRoomWsConnection } = roomWsConnectionStore

const { leaveCall, updateStream, shareScreen, stopStream, toggleMicrophone } = rtcConnectionsStore
const { roomWs, localUuid, isWsConnected, currentRoom } = storeToRefs(roomWsConnectionStore)
const { peerConnections, isMicrophoneOn } = storeToRefs(rtcConnectionsStore)

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
    closeRoomWsConnection()
    leaveCall()
  }
}
</script>

<template>
  <div class="voice-chat-menu-wrapper" v-if="isWsConnected != 'disconnected'">
    <div class="voice-chat-menu-container">
      <div class="row1">
        <div>
          <p
            :class="`connection-state${isWsConnected == 'connecting' ? '-connecting' : '-connected'}`"
          >
            {{ isWsConnected == 'connecting' ? 'Connecting...' : 'Voice Connected' }}
          </p>
          <router-link
            :to="`/channels/${currentRoom.serverId}/voice-channel/${currentRoom.roomId}`"
            class="current-room-title"
            >{{ currentRoom.roomTitle }}</router-link
          >
        </div>
        <div class="row1-controls">
          <div class="button-container">
            <MicrophoneIcon @click="toggleMicrophone" v-if="isMicrophoneOn"></MicrophoneIcon>
            <MicrophoneOffIcon @click="toggleMicrophone" v-else></MicrophoneOffIcon>
          </div>
          <div class="button-container" title="Выйти из звонка">
            <HangPhoneIcon @click="handleLeaveCall"></HangPhoneIcon>
          </div>
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
  bottom: 58px;
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
  width: 30px;
  height: 30px;
  border-radius: 5px;
}
.row1,
.row2 {
  display: flex;
}

.row1-controls {
  display: flex;
  gap: 10px;
}

.row2 button {
  outline: none;
  border: 1px solid white;
  background-color: black;
  color: white;
  border-radius: 5px;
  font-size: 15px;
  cursor: pointer;
  padding: 5px 5px 5px 5px;
}

.row1 {
  width: 100%;
  justify-content: space-between;
}
.row1 p {
  margin-top: 0;
  margin-bottom: 0;
  font-weight: 500;
}
.row1 .button-container {
  justify-self: flex-end;
}
.connection-state-connected {
  color: darkgreen;
}
.connection-state-connecting {
  color: chocolate;
}
.current-room-title {
  text-decoration: none;
  color: black;
}
.current-room-title:hover {
  text-decoration: underline 1px black solid;
}
</style>
