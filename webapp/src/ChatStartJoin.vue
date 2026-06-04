<script setup lang="ts">
import LogoWithConnectionStatus from './LogoWithConnectionStatus.vue'

import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useWebsocket } from './useWebsocket.ts'
import { createChatUpdate } from './createChatUpdate.ts'
import { useAnonymousUser } from './useAnonymousUser.ts'

const router = useRouter()
const { sendChatUpdateOverWebsocket } = useWebsocket()

const { userID, readonlyUsername } = useAnonymousUser()
const joinChatID = ref('')

function joinChat() {
  sendChatUpdateOverWebsocket(
    createChatUpdate('status-join-room', {
      payload: {
        roomID: joinChatID.value,
      },
    }),
  )
  router.push({
    name: 'chat-room',
  })
}
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
        <p class="pb-1 text-2xl font-light">Join chat room</p>
        <div class="flex flex-col gap-4 rounded-lg border border-gray-200 px-2 py-4 shadow">
          <div>
            <div class="flex flex-row items-center justify-between">
              <p>Chat ID</p>
              <button
                class="cursor-pointer text-sm font-light text-red-400"
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
        </div>
        <div class="pt-6 text-right">
          <button
            class="cursor-pointer rounded-2xl border border-green-400 px-4 py-1 text-green-600"
            @click="joinChat"
          >
            Join Chat
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
