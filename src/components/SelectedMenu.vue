<script lang="ts" setup>
import type { ChannelDescription } from '@/types'
import { useRoomWsStore } from '@/stores/wsConnection'
import { storeToRefs } from 'pinia'
import { useRouter, useRoute } from 'vue-router'

interface Props {
  voiceChannels: ChannelDescription[]
  textChannels: ChannelDescription[]
}
defineProps<Props>()

const roomWsConnectionStore = useRoomWsStore()

const route = useRoute()

const router = useRouter()

const { initWebSocket, getMicrophoneTrack } = roomWsConnectionStore

const { currentRoomUrl } = storeToRefs(roomWsConnectionStore)

const handleConnectionStart = async (room: string) => {
  initWebSocket(room)
  await getMicrophoneTrack()
}
</script>

<template>
  <section class="menu-container">
    <article class="menu-controls">
      <div class="menu-controls-item-active">
        <svg
          viewBox="0 0 64 64"
          xmlns="http://www.w3.org/2000/svg"
          stroke-width="3"
          stroke="#000000"
          fill="none"
        >
          <g id="SVGRepo_bgCarrier" stroke-width="0"></g>
          <g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
          <g id="SVGRepo_iconCarrier">
            <circle cx="22.83" cy="22.57" r="7.51"></circle>
            <path d="M38,49.94a15.2,15.2,0,0,0-15.21-15.2h0a15.2,15.2,0,0,0-15.2,15.2Z"></path>
            <circle cx="44.13" cy="27.22" r="6.05"></circle>
            <path
              d="M42.4,49.94h14A12.24,12.24,0,0,0,44.13,37.7h0a12.21,12.21,0,0,0-5.75,1.43"
            ></path>
          </g>
        </svg>
        <span>Друзьяшки</span>
      </div>
      <div class="menu-controls-item">
        <svg viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg" fill="#DFDFDF" stroke="#DFDFDF">
          <g id="SVGRepo_bgCarrier" stroke-width="0"></g>
          <g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
          <g id="SVGRepo_iconCarrier">
            <g id="Layer_2" data-name="Layer 2">
              <g id="invisible_box" data-name="invisible box">
                <rect width="48" height="48" fill="none"></rect>
              </g>
              <g id="Layer_7" data-name="Layer 7">
                <g>
                  <path
                    d="M24.3,6A11.2,11.2,0,0,0,16,9.3a11,11,0,0,0-3.5,8.2,2.5,2.5,0,0,0,5,0,6.5,6.5,0,0,1,2-4.7A6.2,6.2,0,0,1,24.2,11a6.5,6.5,0,0,1,1,12.9,4.4,4.4,0,0,0-3.7,4.4v3.2a2.5,2.5,0,0,0,5,0V28.7a11.6,11.6,0,0,0,9-11.5A11.7,11.7,0,0,0,24.3,6Z"
                  ></path>
                  <circle cx="24" cy="39.5" r="2.5"></circle>
                </g>
              </g>
            </g>
          </g>
        </svg>
        <span>Че-нибудь еще</span>
      </div>
    </article>
    <!-- <article class="menu-options">
      <div class="menu-options-item">
        <h1>Смс-ки</h1>
        <div class="menu-options-list">
          <div class="menu-options-list-item">
            <img
              src="https://e7.pngegg.com/pngimages/719/959/png-clipart-celebes-crested-macaque-monkey-selfie-grapher-people-for-the-ethical-treatment-of-animals-funny-mammal-animals-thumbnail.png"
              alt="avatar"
            />
            <div class="list-item-username">Обезьянка</div>
          </div>
        </div>
      </div>
    </article> -->
    <!-- Это для смсок -->
    <article class="menu-options">
      <div class="menu-options-item">
        <h1>Текстовые каналы</h1>
        <div class="menu-options-list">
          <router-link
            class="menu-options-list-item"
            v-for="i in textChannels"
            :to="`/channels/${route.params.channelId}/${i.url}`"
          >
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <g id="SVGRepo_bgCarrier" stroke-width="0"></g>
              <g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
              <g id="SVGRepo_iconCarrier">
                <path
                  d="M10 4L7 20M17 4L14 20M5 8H20M4 16H19"
                  stroke="#6c6c6c"
                  stroke-width="2"
                  stroke-linecap="round"
                ></path>
              </g>
            </svg>
            <div class="list-item-username">{{ i.channelTitle }}</div>
          </router-link>
        </div>
      </div>
      <div class="menu-options-item">
        <h1>Голосовые каналы</h1>
        <div class="menu-options-list">
          <div
            class="menu-options-list-item"
            v-for="i in voiceChannels"
            :class="
              currentRoomUrl == i.url ? 'menu-options-list-item-active' : 'menu-options-list-item'
            "
            @click="
              currentRoomUrl != i.url
                ? handleConnectionStart(i.url)
                : (() => {
                    console.log('dfffsd')
                    router.push(`/channels/${route.params.channelId}/${i.url}`)
                  })()
            "
          >
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <g id="SVGRepo_bgCarrier" stroke-width="0"></g>
              <g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
              <g id="SVGRepo_iconCarrier">
                <path
                  d="M13 3.7446C13 3.27314 12.8728 2.50021 12.1657 2.14424C11.4151 1.76635 10.7163 2.19354 10.3623 2.51158L4.94661 7.43717H3C1.89543 7.43717 1 8.3326 1 9.43717L1.00001 14.6248C1.00001 15.7293 1.89544 16.6248 3.00001 16.6248H4.95001L10.3623 21.4891C10.7175 21.8081 11.416 22.2331 12.1656 21.8554C12.8717 21.4998 13 20.7286 13 20.2561V3.7446Z"
                  fill="#a1a1a1"
                ></path>
                <path
                  d="M17.336 3.79605L17.0952 3.72886C16.5633 3.58042 16.0117 3.89132 15.8632 4.42329L15.7289 4.90489C15.5804 5.43685 15.8913 5.98843 16.4233 6.13687L16.6641 6.20406C18.9551 6.84336 20.7501 9.14615 20.7501 12.0001C20.7501 14.854 18.9551 17.1568 16.6641 17.7961L16.4233 17.8632C15.8913 18.0117 15.5804 18.5633 15.7289 19.0952L15.8632 19.5768C16.0117 20.1088 16.5633 20.4197 17.0952 20.2713L17.336 20.2041C20.7957 19.2387 23.2501 15.8818 23.2501 12.0001C23.2501 8.11832 20.7957 4.76146 17.336 3.79605Z"
                  fill="#a1a1a1"
                ></path>
                <path
                  d="M16.3581 7.80239L16.1185 7.73078C15.5894 7.57258 15.0322 7.87329 14.874 8.40243L14.7308 8.88148C14.5726 9.41062 14.8733 9.96782 15.4024 10.126L15.642 10.1976C16.1752 10.3571 16.75 11.012 16.75 12C16.75 12.9881 16.1752 13.643 15.642 13.8024L15.4024 13.874C14.8733 14.0322 14.5726 14.5894 14.7308 15.1185L14.874 15.5976C15.0322 16.1267 15.5894 16.4274 16.1185 16.2692L16.3581 16.1976C18.1251 15.6693 19.25 13.8987 19.25 12C19.25 10.1014 18.1251 8.33068 16.3581 7.80239Z"
                  fill="#a1a1a1"
                ></path>
              </g>
            </svg>
            <div class="list-item-username">{{ i.channelTitle }}</div>
          </div>
        </div>
      </div>
    </article>
  </section>
</template>

<style scoped>
.menu-container {
  display: flex;
  flex-direction: column;
  width: 328px;
  align-items: center;
  border-right: 1px solid rgba(255, 255, 255, 0.3);
  cursor: default;
}
.menu-controls {
  display: flex;
  flex-direction: column;
  width: 328px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.3);
}

.menu-controls-item-active {
  cursor: pointer;
  margin-left: 16px;
  width: 296px;
  display: flex;
  gap: 5px;
  font-size: 19px;
  color: white;
  align-items: center;
  height: 36px;
  margin-top: 10px;
  margin-bottom: 10px;
  background-color: #6c6c6c;
  border-radius: 5px;
}

.menu-controls-item-active span {
  margin-bottom: 4px;
}
.menu-controls-item-active svg {
  margin-left: 8px;
  stroke: white;
  width: 28px;
  height: 28px;
}

.menu-options-item svg {
  margin-left: 8px;
  /*   stroke: white; */
  width: 28px;
  height: 28px;
}

a {
  text-decoration: none;
}

.menu-controls-item {
  cursor: pointer;
  margin-left: 16px;
  width: 296px;
  display: flex;
  gap: 5px;
  font-size: 19px;
  color: #dfdfdf;
  align-items: center;
  height: 36px;
  margin-top: 10px;
  margin-bottom: 10px;
}
.menu-controls-item span {
  margin-bottom: 4px;
}
.menu-controls-item svg {
  margin-left: 8px;
  stroke: #dfdfdf;
  width: 28px;
  height: 28px;
}
.menu-options {
  width: 328px;
  display: flex;
  flex-direction: column;
  gap: 15px;
  align-items: center;
}

.menu-options-item {
  display: flex;
  flex-direction: column;
  margin-top: 10px;
  width: 296px;
  text-align: start;
  align-items: start;
}
.menu-options-item h1 {
  margin-bottom: 0;
  margin-top: 0;
  color: white;
  text-align: start;
  font-size: 19px;
  font-weight: 500;
}
.menu-options-list {
  margin-top: 10px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.menu-options-list-item,
.menu-options-list-item-active {
  width: 296px;
  display: flex;
  gap: 8px;
  color: white;
  align-items: center;
  height: 45px;
  border-radius: 5px;
  cursor: pointer;
  font-size: 19px;
}

.menu-options-list-item {
  background-color: #2a2a2a;
}

.menu-options-list-item:hover {
  background-color: #3d3d3d;
}

.menu-options-list-item-active {
  background-color: #6c6c6c !important;
}

.list-item-username {
  margin-bottom: 3px;
}
.menu-options-list-item img {
  margin-left: 10px;
  width: 35px;
  height: 35px;
  border-radius: 50%;
}
</style>
