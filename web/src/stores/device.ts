import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Device } from '@/types'
import axios from 'axios'

const API_BASE = import.meta.env.VITE_API_BASE || ''

export const useDeviceStore = defineStore('device', () => {
  const devices = ref<Device[]>([])
  const isLoading = ref(false)

  async function fetchDevices() {
    isLoading.value = true
    try {
      const response = await axios.get(`${API_BASE}/api/v1/devices`)
      devices.value = response.data.devices
    } finally {
      isLoading.value = false
    }
  }

  return {
    devices,
    isLoading,
    fetchDevices
  }
})