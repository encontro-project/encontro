import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: { render: () => null, functional: true },
  },
  {
    path: '/channels/:channelId',
    name: 'channel',
    component: { render: () => null, functional: true },
  },
  {
    path: `/channels/:channelId/chat/:channelViewId`,
    name: 'chatView',
    component: () => import('./components/Views/Chat/ChatView.vue'),
  },
  {
    path: `/channels/:channelId/voice-channel/:channelViewId`,
    name: 'voiceChannelView',
    component: () => import('./components/Views/VoiceChannel/VoiceChannel.vue'),
  },
]

const router = createRouter({ history: createWebHistory(), routes: routes })

export default router
