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
      // @todo Maybe room ID should be in this link?
      path: '/chat/room',
      name: 'chat-room',
      component: () => import('../Chat.vue'),
    },
  ],
})

export default router
