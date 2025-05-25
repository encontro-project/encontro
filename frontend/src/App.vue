<script setup lang="ts">
import ServerList from './components/ServerList/ServerList.vue'
import SelectedMenu from './components/SelectedMenu.vue'
import PeersVoiceTracks from './components/global/PeersVoiceTracks/PeersVoiceTracks.vue'

import { storeToRefs } from 'pinia'

import { useConnectionsStore } from '@/stores/rtcConnections'

import { useRoomWsStore } from '@/stores/wsConnection'

import { onBeforeMount, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

import { useUserDataStore } from './stores/userData'
import VoiceChatMenu from './components/global/VoiceChatMenu/VoiceChatMenu.vue'
import { useContextMenuStore } from './stores/contextMenu'

const rtcConnectionsStore = useConnectionsStore()

const userDataStore = useUserDataStore()

const roomWsConnectionStore = useRoomWsStore()

const contextMenuStore = useContextMenuStore()

const {
  setupPeer,
  handlePeerDisconnect,
  handleIceCandidate,
  handleSdpSignal,
  shareScreen,
  unsubscribeFromStream,
  leaveCall,
  shareMicrophone,
  setTrackMetadata,
  handleUserMute,
  handleUserUnmute,
} = rtcConnectionsStore

const { isMenuActive } = storeToRefs(contextMenuStore)

const { hideMenu } = contextMenuStore

const { fetchUserData } = userDataStore

const { isLoading, userData } = storeToRefs(userDataStore)

const { roomWs, localUuid, localDisplayName } = storeToRefs(roomWsConnectionStore)

const { peerConnections } = storeToRefs(rtcConnectionsStore)

const gotMessageFromServer = async (message: MessageEvent) => {
  const signal = JSON.parse(message.data)
  const peerUuid = signal.uuid

  console.log(signal)

  if (!roomWs.value) {
    return
  }

  //Убираем пир из peerConnections на его сообщение о дисконнекте
  if (signal.type === 'peer-disconnect') {
    handlePeerDisconnect(peerUuid)
    return
  }
  //Присылаем пиру tracks из демки на его сообщение
  if (signal.type == 'get-tracks' && peerUuid != localUuid.value) {
    if (!peerConnections.value[peerUuid].ssVideoSender) {
      shareScreen(peerUuid)
    }
    return
  }
  if (signal.type == 'user-muted' && signal.dest == 'all' && peerUuid) {
    handleUserMute(peerUuid)
    return
  }
  console.log(signal)
  if (signal.type == 'user-unmuted' && signal.dest == 'all' && peerUuid) {
    handleUserUnmute(peerUuid)
    return
  }
  //Присылаем пиру track микрофона на его сообщение
  if (signal.type == 'get-microphone' && peerUuid != localUuid.value) {
    if (!peerConnections.value[peerUuid].microphoneStream) {
      shareMicrophone(peerUuid)
    }
  }
  //Прекращаем стрим пира на его сообщение
  if (signal.type == 'stop-stream' && peerUuid != localUuid.value) {
    if (peerConnections.value[peerUuid].ssVideoStream) {
      unsubscribeFromStream(peerUuid)
    }
    return
  }
  //Ретурним, чтобы не вмешиваться в логику webRTCNegotiation
  if (peerUuid === localUuid.value || (signal.dest !== localUuid.value && signal.dest !== 'all'))
    return

  //Бродкастим сигнал чтобы подключиться ко всем пирам
  if (signal.displayName && signal.dest === 'all') {
    setupPeer(peerUuid, signal.displayName, false)
    roomWs.value?.send(
      JSON.stringify({
        displayName: localDisplayName.value,
        uuid: localUuid.value,
        dest: peerUuid,
      }),
    )
  } else if (signal.displayName && signal.dest === localUuid.value) {
    //Создаем оффер на пришедший сигнал
    setupPeer(peerUuid, signal.displayName, true)
  } else if (signal.sdp) {
    //Обмениваемся sdp описаниями
    console.log('dest: ', signal.dest, 'uuid: ', peerUuid)
    await handleSdpSignal(signal, peerUuid)
  } else if (signal.ice) {
    //Добавляем ICE кандидатов
    handleIceCandidate(signal, peerUuid)
  }
}

onMounted(async () => {
  await fetchUserData()
})

onUnmounted(() => {
  roomWs.value!.send(
    JSON.stringify({
      uuid: localUuid.value,
      type: 'peer-disconnect',
    }),
  )
})

watch(
  roomWs,
  (newSocket, oldSocket) => {
    if (oldSocket) {
      oldSocket.send(
        JSON.stringify({
          type: 'peer-disconnect',
          uuid: localUuid.value,
        }),
      )
      leaveCall()
      oldSocket.close()
      console.log('cleanup')
    }
    if (newSocket) {
      newSocket.onmessage = (e) => {
        gotMessageFromServer(e)
      }
      console.log(newSocket)
    }
  },
  { deep: false },
)
</script>

<template>
  <div class="loading-screen" v-if="isLoading">E</div>
  <div @click="isMenuActive ? hideMenu() : ''">
    <header>
      <div class="logo">encontro</div>
    </header>
    <div class="main-container">
      <PeersVoiceTracks></PeersVoiceTracks>
      <ServerList></ServerList>
      <SelectedMenu></SelectedMenu>
      <VoiceChatMenu></VoiceChatMenu>
      <router-view> </router-view>
    </div>
  </div>
</template>
<style>
@font-face {
  font-family: Roboto, system-ui;
  src: url('../public/Roboto-VariableFont_wdth/wght.ttf');
}

.loading-screen {
  width: 100vw;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 128px;
  background-color: black;
  color: white;
}
header {
  position: fixed;
  z-index: 999;
  top: 0;
  width: 100vw;
  height: 40px;
  display: flex;
  color: white;
  font-weight: 600;
  font-size: 24px;
  justify-content: center;
  align-items: center;
  background-color: black;
}
.main-container {
  display: flex;
  margin-top: 40px;
  width: 100vw;
  min-height: 100vh;
  background-color: black;
  position: relative;
}
.logo {
  text-align: center;
  display: flex;
  align-items: center;
  margin-bottom: 7px;
}
body {
  overflow: hidden;
  margin: 0 auto;
  font-family: Roboto, system-ui;
}
</style>
