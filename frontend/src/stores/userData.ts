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
    //Олег, можно трогать
    isLoading.value = true
    const data = (await httpClient.get('/user/100')) as ServersResponse
    userData.value = data
    console.log(data)
    isLoading.value = false
  }

  return { userData, isLoading, fetchUserData }
})
