<script setup lang="ts">
import ServerList from './components/ServerList/ServerList.vue'
import SelectedMenu from './components/SelectedMenu.vue'
import PeersVoiceTracks from './components/PeersVoiceTracks/PeersVoiceTracks.vue'

import type { ChannelDescription } from './types'

import { storeToRefs } from 'pinia'

import { useConnectionsStore } from '@/stores/rtcConnections'

import { useRoomWsStore } from '@/stores/wsConnection'

import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const rtcConnectionsStore = useConnectionsStore()

const roomWsConnectionStore = useRoomWsStore()

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
} = rtcConnectionsStore

const { closeRoomWsConnection } = roomWsConnectionStore

const { roomWs, localUuid, localDisplayName } = storeToRefs(roomWsConnectionStore)

const { peerConnections } = storeToRefs(rtcConnectionsStore)

const channelsRef = ref<{
  chats: ChannelDescription[]
  voiceChannels: ChannelDescription[]
}>({ chats: [], voiceChannels: [] })

watch(route, async (newRoute, oldRoute) => {
  if (newRoute.matched[0].name == 'channel') {
    //проверяем, нужный ли path для изменения данных в selectedMenuItem
    console.log(newRoute)
    const data = (await fetch(
      `http://localhost:3000/channel-info/${newRoute.params.channelId}`,
    )) as any
    channelsRef.value = await data.json()
    console.log(channelsRef.value)
  }
})

const gotMessageFromServer = async (message: MessageEvent) => {
  const signal = JSON.parse(message.data)
  const peerUuid = signal.uuid

  if (!roomWs.value) {
    return
  }

  /* console.log(signal, peerConnections.value, peerUuid) */
  if (signal.type === 'peer-disconnect') {
    handlePeerDisconnect(peerUuid)

    return
  }
  if (signal.type == 'get-tracks' && peerUuid != localUuid.value) {
    if (!peerConnections.value[peerUuid].ssVideoSender) {
      shareScreen(peerUuid)
    }
  }
  if (signal.type == 'get-microphone' && peerUuid != localUuid.value) {
    if (!peerConnections.value[peerUuid].microphoneStream) {
      shareMicrophone(peerUuid)
    }
  }
  if (signal.type == 'stop-stream' && peerUuid != localUuid.value) {
    if (peerConnections.value[peerUuid].ssVideoStream) {
      unsubscribeFromStream(peerUuid)
    }
  }
  if (peerUuid === localUuid.value || (signal.dest !== localUuid.value && signal.dest !== 'all'))
    return

  if (signal.metadata) {
    setTrackMetadata(peerUuid, signal.metadata)
  }
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
    setupPeer(peerUuid, signal.displayName, true)
  } else if (signal.sdp) {
    console.log('dest: ', signal.dest, 'uuid: ', peerUuid)
    await handleSdpSignal(signal, peerUuid)
  } else if (signal.ice) {
    handleIceCandidate(signal, peerUuid)
  }
}

watch(
  roomWs,
  (newSocket, oldSocket) => {
    if (oldSocket) {
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
  <header>
    <div class="logo">encontro</div>
  </header>
  <div class="main-container">
    <PeersVoiceTracks></PeersVoiceTracks>
    <ServerList></ServerList>
    <SelectedMenu
      :voice-channels="channelsRef?.voiceChannels"
      :chats="channelsRef.chats"
    ></SelectedMenu>
    <router-view> </router-view>
  </div>
</template>
<style>
@font-face {
  font-family: Roboto, system-ui;
  src: url('../public/Roboto-VariableFont_wdth/wght.ttf');
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
  background-color: #2a2a2a;
}
.logo {
  text-align: center;
  display: flex;
  align-items: center;
  margin-bottom: 7px;
}
body {
  overflow-x: hidden;
  margin: 0 auto;
  font-family: Roboto, system-ui;
}
</style>
