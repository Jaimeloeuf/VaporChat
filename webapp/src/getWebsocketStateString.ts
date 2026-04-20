const websocketStates = {
  undefined: 'Not connected',
  [WebSocket.CONNECTING]: 'Connecting',
  [WebSocket.OPEN]: 'Connected',
  [WebSocket.CLOSING]: 'Connected',
  [WebSocket.CLOSED]: 'Connected',
}

export function getWebsocketStateString(websocketState?: WebSocket['readyState']) {
  return websocketStates[websocketState ?? 'undefined']
}
