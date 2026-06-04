<script setup lang="ts">
import LogoWithConnectionStatus from './LogoWithConnectionStatus.vue'
import RightChevron from './RightChevron.vue'

import { RouterLink } from 'vue-router'
import { useAnonymousUser } from './useAnonymousUser.ts'

const { userID, username } = useAnonymousUser()
</script>

<template>
  <div class="flex flex-row justify-center">
    <div class="flex w-sm flex-col gap-6">
      <LogoWithConnectionStatus />
      <div class="py-4">
        <p class="text-lg">Temporary Anonymous Account</p>
        <p class="pb-1 text-xs">
          *Temporary anonymous account generated for you, this is deleted once you leave the chat
        </p>
        <div class="flex flex-col gap-4 rounded-lg border border-gray-200 p-4">
          <div>
            <p>User ID</p>
            <input
              v-model="userID"
              type="text"
              class="w-full rounded border border-gray-200 bg-gray-50 p-1.5 outline-none"
              disabled
            />
          </div>
          <div>
            <p>Username</p>
            <input
              v-model="username"
              type="text"
              class="w-full rounded border border-gray-200 p-1.5 outline-none"
            />
          </div>
        </div>
      </div>
      <RouterLink
        :to="{
          name: 'chat-start-new',
        }"
      >
        <div
          class="flex cursor-pointer flex-row items-center justify-between rounded-lg border border-gray-200 p-8 shadow hover:bg-green-100"
        >
          <p>Start new Chat</p>
          <RightChevron class="h-4 w-4" />
        </div>
      </RouterLink>
      <RouterLink
        :to="{
          name: 'chat-start-join',
        }"
      >
        <div
          class="flex cursor-pointer flex-row items-center justify-between rounded-lg border border-gray-200 p-8 shadow hover:bg-green-100"
        >
          <p>Join Chat</p>
          <RightChevron class="h-4 w-4" />
        </div>
      </RouterLink>
    </div>
  </div>
</template>
