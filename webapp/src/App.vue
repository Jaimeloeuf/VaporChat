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

const resetConfig = () => window.location.reload()

const leaveChat = () => window.location.reload()
</script>

<template>
  <div class="p-6 md:p-12">
    <div class="flex w-full flex-col">
      <div v-if="ws === null" class="flex h-[80dvh] w-full flex-col items-center justify-center">
        <div class="flex flex-col gap-8">
          <p class="text-2xl text-gray-500">VaporChat</p>
          <div class="w-xs">
            <div class="flex flex-row justify-between pb-1">
              <p class="text-gray-500">Config</p>
              <button
                class="cursor-pointer text-sm font-light text-red-400 underline"
                @click="resetConfig"
              >
                reset
              </button>
            </div>
            <div class="rounded-lg border border-gray-200 p-2">
              <p class="text-gray-500">Max messages to keep in chat</p>
              <p class="text-gray-500">Max message retention time in seconds</p>
            </div>
          </div>
          <div>
            <button
              class="cursor-pointer rounded-2xl border border-green-300 px-4 py-1 text-green-600"
              @click="startNewChat"
            >
              Start new Chat
            </button>
          </div>
        </div>
      </div>
      <div v-else>
        <div class="flex w-full flex-row justify-between pb-4 align-middle">
          <p class="text-2xl text-gray-500">VaporChat</p>
          <div>
            <button
              class="cursor-pointer rounded-2xl border border-red-500 px-4 py-1 text-red-500"
              @click="leaveChat"
            >
              leave
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
