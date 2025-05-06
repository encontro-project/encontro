import { defineStore, storeToRefs } from 'pinia'
import { ref } from 'vue'
import { useConnectionsStore } from './rtcConnections'
import { v4 as uuidv4 } from 'uuid'

const generateUUID = () => {
  const uuid = uuidv4()
  localStorage.setItem('uuid', JSON.stringify(uuid))
  localStorage.setItem('displayName', JSON.stringify(uuid))
  return uuid
}

export const useRoomWsStore = defineStore('roomWsStore', () => {
  const isWsConnected = ref<boolean>(false)
  const roomWs = ref<WebSocket | null>(null)
  const roomUrl = ref<string>('')
  const localUuid = ref<string>(
    localStorage.getItem('uuid')
      ? JSON.parse(localStorage.getItem('uuid') as string)
      : generateUUID(),
  )

  const localDisplayName = ref<string>(
    localStorage.getItem('displayName')
      ? JSON.parse(localStorage.getItem('displayName') as string)
      : generateUUID(),
  )

  const initWebSocket = (room: string) => {
    roomWs.value = new WebSocket(`wss://localhost:8443/ws/${room}`)

    roomWs.value.onopen = () => {
      roomWs.value!.send(
        JSON.stringify({
          displayName: localDisplayName.value,
          uuid: localUuid.value,
          dest: 'all',
        }),
      )

      isWsConnected.value = true
    }

    roomWs.value.onclose = () => {
      isWsConnected.value = false
    }
  }

  const closeRoomWsConnection = () => {
    roomWs.value?.close()
    roomWs.value = null
  }
  return {
    roomWs,
    isWsConnected,
    localUuid,
    localDisplayName,
    initWebSocket,
    closeRoomWsConnection,
  }
})
