interface BaseChatUpdate {
  timestamp: string
  author: string

  // To override in subtypes with a string literal
  type: string
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

export interface ChatUpdateNewMessage extends BaseChatUpdate {
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
  | ChatUpdateNewStatusJoinRoom
  | ChatUpdateNewStatusLeaveRoom
  | ChatUpdateNewMessage
  | ChatUpdateMessageDelete
  | ChatUpdateTyping
