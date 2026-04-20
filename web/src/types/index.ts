// User types
export interface User {
  id: string
  username: string
  created_at: string
  updated_at: string
}

// Clip types
export type ClipType = 'text' | 'image'

export interface Clip {
  id: string
  user_id: string
  device_id: string
  type: ClipType
  content: string
  meta?: string
  created_at: string
  updated_at: string
}

export interface ClipCreate {
  type: ClipType
  content: string
  meta?: string
}

export interface ClipMetaImage {
  width: number
  height: number
  size: number
  format: string
  thumb_path: string
}

// Device types
export interface Device {
  id: string
  user_id: string
  name: string
  ip: string
  last_seen: string
  is_online: boolean
}

// WebSocket message types
export interface WSMessage<T = any> {
  type: string
  data: T
}

export interface ClipCreatedData {
  clip: Clip
}

export interface ClipDeletedData {
  clip_id: string
}

export interface DeviceStatusData {
  device_id: string
  device_name: string
  is_online: boolean
}