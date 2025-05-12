<script setup lang="ts">
import ServerList from './components/ServerList/ServerList.vue'
import SelectedMenu from './components/SelectedMenu.vue'

import type { ChannelDescription } from './types'

import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const channelsRef = ref<{
  textChannels: ChannelDescription[]
  voiceChannels: ChannelDescription[]
}>({ textChannels: [], voiceChannels: [] })

watch(route, async (newRoute, oldRoute) => {
  if (newRoute) {
    console.log(newRoute.params.channelId)
    const data = (await fetch(
      `http://localhost:3000/channelInfo/${newRoute.params.channelId}`,
    )) as any
    channelsRef.value = await data.json()
    console.log(channelsRef.value)
  }
})
</script>

<template>
  <header>
    <div class="logo">encontro</div>
  </header>
  <div class="main-container">
    <ServerList></ServerList>
    <SelectedMenu
      :voice-channels="channelsRef?.voiceChannels"
      :text-channels="channelsRef.textChannels"
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
