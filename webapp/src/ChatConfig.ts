import { z } from 'zod'

export const ChatConfigSchema = z.object({
  chatRoomTTL: z.number().positive(),
  maxNumberOfParticipants: z.number().positive(),
  maxHistoryDurationInSeconds: z.number().positive(),
  maxMessagesLength: z.number().positive(),
})

export type ChatConfig = z.infer<typeof ChatConfigSchema>
