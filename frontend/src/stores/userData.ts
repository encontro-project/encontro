import httpClient from '@/httpClient/httpClient'
import type { Server } from '@/types'
import { defineStore } from 'pinia'
import { ref } from 'vue'

type ServersResponse = {
  servers: Server[]
}

export const useUserDataStore = defineStore('userData', () => {
  const userData = ref<ServersResponse>({ servers: [] })
  const isLoading = ref<boolean>(true)

  async function fetchUserData() {
    isLoading.value = true
    const data = (await httpClient.get('/api/user/100')) as ServersResponse
    userData.value = data
    isLoading.value = false
  }

  return { userData, isLoading, fetchUserData }
})
