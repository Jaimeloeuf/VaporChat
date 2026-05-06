<script setup lang="ts">
import type { ChatConfig } from './ChatConfig'

import Logo from './Logo.vue'
import Chat from './Chat.vue'

import { reactive, ref, computed } from 'vue'
import { getWebsocketStateString } from './getWebsocketStateString'

const chatConfig = reactive<ChatConfig>({
  maxHistoryDurationInSeconds: 120,
  maxMessagesLength: 20,
})

const joinChatID = ref('')

const ws = ref<WebSocket | null>(null)

const wsConnectionState = ref<WebSocket['readyState'] | undefined>(undefined)
const isWebsocketConnected = computed(() => wsConnectionState.value === WebSocket.OPEN)

const websocketServerUrlBuilder = (chatID: string) =>
  `ws://localhost:3000/api/chat/join/${chatID}/websocket`

async function joinChat() {
  ws.value = new WebSocket(websocketServerUrlBuilder(joinChatID.value))
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
  <div class="mx-auto max-w-7xl p-4 text-gray-500">
    <div class="flex w-full flex-col">
      <div v-if="ws === null" class="flex h-[80dvh] w-full flex-col items-center justify-center">
        <div class="flex flex-col gap-8">
          <Logo />
          <div class="flex w-xs flex-col gap-4 rounded-lg border border-gray-200 px-2 py-4 shadow">
            <div>
              <div class="flex flex-row items-center justify-between">
                <p>Chat ID</p>
                <button
                  class="cursor-pointer text-sm font-light text-red-400 underline"
                  @click="joinChatID = ''"
                >
                  reset
                </button>
              </div>
              <input
                v-model="joinChatID"
                type="text"
                class="w-full rounded border border-gray-200 p-1.5 outline-none"
              />
            </div>
            <button
              class="cursor-pointer rounded-2xl border border-green-400 px-4 py-1 text-green-600"
              @click="joinChat"
            >
              Join Chat
            </button>
          </div>
          <div class="w-xs">
            <div class="flex flex-row justify-between pb-1">
              <p>Chat Config</p>
              <button
                class="cursor-pointer text-sm font-light text-red-400 underline"
                @click="resetConfig"
              >
                reset
              </button>
            </div>
            <div class="flex flex-col gap-4 rounded-lg border border-gray-200 px-2 py-4 shadow">
              <div>
                <p>Max messages to keep in chat</p>
                <input
                  v-model="chatConfig.maxMessagesLength"
                  type="number"
                  step="1"
                  min="1"
                  class="w-full rounded border border-gray-200 p-1.5 outline-none"
                />
              </div>
              <div>
                <p>Max message retention time in seconds</p>
                <input
                  v-model="chatConfig.maxHistoryDurationInSeconds"
                  type="number"
                  step="1"
                  min="1"
                  class="w-full rounded border border-gray-200 p-1.5 outline-none"
                />
              </div>
            </div>
          </div>
          <div class="text-right">
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
        <div class="flex w-full flex-row items-center justify-between pb-4 align-middle">
          <div>
            <Logo />
            <p class="text-sm font-light">
              {{ getWebsocketStateString(wsConnectionState) }}
            </p>
          </div>
          <div>
            <button
              class="cursor-pointer rounded-2xl border border-red-500 px-4 py-1 text-red-500"
              @click="leaveChat"
            >
              leave
            </button>
          </div>
        </div>
        <Chat v-if="isWebsocketConnected" :ws="ws" :chat-config="Object.freeze(chatConfig)" />
      </div>
    </div>
  </div>
</template>
