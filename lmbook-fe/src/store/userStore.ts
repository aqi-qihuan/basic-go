import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User } from '@/types/user'
import request from '@/utils/request'

interface UserState {
  user: User | null
  token: string | null
  isLogin: boolean
  login: (email: string, password: string) => Promise<void>
  register: (data: { email: string; password: string; nickname: string }) => Promise<void>
  logout: () => void
  checkAuth: () => boolean
  setUser: (user: User | null) => void
  setToken: (token: string | null) => void
}

export const useUserStore = create<UserState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isLogin: false,

      login: async (email: string, password: string) => {
        try {
          const response = await request.post<any>('/user/login', {
            email,
            password,
          })
          const { token, user } = response.data.data
          localStorage.setItem('token', token)
          set({ user, token, isLogin: true })
        } catch (error) {
          throw error
        }
      },

      register: async (data: { email: string; password: string; nickname: string }) => {
        try {
          await request.post('/user/register', data)
        } catch (error) {
          throw error
        }
      },

      logout: () => {
        localStorage.removeItem('token')
        set({ user: null, token: null, isLogin: false })
      },

      checkAuth: () => {
        const token = localStorage.getItem('token')
        if (token) {
          set({ token, isLogin: true })
          return true
        }
        return false
      },

      setUser: (user: User | null) => {
        set({ user })
      },

      setToken: (token: string | null) => {
        set({ token, isLogin: !!token })
      },
    }),
    {
      name: 'user-storage',
      partialize: (state) => ({ user: state.user, token: state.token, isLogin: state.isLogin }),
    }
  )
)
