<script setup lang="ts">
import LogoWithConnectionStatus from './LogoWithConnectionStatus.vue'

import { useRouter } from 'vue-router'
import { useWebsocket } from './useWebsocket.ts'
import { createWsRequest } from './createWsRequest.ts'
import { useAnonymousUser } from './useAnonymousUser.ts'
import { useChatConfig } from './useChatConfig.ts'
const { chatConfig } = useChatConfig()

const router = useRouter()
const { sendWsRequest } = useWebsocket()

const { userID, readonlyUsername } = useAnonymousUser()

function startNewChat() {
  sendWsRequest(
    createWsRequest('room-create', {
      payload: {
        chatConfig,
      },
    }),
  )
  router.push({
    name: 'chat-room',
  })
}

const resetConfig = () => window.location.reload()
</script>

<template>
  <div class="flex flex-row justify-center">
    <div class="flex w-sm flex-col gap-6">
      <LogoWithConnectionStatus />
      <div class="py-4">
        <p class="font-light">Username: {{ readonlyUsername }}</p>
        <p class="text-sm font-extralight">ID: {{ userID }}</p>
      </div>
      <div>
        <div class="flex flex-row justify-between pb-1">
          <p class="text-lg font-light">New chat room config</p>
          <button class="cursor-pointer text-sm font-light text-red-400" @click="resetConfig">
            reset
          </button>
        </div>
        <div class="flex flex-col gap-4 rounded-lg border border-gray-200 px-2 py-4 shadow">
          <div>
            <p>Chat Room Time To Live (TTL)</p>
            <p class="text-xs">
              *Default to 300 seconds / 5 mins, after which chat room will be permanently destroyed
            </p>
            <input
              v-model="chatConfig.chatRoomTTL"
              type="number"
              step="1"
              min="1"
              max="86400"
              class="w-full rounded border border-gray-200 p-1.5 outline-none"
            />
          </div>
          <div>
            <p>Max number of participants</p>
            <p class="text-xs">*Default is a 2 person peer to peer chat</p>
            <input
              v-model="chatConfig.maxNumberOfParticipants"
              type="number"
              step="1"
              min="1"
              class="w-full rounded border border-gray-200 p-1.5 outline-none"
            />
          </div>
          <div>
            <p>Max messages to keep in chat</p>
            <p class="text-xs">*Older messages will be auto deleted</p>
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
            <p class="text-xs">*Expired messages will be auto deleted</p>
            <input
              v-model="chatConfig.maxHistoryDurationInSeconds"
              type="number"
              step="1"
              min="1"
              class="w-full rounded border border-gray-200 p-1.5 outline-none"
            />
          </div>
        </div>
        <div class="pt-6 text-right">
          <button
            class="cursor-pointer rounded-2xl border border-green-300 px-4 py-1 text-green-600"
            @click="startNewChat"
          >
            Start new Chat
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
