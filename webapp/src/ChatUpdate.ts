import type { ChatConfig } from './ChatConfig'

export interface BaseChatUpdate {
  id: string
  timestamp: string
  userID: string
  username: string

  // To override in subtypes with a string literal
  type: string
}

export interface ChatUpdateRoomCreate extends BaseChatUpdate {
  type: 'room-create'
  payload: {
    chatConfig: ChatConfig
  }
}

export interface ChatUpdateRoomDestroy extends BaseChatUpdate {
  type: 'room-destroy'
  payload: {
    roomID: string
  }
}

export interface ChatUpdateNewStatusJoinRoom extends BaseChatUpdate {
  type: 'status-join-room'
  payload: {
    roomID: string
  }
}

export interface ChatUpdateNewStatusLeaveRoom extends BaseChatUpdate {
  type: 'status-leave-room'
  payload: {
    roomID: string
  }
}

// @todo In UI show for 2s from the last time of receiving this update
export interface ChatUpdateTyping extends BaseChatUpdate {
  type: 'typing'
}

export interface ChatUpdateMessageNew extends BaseChatUpdate {
  type: 'message-new'
  payload: {
    roomID: string
    message: string
  }
}

export interface ChatUpdateMessageDelete extends BaseChatUpdate {
  type: 'message-delete'
  payload: {
    messageID: string
  }
}

/**
 * `ChatUpdate` is both what you send out and what you receive
 */
export type ChatUpdate =
  | ChatUpdateRoomCreate
  | ChatUpdateRoomDestroy
  | ChatUpdateNewStatusJoinRoom
  | ChatUpdateNewStatusLeaveRoom
  | ChatUpdateTyping
  | ChatUpdateMessageNew
  | ChatUpdateMessageDelete
