import { reactive, readonly } from 'vue'
import type { ChatConfig } from './ChatConfig'

const chatConfig = reactive<ChatConfig>({
  chatRoomTTL: 300,
  maxNumberOfParticipants: 2,
  maxHistoryDurationInSeconds: 120,
  maxMessagesLength: 20,
})

export function useChatConfig() {
  return {
    chatConfig,
    readonlyChatConfig: readonly(chatConfig),
  }
}
