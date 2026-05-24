<script setup lang="ts">
import type { ChatConfig } from './ChatConfig'

import Logo from './Logo.vue'
import LogoWithConnectionStatus from './LogoWithConnectionStatus.vue'
import Chat from './Chat.vue'

import { reactive, ref } from 'vue'
import { getRandomAnimalName } from './getRandomAnimalName'
import { useWebsocket } from './useWebsocket.ts'

const chatConfig = reactive<ChatConfig>({
  chatRoomTTL: 300,
  maxNumberOfParticipants: 2,
  maxHistoryDurationInSeconds: 120,
  maxMessagesLength: 20,
})

const joinChatID = ref('')
const username = ref(`Anonymous ${getRandomAnimalName()}`)

const { setupWebsocket, websocket, isWebsocketConnected } = useWebsocket()

function joinChat() {
  setupWebsocket()
}

async function startNewChat() {
  setupWebsocket()
}

const resetConfig = () => window.location.reload()

const leaveChat = () => window.location.reload()
</script>

<template>
  <div>
    <div>
      <LogoWithConnectionStatus />
      <div class="flex flex-col justify-center gap-8 align-middle md:flex-row">
        <div class="h-full w-full max-w-sm">
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
          <div class="pt-4 text-right">
            <button
              class="cursor-pointer rounded-2xl border border-green-400 px-4 py-1 text-green-600"
              @click="joinChat"
            >
              Join Chat
            </button>
          </div>
        </div>

        <div class="h-full w-full max-w-sm">
          <p class="pb-1 text-2xl font-light">New chat room</p>
          <div class="flex flex-col gap-4 rounded-lg border border-gray-200 px-2 py-4 shadow">
            <div class="flex flex-row justify-between border-b border-gray-200 pb-1">
              <p>Chat Config</p>
              <button class="cursor-pointer text-sm font-light text-red-400" @click="resetConfig">
                reset
              </button>
            </div>
            <div>
              <p>Username</p>
              <p class="text-xs">
                *Temporary username that will be deleted once you leave the chat room
              </p>
              <input
                v-model="username"
                type="text"
                class="w-full rounded border border-gray-200 p-1.5 outline-none"
              />
            </div>
            <div>
              <p>Chat Room Time To Live (TTL)</p>
              <p class="text-xs">
                *Default to 300 seconds / 5 mins, after which chat room will be permanently
                destroyed
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
          <div class="pt-4 text-right">
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

    <div>
      <div class="flex w-full flex-row items-center justify-between pb-4 align-middle">
        <LogoWithConnectionStatus />
        <div>
          <button
            class="cursor-pointer rounded-2xl border border-red-500 px-4 py-1 text-red-500"
            @click="leaveChat"
          >
            leave
          </button>
        </div>
      </div>
      <Chat :chat-config="Object.freeze(chatConfig)" />
    </div>
  </div>
</template>
