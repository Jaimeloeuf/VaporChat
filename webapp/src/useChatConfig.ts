import { reactive, readonly } from 'vue'
import { ChatConfigSchema } from './ChatConfig'

const chatConfig = reactive(
  ChatConfigSchema.parse({
    chatRoomTTL: 300,
    maxNumberOfParticipants: 2,
    maxHistoryDurationInSeconds: 120,
    maxMessagesLength: 20,
  }),
)

const readonlyChatConfig = readonly(chatConfig)

export function useChatConfig() {
  return {
    chatConfig,
    readonlyChatConfig,
  }
}
