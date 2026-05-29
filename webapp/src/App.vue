<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterView } from 'vue-router'

import Logo from './Logo.vue'
import { useWebsocket } from './useWebsocket.ts'

const { setupWebsocket, isWebsocketConnected } = useWebsocket()

onMounted(() => {
  setupWebsocket()
})
</script>

<template>
  <div v-if="!isWebsocketConnected" class="flex h-dvh w-full flex-col items-center justify-center">
    <div class="text-gray-700 italic">
      <Logo />
    </div>
    <p class="text-2xl font-thin text-gray-400">... loading ...</p>
  </div>
  <div v-else class="mx-auto max-w-7xl p-4 text-gray-500">
    <RouterView />
  </div>
</template>
