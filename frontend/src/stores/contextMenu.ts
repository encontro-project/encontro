import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useContextMenuStore = defineStore('contextMenu', () => {
  const isMenuActive = ref<boolean>(false)

  function hideMenu() {
    console.log('menu hided')
    isMenuActive.value = false
  }

  function openMenu() {
    console.log('menu opened')
    isMenuActive.value = true
  }

  return {
    isMenuActive,
    openMenu,
    hideMenu,
  }
})
