import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Clip, ClipCreate, ClipMetaImage } from '@/types'
import axios from 'axios'

const API_BASE = import.meta.env.VITE_API_BASE || ''

export const useClipStore = defineStore('clip', () => {
  const clips = ref<Clip[]>([])
  const isLoading = ref(false)

  async function fetchClips() {
    isLoading.value = true
    try {
      const response = await axios.get(`${API_BASE}/api/v1/clips`)
      clips.value = response.data.clips
    } finally {
      isLoading.value = false
    }
  }

  async function createClip(data: ClipCreate, deviceId: string) {
    try {
      const response = await axios.post(`${API_BASE}/api/v1/clips?device_id=${deviceId}`, data)
      const newClip = response.data.clip
      clips.value.unshift(newClip)
      return newClip
    } catch (error: any) {
      throw error.response?.data?.error || 'Failed to create clip'
    }
  }

  async function uploadImage(file: File, deviceId: string): Promise<Clip> {
    const formData = new FormData()
    formData.append('image', file)

    try {
      const response = await axios.post(`${API_BASE}/api/v1/clips/upload?device_id=${deviceId}`, formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      })
      const newClip = response.data.clip
      clips.value.unshift(newClip)
      return newClip
    } catch (error: any) {
      throw error.response?.data?.error || 'Failed to upload image'
    }
  }

  async function deleteClip(clipId: string) {
    try {
      await axios.delete(`${API_BASE}/api/v1/clips/${clipId}`)
      clips.value = clips.value.filter(c => c.id !== clipId)
    } catch (error: any) {
      throw error.response?.data?.error || 'Failed to delete clip'
    }
  }

  function addClip(clip: Clip) {
    const exists = clips.value.find(c => c.id === clip.id)
    if (!exists) {
      clips.value.unshift(clip)
    }
  }

  function removeClip(clipId: string) {
    clips.value = clips.value.filter(c => c.id !== clipId)
  }

  function getThumbUrl(clip: Clip): string {
    if (clip.type !== 'image' || !clip.meta) return ''
    try {
      const meta: ClipMetaImage = JSON.parse(clip.meta)
      return `${API_BASE}/uploads/${meta.thumb_path}`
    } catch {
      return ''
    }
  }

  function getOrigUrl(clip: Clip): string {
    if (clip.type !== 'image') return ''
    return `${API_BASE}/uploads/${clip.content}`
  }

  return {
    clips,
    isLoading,
    fetchClips,
    createClip,
    uploadImage,
    deleteClip,
    addClip,
    removeClip,
    getThumbUrl,
    getOrigUrl
  }
})