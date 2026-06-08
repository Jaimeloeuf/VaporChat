import { useAnonymousUser } from './useAnonymousUser'
import type { BaseChatUpdate, ChatUpdate } from './ChatUpdate'

// Helper type to check if an object has no keys
type Prettify<T> = { [K in keyof T]: T[K] } & {}
type ExtraFields<T extends ChatUpdate['type']> = Prettify<
  Omit<Extract<ChatUpdate, { type: T }>, keyof BaseChatUpdate>
>

export function createChatUpdate<T extends ChatUpdate['type']>(
  ...args: ExtraFields<T> extends Record<string, never>
    ? [type: T]
    : [type: T, additionalData: ExtraFields<T>]
): Extract<ChatUpdate, { type: T }> {
  const [type, additionalData] = args
  const { userID, readonlyUsername } = useAnonymousUser()
  return {
    id: crypto.randomUUID(),
    timestamp: new Date().toISOString(),
    userID,
    username: readonlyUsername.value,
    type,
    ...(additionalData || {}),
  } as Extract<ChatUpdate, { type: T }>
}
