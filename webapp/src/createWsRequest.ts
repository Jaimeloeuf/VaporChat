import { useAnonymousUser } from './useAnonymousUser'
import type { BaseWsRequest, WsRequest } from './WsRequest'

// Helper type to check if an object has no keys
type Prettify<T> = { [K in keyof T]: T[K] } & {}
type ExtraFields<T extends WsRequest['type']> = Prettify<
  Omit<Extract<WsRequest, { type: T }>, keyof BaseWsRequest>
>

export function createWsRequest<T extends WsRequest['type']>(
  ...args: ExtraFields<T> extends Record<string, never>
    ? [type: T]
    : [type: T, additionalData: ExtraFields<T>]
): Extract<WsRequest, { type: T }> {
  const [type, additionalData] = args
  const { userID, readonlyUsername } = useAnonymousUser()
  return {
    id: crypto.randomUUID(),
    timestamp: new Date().toISOString(),
    userID,
    username: readonlyUsername.value,
    type,
    ...(additionalData || {}),
  } as Extract<WsRequest, { type: T }>
}
