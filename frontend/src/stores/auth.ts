import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

interface AuthUser {
  id: number
  username: string
  email: string
}

const TOKEN_KEY = 'token'
const USER_KEY = 'user'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>('')
  const user = ref<AuthUser | null>(null)

  const isAuthenticated = computed(() => !!token.value)

  const initFromStorage = () => {
    const savedToken = localStorage.getItem(TOKEN_KEY)
    const savedUser = localStorage.getItem(USER_KEY)

    token.value = savedToken || ''

    if (savedUser) {
      try {
        user.value = JSON.parse(savedUser) as AuthUser
      } catch {
        user.value = null
        localStorage.removeItem(USER_KEY)
      }
    } else {
      user.value = null
    }
  }

  const setAuth = (nextToken: string, nextUser: AuthUser) => {
    token.value = nextToken
    user.value = nextUser
    localStorage.setItem(TOKEN_KEY, nextToken)
    localStorage.setItem(USER_KEY, JSON.stringify(nextUser))
  }

  const logout = () => {
    token.value = ''
    user.value = null
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
  }

  return {
    token,
    user,
    isAuthenticated,
    initFromStorage,
    setAuth,
    logout,
  }
})
