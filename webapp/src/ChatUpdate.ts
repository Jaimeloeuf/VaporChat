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

export interface ChatUpdateNewMessage extends BaseChatUpdate {
  type: 'message'
  payload: {
    message: string
  }
}

// @todo In UI show for 2s from the last time of receiving this update
export interface ChatUpdateTyping extends BaseChatUpdate {
  type: 'typing'
}

// @todo Support Delete message

export type ChatUpdate =
  | ChatUpdateNewStatusJoinRoom
  | ChatUpdateNewStatusLeaveRoom
  | ChatUpdateNewMessage
  | ChatUpdateTyping
