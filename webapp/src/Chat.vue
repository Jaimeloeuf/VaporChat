<script setup lang="ts">
import type { Message } from './Message'

import { ref, reactive, onMounted } from 'vue'
import { isIsoDatetimeOlderThan } from './isIsoDatetimeOlderThan'

const props = defineProps<{
  ws: WebSocket
}>()

const currentMessageDraft = ref('')
const messages = reactive<Array<Message>>([
  {
    author: 'system',
    timestamp: new Date().toISOString(),
    message: 'This is a system message',
  },
])

props.ws?.addEventListener('message', function (event) {
  addNewLocalMessage({
    author: 'other-user',
    timestamp: new Date().toISOString(),
    message: event.data,
  })
})

// @todo Allow user to change this
const maxHistoryDurationInMs = 120000

/**
 * Check for and remove old messages every second
 */
onMounted(() => {
  setInterval(() => {
    while (
      messages.length > 0 &&
      messages[0]?.timestamp !== undefined &&
      isIsoDatetimeOlderThan(messages[0]?.timestamp, maxHistoryDurationInMs)
    ) {
      messages.shift()
    }
  }, 1000)
})

// @todo Allow user to change this
const maxMessagesLength = 20

/**
 * Adds new message to the local messages array and ensure that it is not too
 * long.
 */
function addNewLocalMessage(message: Message) {
  messages.push(message)

  while (messages.length > maxMessagesLength) {
    messages.shift()
  }

  // Clear text by char count.... but leave at least 2 last messages
}

async function sendNewMessage() {
  if (currentMessageDraft.value === '') {
    return
  }

  const message: Message = {
    author: 'current-user',
    timestamp: new Date().toISOString(),
    message: currentMessageDraft.value,
  }
  addNewLocalMessage(message)
  currentMessageDraft.value = ''
}
</script>

<template>
  <div class="text-gray-700">
    <div class="flex flex-col gap-4 md:flex-row md:justify-center">
      <div class="w-full max-w-lg">
        <p class="pb-1 text-sm font-medium">Chat</p>
        <!-- @todo Make this fixed width but scrollable within -->
        <div
          class="flex min-h-8 flex-col gap-2 rounded-lg border border-gray-200 px-4 py-2 shadow-sm"
        >
          <p v-if="messages.length === 0" class="text-center text-sm">... no messages ...</p>
          <div
            v-for="message in messages"
            :key="message.timestamp"
            class="rounded-xl border border-gray-200 px-2 py-0.5"
          >
            <pre class="whitespace-pre-wrap">{{ message.message }}</pre>
          </div>
        </div>
      </div>

      <div class="w-full max-w-lg">
        <label for="message" class="block pb-1 text-sm font-medium"> Your Message </label>
        <textarea
          id="message"
          rows="4"
          class="w-full resize-y rounded-lg border border-gray-200 px-4 py-2 placeholder-gray-400 shadow-sm outline-none"
          placeholder="Type your message here..."
          v-model="currentMessageDraft"
          @keydown.enter.exact.prevent="sendNewMessage"
        />
      </div>
    </div>
  </div>
</template>
