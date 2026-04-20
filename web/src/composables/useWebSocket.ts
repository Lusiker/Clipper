import { ref, computed } from 'vue'

const ws = ref<WebSocket | null>(null)
const isConnected = ref(false)
const reconnectAttempts = ref(0)
const maxReconnectAttempts = 5

export function useWebSocket() {
  let messageHandler: ((event: MessageEvent) => void) | null = null

  function connect(deviceId: string, deviceName: string, handler: (event: MessageEvent) => void) {
    if (ws.value) {
      disconnect()
    }

    messageHandler = handler

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = import.meta.env.VITE_WS_HOST || window.location.host
    const wsUrl = `${protocol}//${host}/api/v1/devices/ws?device_id=${deviceId}&device_name=${encodeURIComponent(deviceName)}`

    ws.value = new WebSocket(wsUrl)

    ws.value.onopen = () => {
      isConnected.value = true
      reconnectAttempts.value = 0
      console.log('[WS] Connected')
    }

    ws.value.onmessage = (event) => {
      if (messageHandler) {
        messageHandler(event)
      }
    }

    ws.value.onclose = () => {
      isConnected.value = false
      console.log('[WS] Disconnected')

      if (reconnectAttempts.value < maxReconnectAttempts) {
        reconnectAttempts.value++
        setTimeout(() => {
          if (deviceId && deviceName && messageHandler) {
            connect(deviceId, deviceName, messageHandler)
          }
        }, 2000 * reconnectAttempts.value)
      }
    }

    ws.value.onerror = (error) => {
      console.error('[WS] Error:', error)
    }
  }

  function disconnect() {
    if (ws.value) {
      ws.value.close()
      ws.value = null
    }
    isConnected.value = false
    messageHandler = null
  }

  function send(type: string, data: any) {
    if (ws.value && isConnected.value) {
      ws.value.send(JSON.stringify({ type, data }))
    }
  }

  return {
    ws,
    isConnected,
    connect,
    disconnect,
    send
  }
}