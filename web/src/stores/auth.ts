import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'
import axios from 'axios'

const API_BASE = import.meta.env.VITE_API_BASE || ''

// 确保 axios 发送 cookie
axios.defaults.withCredentials = true

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const isLoading = ref(false)

  const isLoggedIn = computed(() => user.value !== null)

  async function checkAuth() {
    try {
      const response = await axios.get(`${API_BASE}/api/v1/auth/me`)
      user.value = response.data.user
    } catch {
      user.value = null
    }
  }

  async function login(username: string, password: string) {
    isLoading.value = true
    try {
      const response = await axios.post(`${API_BASE}/api/v1/auth/login`, {
        username,
        password
      })
      user.value = response.data.user
      return true
    } catch (error: any) {
      throw error.response?.data?.error || 'Login failed'
    } finally {
      isLoading.value = false
    }
  }

  async function register(username: string, password: string) {
    isLoading.value = true
    try {
      const response = await axios.post(`${API_BASE}/api/v1/auth/register`, {
        username,
        password
      })
      user.value = response.data.user
      return true
    } catch (error: any) {
      throw error.response?.data?.error || 'Registration failed'
    } finally {
      isLoading.value = false
    }
  }

  async function logout() {
    try {
      await axios.post(`${API_BASE}/api/v1/auth/logout`)
    } catch {
      // Ignore logout errors
    }
    user.value = null
  }

  return {
    user,
    isLoading,
    isLoggedIn,
    checkAuth,
    login,
    register,
    logout
  }
})