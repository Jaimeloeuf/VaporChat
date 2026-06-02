<script setup lang="ts">
import LogoWithConnectionStatus from './LogoWithConnectionStatus.vue'
import RightChevron from './RightChevron.vue'

import { ref } from 'vue'
import { RouterLink } from 'vue-router'
import { getRandomAnonymousAnimalName } from './getRandomAnimalName'

// @todo this should be an external state
const username = ref(getRandomAnonymousAnimalName())
</script>

<template>
  <div class="flex flex-row justify-center">
    <div class="flex max-w-sm flex-col gap-6">
      <LogoWithConnectionStatus />
      <div class="py-8">
        <p>Username</p>
        <p class="pb-1 text-xs">*Temporary username, this is deleted once you leave the chat</p>
        <input
          v-model="username"
          type="text"
          class="w-full rounded border border-gray-200 p-1.5 outline-none"
        />
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
