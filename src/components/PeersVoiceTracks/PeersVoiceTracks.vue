<script lang="ts" setup>
import { useConnectionsStore } from '@/stores/rtcConnections'
import { storeToRefs } from 'pinia'
import PeerVoiceTrack from './PeerVoiceTrack.vue'

const rtcConnectionsStore = useConnectionsStore()

const { peerConnections } = storeToRefs(rtcConnectionsStore)
</script>

<template>
  <PeerVoiceTrack
    :peer-connection="peerConnection"
    v-for="peerConnection in Object.values(peerConnections)"
  ></PeerVoiceTrack>
</template>
