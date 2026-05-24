import type { ChatConfig } from './ChatConfig'

interface BaseChatUpdate {
  timestamp: string
  author: string

  // To override in subtypes with a string literal
  type: string
}

export interface ChatUpdateRoomCreate extends BaseChatUpdate {
  type: 'room-create'
  payload: {
    chatConfig: ChatConfig
  }
}

export interface ChatUpdateNewStatusJoinRoom extends BaseChatUpdate {
  type: 'status-join-room'
}

export interface ChatUpdateNewStatusLeaveRoom extends BaseChatUpdate {
  type: 'status-leave-room'
}

// @todo In UI show for 2s from the last time of receiving this update
export interface ChatUpdateTyping extends BaseChatUpdate {
  type: 'typing'
}

export interface ChatUpdateMessageNew extends BaseChatUpdate {
  type: 'message-new'
  payload: {
    message: string
  }
}

export interface ChatUpdateMessageDelete extends BaseChatUpdate {
  type: 'message-delete'
  payload: {
    messageID: string
  }
}

export type ChatUpdate =
  | ChatUpdateRoomCreate
  | ChatUpdateNewStatusJoinRoom
  | ChatUpdateNewStatusLeaveRoom
  | ChatUpdateTyping
  | ChatUpdateMessageNew
  | ChatUpdateMessageDelete

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
  return {
    timestamp: new Date().toISOString(),
    // @todo
    author: 'author',
    type,
    ...(additionalData || {}),
  } as Extract<ChatUpdate, { type: T }>
}
