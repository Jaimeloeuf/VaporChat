export type Message = {
  timestamp: string
  author: 'current-user' | 'other-user' | 'system'
  message: string
}
