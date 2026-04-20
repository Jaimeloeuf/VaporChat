<script setup lang="ts">
import Chat from './Chat.vue'

import { ref, onMounted } from 'vue'
import { getWebsocketStateString } from './getWebsocketStateString'

const ws = ref<WebSocket | null>(null)

const wsConnectionState = ref<WebSocket['readyState'] | undefined>(undefined)

// @todo Allow user to customise
const websocketServerUrl = 'ws://localhost:3000/ws'

onMounted(() => {
  ws.value = new WebSocket(websocketServerUrl)
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
})

// onBeforeUnmount(() => {
//   if (socket) socket.close();
// });

async function startNewChat() {
  if (ws.value === null) {
    console.log('no webserver connection now')
    return
  }

  ws.value!.send('Hello Server!')

  // @todo This creates a new listener everytime...
  ws.value.addEventListener('message', function (event) {
    console.log('Server says:', event.data)
  })
}
</script>

<template>
  <div class="p-6 md:p-12">
    <div class="flex w-full flex-col">
      <div class="flex w-full flex-row justify-between pb-4 align-middle">
        <div class="flex flex-col">
          <p class="text-2xl text-gray-500">VaporChat</p>
          <p class="text-sm font-light text-gray-500">
            {{ getWebsocketStateString(wsConnectionState) }}
          </p>
        </div>
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
