<script setup lang="ts">
import Chat from './Chat.vue'

async function startNewChat() {
  //
}
</script>

<template>
  <div class="p-6 md:p-12">
    <div class="flex w-full flex-col">
      <div class="flex w-full flex-row justify-between pb-4 align-middle">
        <p class="text-2xl text-gray-500">VaporChat</p>
        <button
          class="cursor-pointer rounded-2xl bg-cyan-200 px-4 py-1 text-gray-600"
          @click="startNewChat"
        >
          Start new Chat
        </button>
      </div>
      <Chat />
    </div>
  </div>
</template>
