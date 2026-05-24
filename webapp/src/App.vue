<script setup lang="ts">
import Logo from './Logo.vue'
import ChatStart from './ChatStart.vue'

import { onMounted } from 'vue'
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
  <ChatStart v-else />
</template>
