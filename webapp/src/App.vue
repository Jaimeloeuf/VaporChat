<script setup lang="ts">
import Chat from './Chat.vue'

import { ref, computed } from 'vue'
import { getWebsocketStateString } from './getWebsocketStateString'

const ws = ref<WebSocket | null>(null)

const wsConnectionState = ref<WebSocket['readyState'] | undefined>(undefined)
const isWebsocketConnected = computed(() => wsConnectionState.value === WebSocket.OPEN)

const websocketServerUrlBuilder = (chatID: string) =>
  `ws://localhost:3000/api/chat/join/${chatID}/websocket`

async function startNewChat() {
  const chatID = crypto.randomUUID()

  ws.value = new WebSocket(websocketServerUrlBuilder(chatID))
  wsConnectionState.value = ws.value.readyState

  ws.value.addEventListener('open', () => {
    wsConnectionState.value = ws.value?.readyState
  })
  ws.value.addEventListener('close', () => {
    wsConnectionState.value = ws.value?.readyState
  })
  ws.value.addEventListener('error', () => {
    wsConnectionState.value = ws.value?.readyState
  })
}
</script>

<template>
  <div class="p-6 md:p-12">
    <div class="flex w-full flex-col">
      <div v-if="ws === null" class="flex w-full flex-row justify-between pb-4 align-middle">
        <p class="text-2xl text-gray-500">VaporChat</p>
        <div>
          <button
            class="cursor-pointer rounded-2xl bg-cyan-200 px-4 py-2 text-gray-600"
            @click="startNewChat"
          >
            Start new Chat
          </button>
        </div>
      </div>
      <div v-else>
        <div class="flex w-full flex-row justify-between pb-4 align-middle">
          <p class="text-2xl text-gray-500">VaporChat</p>
          <div>
            <button
              class="cursor-pointer rounded-2xl bg-cyan-200 px-4 py-2 text-gray-600"
              @click="startNewChat"
            >
              Start new Chat
            </button>
          </div>
        </div>
        <div>
          <p class="text-sm font-light text-gray-500">
            {{ getWebsocketStateString(wsConnectionState) }}
          </p>
        </div>
        <Chat v-if="isWebsocketConnected" :ws="ws" />
      </div>
    </div>
  </div>
</template>
