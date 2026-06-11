import { z } from 'zod'
import { ChatConfigSchema } from './ChatConfig'

const BaseWsRequestSchema = z.object({
  id: z.uuidv4(),
  timestamp: z.number().int().positive(),
  userID: z.uuidv4(),
  username: z.string().nonempty(),

  /**
   * To override in subtypes with a string literal
   */
  type: z.string().nonempty(),
})
export type BaseWsRequest = z.infer<typeof BaseWsRequestSchema>

export const WsRequestRoomCreateSchema = BaseWsRequestSchema.extend({
  type: z.literal('room-create'),
  payload: z.object({
    chatConfig: ChatConfigSchema,
  }),
})
export type WsRequestRoomCreate = z.infer<typeof WsRequestRoomCreateSchema>

export const WsRequestRoomDestroySchema = BaseWsRequestSchema.extend({
  type: z.literal('room-destroy'),
  payload: z.object({
    roomID: z.uuidv4(),
  }),
})
export type WsRequestRoomDestroy = z.infer<typeof WsRequestRoomDestroySchema>

export const WsRequestNewStatusJoinRoomSchema = BaseWsRequestSchema.extend({
  type: z.literal('status-join-room'),
  payload: z.object({
    roomID: z.uuidv4(),
  }),
})
export type WsRequestNewStatusJoinRoom = z.infer<typeof WsRequestNewStatusJoinRoomSchema>

export const WsRequestNewStatusLeaveRoomSchema = BaseWsRequestSchema.extend({
  type: z.literal('status-leave-room'),
  payload: z.object({
    roomID: z.uuidv4(),
  }),
})
export type WsRequestNewStatusLeaveRoom = z.infer<typeof WsRequestNewStatusLeaveRoomSchema>

// @todo In UI show for 2s from the last time of receiving this update
export const WsRequestTypingSchema = BaseWsRequestSchema.extend({
  type: z.literal('typing'),
})
export type WsRequestTyping = z.infer<typeof WsRequestTypingSchema>

export const WsRequestMessageNewSchema = BaseWsRequestSchema.extend({
  type: z.literal('message-new'),
  payload: z.object({
    roomID: z.uuidv4(),
    message: z.string(),
  }),
})
export type WsRequestMessageNew = z.infer<typeof WsRequestMessageNewSchema>

export const WsRequestMessageDeleteSchema = BaseWsRequestSchema.extend({
  type: z.literal('message-delete'),
  payload: z.object({
    roomID: z.uuidv4(),
    messageID: z.uuidv4(),
  }),
})
export type WsRequestMessageDelete = z.infer<typeof WsRequestMessageDeleteSchema>

/**
 * Zod discriminated union for efficient request routing and parsing
 */
export const WsRequestSchema = z.discriminatedUnion('type', [
  WsRequestRoomCreateSchema,
  WsRequestRoomDestroySchema,
  WsRequestNewStatusJoinRoomSchema,
  WsRequestNewStatusLeaveRoomSchema,
  WsRequestTypingSchema,
  WsRequestMessageNewSchema,
  WsRequestMessageDeleteSchema,
])

/**
 * `WebSocket Request` is what the client sends to the server
 */
export type WsRequest = z.infer<typeof WsRequestSchema>
