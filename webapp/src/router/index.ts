import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'chat-start',
      component: () => import('../ChatStart.vue'),
    },
    {
      path: '/chat/start/new',
      name: 'chat-start-new',
      component: () => import('../ChatStartNew.vue'),
    },
    {
      path: '/chat/start/join',
      name: 'chat-start-join',
      component: () => import('../ChatStartJoin.vue'),
    },
    {
      // @todo Maybe room ID should be in this link?
      path: '/chat/room',
      name: 'chat-room',
      component: () => import('../Chat.vue'),
    },
  ],
})

export default router
