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

// @todo Support Delete message

export type ChatUpdate =
  | ChatUpdateNewStatusJoinRoom
  | ChatUpdateNewStatusLeaveRoom
  | ChatUpdateNewMessage
  | ChatUpdateTyping
