import type { ChatConfig } from './ChatConfig'

export interface BaseWsRequest {
  id: string
  timestamp: string
  userID: string
  username: string

  // To override in subtypes with a string literal
  type: string
}

export interface WsRequestRoomCreate extends BaseWsRequest {
  type: 'room-create'
  payload: {
    chatConfig: ChatConfig
  }
}

export interface WsRequestRoomDestroy extends BaseWsRequest {
  type: 'room-destroy'
  payload: {
    roomID: string
  }
}

export interface WsRequestNewStatusJoinRoom extends BaseWsRequest {
  type: 'status-join-room'
  payload: {
    roomID: string
  }
}

export interface WsRequestNewStatusLeaveRoom extends BaseWsRequest {
  type: 'status-leave-room'
  payload: {
    roomID: string
  }
}

// @todo In UI show for 2s from the last time of receiving this update
export interface WsRequestTyping extends BaseWsRequest {
  type: 'typing'
}

export interface WsRequestMessageNew extends BaseWsRequest {
  type: 'message-new'
  payload: {
    roomID: string
    message: string
  }
}

export interface WsRequestMessageDelete extends BaseWsRequest {
  type: 'message-delete'
  payload: {
    messageID: string
  }
}

/**
 * `WebSocket Request` is what the client sends to the server
 */
export type WsRequest =
  | WsRequestRoomCreate
  | WsRequestRoomDestroy
  | WsRequestNewStatusJoinRoom
  | WsRequestNewStatusLeaveRoom
  | WsRequestTyping
  | WsRequestMessageNew
  | WsRequestMessageDelete
