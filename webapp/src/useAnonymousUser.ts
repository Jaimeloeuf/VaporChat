import { ref, readonly } from 'vue'

import { getRandomAnonymousAnimalName } from './getRandomAnimalName'

const userID = crypto.randomUUID()

const username = ref(getRandomAnonymousAnimalName())

const readonlyUsername = readonly(username)

/**
 * A single consistent anonymous user "account" per "session"
 */
export function useAnonymousUser() {
  return Object.freeze({
    userID,
    username,
    readonlyUsername,
  })
}
