import { ref, computed } from 'vue'
import { getWebsocketStateString } from './getWebsocketStateString'

const websocket = ref<WebSocket | null>(null)
const websocketConnectionState = ref<WebSocket['readyState'] | undefined>(undefined)
const websocketConnectionStateString = computed(() =>
  getWebsocketStateString(websocketConnectionState.value),
)

const isWebsocketConnected = computed(() => websocketConnectionState.value === WebSocket.OPEN)

function setupWebsocket() {
  // @todo Move the URL into a config
  websocket.value = new WebSocket('ws://localhost:3000/api/websocket')

  websocketConnectionState.value = websocket.value.readyState

  websocket.value.addEventListener('open', () => {
    websocketConnectionState.value = websocket.value?.readyState
  })
  websocket.value.addEventListener('close', () => {
    websocketConnectionState.value = websocket.value?.readyState
  })
  websocket.value.addEventListener('error', () => {
    websocketConnectionState.value = websocket.value?.readyState
  })
}

function getWebsocket() {
  if (websocket.value === null) {
    throw new Error('Websocket not expected to be null')
  }

  return websocket.value
}

export function useWebsocket() {
  return {
    setupWebsocket,
    websocket,
    getWebsocket,
    websocketConnectionState,
    websocketConnectionStateString,
    isWebsocketConnected,
  }
}
