<script setup lang="ts">
import type { ChatConfig } from './ChatConfig'
import type { Message } from './Message'

import { ref, reactive, onMounted, watch, nextTick } from 'vue'
import { isIsoDatetimeOlderThan } from './isIsoDatetimeOlderThan'
import { createChatUpdate } from './ChatUpdate'
import { useWebsocket } from './useWebsocket.ts'

const { getWebsocket, sendChatUpdateOverWebsocket } = useWebsocket()
const websocket = getWebsocket()

const props = defineProps<{
  chatConfig: Readonly<ChatConfig>
}>()

const messageContainer = ref<HTMLDivElement | null>(null)
const currentMessageDraft = ref('')
const messages = reactive<Array<Message>>([])

// Watch the messages array and scroll whenever it changes
watch(
  messages,
  async function () {
    if (messageContainer.value === null) {
      return
    }

    // Wait for Vue to finish rendering new message in the DOM before scrolling
    await nextTick()

    messageContainer.value.scrollTo({
      top: messageContainer.value.scrollHeight,
      behavior: 'smooth',
    })
  },
  {
    // Ensures it fires when new items are pushed into the array
    deep: true,
  },
)

websocket.addEventListener('message', function (event) {
  addNewLocalMessage({
    author: 'other-user',
    timestamp: new Date().toISOString(),
    message: event.data,
  })
})

/**
 * Check for and remove old messages every second
 */
onMounted(() => {
  const maxHistoryDurationInMs = props.chatConfig.maxHistoryDurationInSeconds * 1000

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

/**
 * Adds new message to the local messages array and ensure that it is not too
 * long.
 */
function addNewLocalMessage(message: Message) {
  messages.push(message)

  while (messages.length > props.chatConfig.maxMessagesLength) {
    messages.shift()
  }

  // Clear text by char count.... but leave at least 2 last messages
}

async function sendNewMessage() {
  if (currentMessageDraft.value === '') {
    return
  }

  const chatUpdate = createChatUpdate('message-new', {
    payload: {
      roomID: props.chatConfig.chatRoomTTL.toString(),
      message: currentMessageDraft.value,
    },
  })

  const message: Message = {
    author: 'current-user',
    timestamp: new Date().toISOString(),
    message: currentMessageDraft.value,
  }
  addNewLocalMessage(message)
  currentMessageDraft.value = ''

  sendChatUpdateOverWebsocket(chatUpdate)
}
</script>

<template>
  <div class="text-gray-500">
    <div class="flex flex-col gap-4 md:flex-row md:justify-center">
      <div class="basis-1/5">
        <p class="pb-1 text-sm font-medium">Chat Config</p>
        <div class="flex flex-col gap-4 rounded-lg border border-gray-200 px-2 py-4 shadow">
          <div>
            <p>Chat Room Time To Live (TTL)</p>
            <p class="text-xs">
              *Default to 300 seconds / 5 mins, after which chat room will be permanently destroyed
            </p>
            <input
              :value="props.chatConfig.chatRoomTTL"
              class="w-full rounded border border-gray-200 p-1.5 outline-none"
              disabled
            />
          </div>
          <div>
            <p>Max number of participants</p>
            <p class="text-xs">*Default is a 2 person peer to peer chat</p>
            <input
              :value="props.chatConfig.maxNumberOfParticipants"
              class="w-full rounded border border-gray-200 p-1.5 outline-none"
              disabled
            />
          </div>
          <div>
            <p>Max messages to keep in chat</p>
            <p class="text-xs">*Older messages will be auto deleted</p>
            <input
              :value="props.chatConfig.maxMessagesLength"
              class="w-full rounded border border-gray-200 p-1.5 outline-none"
              disabled
            />
          </div>
          <div>
            <p>Max message retention time in seconds</p>
            <p class="text-xs">*Expired messages will be auto deleted</p>
            <input
              :value="props.chatConfig.maxHistoryDurationInSeconds"
              class="w-full rounded border border-gray-200 p-1.5 outline-none"
              disabled
            />
          </div>
        </div>
      </div>

      <div class="basis-2/5">
        <p class="pb-1 text-sm font-medium">Chat</p>
        <div
          ref="messageContainer"
          class="no-scrollbar flex h-[70dvh] flex-col gap-2 overflow-y-scroll rounded-lg border border-gray-200 p-4 shadow-sm"
          :class="{
            'justify-center': messages.length === 0,
            'justify-end': messages.length !== 0,
          }"
        >
          <p v-if="messages.length === 0" class="text-center font-thin">... no messages ...</p>
          <div
            v-for="message in messages"
            :key="message.timestamp"
            class="rounded-xl border border-gray-200 px-2 py-0.5"
          >
            <pre class="whitespace-pre-wrap">{{ message.message }}</pre>
          </div>
        </div>
      </div>

      <div class="basis-2/5">
        <label for="message" class="block pb-1 text-sm font-medium">Your Message</label>
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

<style scoped>
.no-scrollbar {
  -ms-overflow-style: none; /* IE and Edge */
  scrollbar-width: none; /* Firefox */

  &::-webkit-scrollbar {
    display: none; /* Chrome, Safari, and Opera */
  }
}
</style>
