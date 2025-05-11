import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    redirect: '',
  },
  { path: '/channels/:channelId', name: 'channel', redirect: '' },
  {
    path: `/channels/:channelId/:channelViewId`,
    name: 'channelView',
    component: () => import('./components/Views/VoiceChannel/VoiceChannel.vue'),
  },
]

const router = createRouter({ history: createWebHistory(), routes: routes })

export default router
