<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'

const currentMessageDraft = ref('')
const messages = reactive<Array<Message>>([
  {
    author: 'system',
    timestamp: new Date().toISOString(),
    message: 'This is a system message',
  },
])

type Message = {
  timestamp: string
  author: 'current-user' | 'other-user' | 'system'
  message: string
}

// @todo Allow user to change this
const maxHistoryDurationInMs = 120000

function isIsoDatetimeOlderThan(isoDatetime: string, olderThanTimeInSeconds: number): boolean {
  const inputTime = new Date(isoDatetime).getTime()
  if (isNaN(inputTime)) {
    throw new Error('Invalid ISO datetime string')
  }

  const currentTime = Date.now()
  return currentTime - inputTime > olderThanTimeInSeconds
}

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

async function startNewChat() {
  //
}

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
  <div class="flex w-full flex-row justify-between p-6 md:p-12">
    <div class="w-full">
      <div class="flex w-full flex-row justify-between pb-4 align-middle">
        <p class="text-2xl text-gray-500">VaporChat</p>
        <button
          class="cursor-pointer rounded-2xl bg-cyan-200 px-4 py-1 text-gray-600"
          @click="startNewChat"
        >
          Start new Chat
        </button>
      </div>

      <div class="flex flex-col gap-4 md:flex-row md:justify-center">
        <div class="w-full max-w-lg">
          <p class="pb-1 text-sm font-medium text-gray-700">Chat</p>
          <!-- @todo Make this fixed width but scrollable within -->
          <div
            class="flex min-h-8 flex-col gap-2 rounded-lg border border-gray-200 px-4 py-2 text-gray-700 shadow-sm"
          >
            <p v-if="messages.length === 0" class="text-center text-sm text-gray-500">
              ... no messages ...
            </p>
            <div
              v-for="message in messages"
              :key="message.timestamp"
              class="rounded-xl border border-gray-200 px-2 py-0.5 text-gray-700"
            >
              <pre class="whitespace-pre-wrap">{{ message.message }}</pre>
            </div>
          </div>
        </div>

        <div class="w-full max-w-lg">
          <label for="message" class="block pb-1 text-sm font-medium text-gray-700">
            Your Message
          </label>
          <textarea
            id="message"
            rows="4"
            class="w-full resize-y rounded-lg border border-gray-200 px-4 py-2 text-gray-700 placeholder-gray-400 shadow-sm outline-none"
            placeholder="Type your message here..."
            v-model="currentMessageDraft"
            @keydown.enter.exact.prevent="sendNewMessage"
          />
        </div>
      </div>
    </div>
  </div>
</template>
