const websocketStates = {
  undefined: 'not connected',
  [WebSocket.CONNECTING]: '... connecting ...',
  [WebSocket.OPEN]: 'connected',
  [WebSocket.CLOSING]: '... closing...',
  [WebSocket.CLOSED]: 'closed',
}

export function getWebsocketStateString(websocketState?: WebSocket['readyState']) {
  return websocketStates[websocketState ?? 'undefined']
}
