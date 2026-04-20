<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { useClipStore } from '@/stores/clip'
import { useWebSocket } from '@/composables/useWebSocket'
import { resetAuthCheck } from '@/router'
import { DocumentCopy, Delete, Plus, Connection, User, SwitchButton, Download, Picture, Monitor, Edit } from '@element-plus/icons-vue'
import type { Clip } from '@/types'

const router = useRouter()
const authStore = useAuthStore()
const clipStore = useClipStore()
const { isConnected, connect, disconnect } = useWebSocket()

const deviceId = ref('')
const deviceName = ref('')
const newClipContent = ref('')
const newClipType = ref<'text' | 'image'>('text')
const deviceNameDialogVisible = ref(false)
const editingDeviceName = ref('')

// 图片上传
const uploadProgress = ref(0)
const isUploading = ref(false)
const previewVisible = ref(false)
const previewUrl = ref('')

// 分页
const currentPage = ref(1)
const pageSize = ref(5)

onMounted(async () => {
  deviceId.value = localStorage.getItem('deviceId') || `device-${Date.now()}`
  localStorage.setItem('deviceId', deviceId.value)

  // 读取或生成设备名称
  deviceName.value = localStorage.getItem('deviceName') || ''
  if (!deviceName.value) {
    const platform = navigator.platform || 'Unknown'
    const isMobile = /iPhone|iPad|Android/i.test(navigator.userAgent)
    deviceName.value = isMobile ? 'Mobile Device' : platform.includes('Win') ? 'Windows PC' : platform.includes('Mac') ? 'Mac' : 'Device'
    localStorage.setItem('deviceName', deviceName.value)
  }

  await clipStore.fetchClips()
  connect(deviceId.value, deviceName.value, handleWSMessage)
})

const paginatedClips = computed(() => {
  const clips = clipStore.clips || []
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return clips.slice(start, end)
})

const totalClips = computed(() => clipStore.clips?.length || 0)

function handleWSMessage(event: MessageEvent) {
  try {
    const message = JSON.parse(event.data)
    console.log('[WS] Message:', message.type)

    switch (message.type) {
      case 'clip_created':
        clipStore.addClip(message.data.clip)
        ElMessage.success('New clip received from another device')
        break
      case 'clip_deleted':
        clipStore.removeClip(message.data.clip_id)
        break
      case 'device_online':
        ElMessage.info(`Device "${message.data.device_name}" connected`)
        break
      case 'device_offline':
        ElMessage.info(`Device "${message.data.device_name}" disconnected`)
        break
    }
  } catch (error) {
    console.error('[WS] Parse error:', error)
  }
}

async function createClip() {
  if (!newClipContent.value) {
    ElMessage.warning('Please enter content')
    return
  }

  try {
    await clipStore.createClip({
      type: newClipType.value,
      content: newClipContent.value
    }, deviceId.value)
    newClipContent.value = ''
    ElMessage.success('Clip created and synced')
  } catch (error: any) {
    ElMessage.error(error)
  }
}

async function handleImageUpload(uploadFile: any) {
  if (!uploadFile?.raw) {
    return false
  }

  const file = uploadFile.raw as File
  isUploading.value = true
  uploadProgress.value = 0

  try {
    await clipStore.uploadImage(file, deviceId.value)
    ElMessage.success('Image uploaded and synced')
  } catch (error: any) {
    ElMessage.error(error)
  } finally {
    isUploading.value = false
    uploadProgress.value = 0
  }

  return false // Prevent el-upload default behavior
}

function showImagePreview(clip: Clip) {
  previewUrl.value = clipStore.getOrigUrl(clip)
  previewVisible.value = true
}

async function copyImageToClipboard(clip: Clip) {
  try {
    const response = await fetch(clipStore.getOrigUrl(clip))
    const blob = await response.blob()
    await navigator.clipboard.write([
      new ClipboardItem({ [blob.type]: blob })
    ])
    ElMessage.success('Image copied to clipboard')
  } catch (error) {
    ElMessage.error('Failed to copy image')
  }
}

function downloadImage(clip: Clip) {
  const url = clipStore.getOrigUrl(clip)
  const link = document.createElement('a')
  link.href = url
  link.download = `clip-${clip.id}.jpg`
  link.click()
}

async function copyPreviewImage() {
  try {
    const response = await fetch(previewUrl.value)
    const blob = await response.blob()
    await navigator.clipboard.write([
      new ClipboardItem({ [blob.type]: blob })
    ])
    ElMessage.success('Image copied to clipboard')
  } catch (error) {
    ElMessage.error('Failed to copy image')
  }
}

function downloadPreviewImage() {
  const link = document.createElement('a')
  link.href = previewUrl.value
  link.download = 'clipper-image.jpg'
  link.click()
}

async function deleteClip(clip: Clip) {
  const message = clip.type === 'image'
    ? 'Delete this image clip and its files?'
    : 'Delete this clip?'

  try {
    await ElMessageBox.confirm(message, 'Confirm Delete', {
      confirmButtonText: 'Delete',
      cancelButtonText: 'Cancel',
      type: 'warning'
    })

    await clipStore.deleteClip(clip.id)
    // 如果当前页删除后没有数据，返回上一页
    if (paginatedClips.value.length === 0 && currentPage.value > 1) {
      currentPage.value--
    }
    // 如果正在预览这个图片，关闭预览
    if (previewUrl.value.includes(clip.id)) {
      previewVisible.value = false
    }
    ElMessage.success('Clip deleted')
  } catch {
    // 用户取消
  }
}

function copyToClipboard(content: string) {
  const textarea = document.createElement('textarea')
  textarea.value = content
  textarea.style.position = 'fixed'
  textarea.style.left = '-9999px'
  textarea.style.top = '0'
  textarea.setAttribute('readonly', '')
  document.body.appendChild(textarea)

  textarea.focus()
  textarea.select()

  try {
    const success = document.execCommand('copy')
    if (success) {
      ElMessage.success('Copied to clipboard')
    } else {
      ElMessage.error('Copy failed')
    }
  } catch (err) {
    ElMessage.error('Copy failed')
  }

  document.body.removeChild(textarea)
}

async function logout() {
  disconnect()
  await authStore.logout()
  resetAuthCheck()
  router.push('/login')
}

function openDeviceNameDialog() {
  editingDeviceName.value = deviceName.value
  deviceNameDialogVisible.value = true
}

async function saveDeviceName() {
  if (!editingDeviceName.value.trim()) {
    ElMessage.warning('Device name cannot be empty')
    return
  }

  deviceName.value = editingDeviceName.value.trim()
  localStorage.setItem('deviceName', deviceName.value)

  // 重新连接WebSocket以更新设备名称
  disconnect()
  connect(deviceId.value, deviceName.value, handleWSMessage)

  deviceNameDialogVisible.value = false
  ElMessage.success('Device name updated')
}
</script>

<template>
  <div class="dashboard">
    <header class="header">
      <div class="header-left">
        <div class="logo-icon-small">
          <svg viewBox="0 0 24 24" width="28" height="28" fill="#409eff">
            <path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-7 14H6v-2h6v2zm4-4H6v-2h10v2zm0-4H6V7h10v2z"/>
          </svg>
        </div>
        <h1>Clipper</h1>
      </div>
      <div class="header-right">
        <div class="user-info">
          <el-icon><User /></el-icon>
          <span>{{ authStore.user?.username }}</span>
        </div>
        <div class="device-info" @click="openDeviceNameDialog">
          <el-icon><Monitor /></el-icon>
          <span>{{ deviceName }}</span>
          <el-icon class="edit-icon"><Edit /></el-icon>
        </div>
        <div class="connection-status" :class="{ connected: isConnected }">
          <el-icon><Connection /></el-icon>
          <span>{{ isConnected ? 'Connected' : 'Disconnected' }}</span>
        </div>
        <el-button type="danger" size="small" @click="logout">
          <el-icon><SwitchButton /></el-icon>
          Logout
        </el-button>
      </div>
    </header>

    <main class="main">
      <div class="layout-container">
        <!-- 左侧创建区域 -->
        <div class="create-section">
          <el-card class="create-card">
            <template #header>
              <div class="card-header">
                <el-icon class="header-icon"><Plus /></el-icon>
                <span>Create New Clip</span>
              </div>
            </template>

            <div class="type-selector">
              <div class="segment-control">
                <div class="segment-slider" :class="{ 'slider-right': newClipType === 'image' }"></div>
                <div
                  class="segment-option"
                  :class="{ active: newClipType === 'text' }"
                  @click="newClipType = 'text'"
                >
                  <el-icon><DocumentCopy /></el-icon>
                  <span>Text</span>
                </div>
                <div
                  class="segment-option"
                  :class="{ active: newClipType === 'image' }"
                  @click="newClipType = 'image'"
                >
                  <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
                    <path d="M21 19V5c0-1.1-.9-2-2-2H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2zM8.5 13.5l2.5 3.01L14.5 12l4.5 6H5l3.5-4.5z"/>
                  </svg>
                  <span>Image</span>
                </div>
              </div>
            </div>

            <div class="content-wrapper">
              <div class="content-slider" :class="{ 'slide-right': newClipType === 'image' }">
                <div class="content-panel text-panel">
                  <el-input
                    v-model="newClipContent"
                    type="textarea"
                    :autosize="{ minRows: 6, maxRows: 20 }"
                    placeholder="Enter text content to sync across devices..."
                    resize="none"
                  />
                </div>
                <div class="content-panel image-panel">
                  <el-upload
                    class="image-upload"
                    drag
                    :auto-upload="false"
                    :show-file-list="false"
                    accept="image/jpeg,image/png,image/gif,image/webp"
                    :on-change="(uploadFile: any) => handleImageUpload(uploadFile)"
                    :disabled="isUploading"
                  >
                    <div class="upload-area" v-if="!isUploading">
                      <el-icon size="48"><Picture /></el-icon>
                      <p>Drag image here or click to upload</p>
                      <span class="upload-tip">Supports JPEG, PNG, GIF, WebP (max 20MB)</span>
                    </div>
                    <div class="upload-progress" v-else>
                      <el-progress type="circle" :percentage="uploadProgress" :width="80" />
                      <p>Uploading...</p>
                    </div>
                  </el-upload>
                </div>
              </div>
            </div>

            <el-button
              type="primary"
              size="large"
              @click="createClip"
              class="create-btn"
              :disabled="!newClipContent"
            >
              <el-icon><Plus /></el-icon>
              Create & Sync
            </el-button>
          </el-card>
        </div>

        <!-- 右侧 Clips 列表 -->
        <div class="clips-section">
          <el-card class="clips-card">
            <template #header>
              <div class="clips-card-header">
                <div class="header-left-info">
                  <el-icon class="header-icon"><DocumentCopy /></el-icon>
                  <span>Your Clips</span>
                </div>
                <span class="clip-count">{{ totalClips }} items</span>
              </div>
            </template>

            <div v-if="!clipStore.clips?.length" class="empty-state">
              <el-icon size="48" color="#c0c4cc"><DocumentCopy /></el-icon>
              <p>No clips yet</p>
              <span>Create your first clip</span>
            </div>

            <div v-else class="clips-container">
              <div class="clips-list">
                <el-card v-for="clip in paginatedClips" :key="clip.id" class="clip-card" shadow="hover">
                  <div class="clip-header">
                    <div class="clip-meta">
                      <span class="clip-type">
                        <el-icon size="12"><DocumentCopy /></el-icon>
                        {{ clip.type }}
                      </span>
                      <span class="clip-time">{{ new Date(clip.created_at).toLocaleString() }}</span>
                    </div>
                  </div>

                  <div class="clip-content">
                    <template v-if="clip.type === 'text'">
                      <p class="text-content">{{ clip.content }}</p>
                    </template>
                    <template v-else>
                      <img
                        :src="clipStore.getThumbUrl(clip)"
                        alt="Clip image"
                        class="image-thumb"
                        @click="showImagePreview(clip)"
                      />
                    </template>
                  </div>

                  <div class="clip-actions">
                    <el-button
                      v-if="clip.type === 'text'"
                      type="primary"
                      size="small"
                      plain
                      @click="copyToClipboard(clip.content)"
                    >
                      <el-icon><DocumentCopy /></el-icon>
                      Copy
                    </el-button>
                    <el-button
                      v-if="clip.type === 'image'"
                      type="primary"
                      size="small"
                      plain
                      @click="copyImageToClipboard(clip)"
                    >
                      <el-icon><DocumentCopy /></el-icon>
                      Copy
                    </el-button>
                    <el-button
                      v-if="clip.type === 'image'"
                      type="success"
                      size="small"
                      plain
                      @click="downloadImage(clip)"
                    >
                      <el-icon><Download /></el-icon>
                      Download
                    </el-button>
                    <el-button
                      type="danger"
                      size="small"
                      plain
                      @click="deleteClip(clip)"
                    >
                      <el-icon><Delete /></el-icon>
                      Delete
                    </el-button>
                  </div>
                </el-card>
              </div>

              <div class="pagination-wrapper">
                <el-pagination
                  v-model:current-page="currentPage"
                  v-model:page-size="pageSize"
                  :page-sizes="[5, 10, 20]"
                  :total="totalClips"
                  layout="sizes, prev, pager, next"
                  small
                  background
                />
              </div>
            </div>
          </el-card>
        </div>
      </div>
    </main>

    <!-- 图片预览弹窗 -->
    <el-dialog
      v-model="previewVisible"
      width="80%"
      :show-close="false"
      destroy-on-close
      class="preview-dialog"
    >
      <div class="preview-content">
        <div class="preview-header">
          <div class="preview-header-left">
            <svg viewBox="0 0 24 24" width="20" height="20" fill="white">
              <path d="M21 19V5c0-1.1-.9-2-2-2H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2zM8.5 13.5l2.5 3.01L14.5 12l4.5 6H5l3.5-4.5z"/>
            </svg>
            <span>Image Preview</span>
          </div>
          <div class="preview-header-actions">
            <el-button type="primary" size="small" @click="copyPreviewImage">
              <el-icon><DocumentCopy /></el-icon>
              Copy
            </el-button>
            <el-button type="success" size="small" @click="downloadPreviewImage">
              <el-icon><Download /></el-icon>
              Download
            </el-button>
            <el-button size="small" @click="previewVisible = false">
              <el-icon><SwitchButton /></el-icon>
              Close
            </el-button>
          </div>
        </div>

        <div class="preview-body">
          <img :src="previewUrl" alt="Preview" class="preview-image" />
        </div>
      </div>
    </el-dialog>

    <!-- 设备名称编辑弹窗 -->
    <el-dialog
      v-model="deviceNameDialogVisible"
      title="Edit Device Name"
      width="400px"
      destroy-on-close
    >
      <el-input
        v-model="editingDeviceName"
        placeholder="Enter device name"
        size="large"
        clearable
      />
      <template #footer>
        <el-button @click="deviceNameDialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="saveDeviceName">Save</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.dashboard {
  min-height: 100vh;
  background: linear-gradient(135deg, #f0f7ff 0%, #e8f4fc 100%);
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: white;
  border-bottom: 1px solid #e4e7ed;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon-small {
  display: flex;
  align-items: center;
}

.header h1 {
  color: #409eff;
  margin: 0;
  font-size: 22px;
  font-weight: 600;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #606266;
  font-size: 14px;
  padding: 6px 12px;
  background: #f4f4f5;
  border-radius: 6px;
}

.device-info {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #909399;
  font-size: 13px;
  padding: 6px 10px;
  background: #f4f4f5;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.2s, color 0.2s;
}

.device-info:hover {
  background: #e9e9eb;
  color: #606266;
}

.device-info .edit-icon {
  font-size: 12px;
  opacity: 0.6;
}

.device-info:hover .edit-icon {
  opacity: 1;
}

.header-right {
  display: flex;
  gap: 12px;
  align-items: center;
}

.connection-status {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border-radius: 6px;
  background: #fef0f0;
  color: #f56c6c;
  font-size: 13px;
  font-weight: 500;
}

.connection-status.connected {
  background: #e6f7e6;
  color: #52c41a;
}

.main {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.layout-container {
  display: flex;
  gap: 24px;
  align-items: flex-start;
}

/* 左侧创建区域 */
.create-section {
  flex: 0 0 480px;
  min-width: 420px;
}

.create-card {
  border-radius: 12px;
  border: none;
}

.create-card :deep(.el-card__header) {
  padding: 16px 20px;
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
  border-radius: 12px 12px 0 0;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  color: white;
  font-size: 16px;
  font-weight: 500;
}

.header-icon {
  font-size: 18px;
}

.create-card :deep(.el-card__body) {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.type-selector {
  flex-shrink: 0;
}

.segment-control {
  display: flex;
  position: relative;
  background: #f0f2f5;
  border-radius: 12px;
  padding: 4px;
  gap: 4px;
}

.segment-slider {
  position: absolute;
  top: 4px;
  left: 4px;
  width: calc(50% - 4px);
  height: calc(100% - 8px);
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
  border-radius: 10px;
  transition: left 0.5s cubic-bezier(0.68, -0.55, 0.265, 1.55);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.3);
}

.segment-slider.slider-right {
  left: calc(50%);
}

.segment-option {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 12px 20px;
  cursor: pointer;
  color: #606266;
  font-size: 14px;
  font-weight: 500;
  transition: color 0.3s ease;
  z-index: 1;
  border-radius: 10px;
  user-select: none;
}

.segment-option:hover {
  color: #303133;
}

.segment-option.active {
  color: white;
}

.segment-option.active:hover {
  color: white;
}

/* 内容滑块容器 - 防止穿帮 */
.content-wrapper {
  position: relative;
  overflow: hidden;
  border-radius: 8px;
}

.content-slider {
  display: flex;
  width: 200%;
  transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.content-slider.slide-right {
  transform: translateX(-50%);
}

.content-panel {
  flex: 0 0 50%;
  width: 50%;
  box-sizing: border-box;
}

.text-panel :deep(.el-textarea__inner) {
  border-radius: 8px;
  padding: 14px;
  font-size: 14px;
  min-height: 150px !important;
  transition: height 0.3s ease;
}

.image-upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 120px;
  padding: 32px;
  background: #f4f4f5;
  border-radius: 8px;
  color: #909399;
}

.image-upload-placeholder p {
  margin: 8px 0 0;
  font-size: 14px;
}

.image-upload {
  width: 100%;
}

.image-upload :deep(.el-upload-dragger) {
  width: 100%;
  height: 180px;
  border-radius: 8px;
  border: 2px dashed #d9d9d9;
  background: #f4f4f5;
  transition: border-color 0.3s;
}

.image-upload :deep(.el-upload-dragger:hover) {
  border-color: #409eff;
}

.upload-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #909399;
}

.upload-area p {
  margin: 12px 0 4px;
  font-size: 14px;
  color: #606266;
}

.upload-tip {
  font-size: 12px;
  color: #c0c4cc;
}

.upload-progress {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.upload-progress p {
  margin-top: 12px;
  color: #606266;
}

.image-thumb {
  max-width: 100%;
  max-height: 200px;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.2s;
}

.image-thumb:hover {
  transform: scale(1.02);
}

/* 图片预览弹窗样式 */
.preview-dialog :deep(.el-dialog) {
  background: transparent !important;
  box-shadow: none !important;
  border: none !important;
}

.preview-dialog :deep(.el-dialog__header) {
  display: none !important;
}

.preview-dialog :deep(.el-dialog__body) {
  padding: 0 !important;
  background: transparent !important;
}

.preview-content {
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  overflow: hidden;
}

.preview-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 20px;
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
}

.preview-header-left {
  display: flex;
  align-items: center;
  gap: 10px;
  color: white;
  font-size: 16px;
  font-weight: 500;
}

.preview-header-actions {
  display: flex;
  gap: 10px;
}

.preview-header-actions :deep(.el-button) {
  border-radius: 8px;
  font-weight: 500;
}

.preview-header-actions :deep(.el-button--primary) {
  background: rgba(255, 255, 255, 0.9);
  border-color: transparent;
  color: #409eff;
}

.preview-header-actions :deep(.el-button--primary:hover) {
  background: white;
}

.preview-header-actions :deep(.el-button--success) {
  background: rgba(255, 255, 255, 0.9);
  border-color: transparent;
  color: #67c23a;
}

.preview-header-actions :deep(.el-button--success:hover) {
  background: white;
}

.preview-header-actions :deep(.el-button:not(.el-button--primary):not(.el-button--success)) {
  background: rgba(255, 255, 255, 0.2);
  border-color: rgba(255, 255, 255, 0.3);
  color: white;
}

.preview-header-actions :deep(.el-button:not(.el-button--primary):not(.el-button--success):hover) {
  background: rgba(255, 255, 255, 0.3);
}

.preview-body {
  padding: 20px;
  background: linear-gradient(135deg, #f0f7ff 0%, #e8f4fc 100%);
  min-height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.preview-image {
  max-width: 100%;
  max-height: 70vh;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease;
}

.preview-image:hover {
  transform: scale(1.02);
}

.create-btn {
  width: 100%;
  height: 44px;
  font-size: 15px;
  font-weight: 500;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  flex-shrink: 0;
}

/* 右侧 Clips 列表 */
.clips-section {
  flex: 1;
  min-width: 0;
}

.clips-card {
  border-radius: 12px;
  border: none;
  max-height: calc(100vh - 180px);
  display: flex;
  flex-direction: column;
}

.clips-card :deep(.el-card__header) {
  padding: 16px 20px;
  background: linear-gradient(135deg, #67c23a 0%, #85ce61 100%);
  border-radius: 12px 12px 0 0;
}

.clips-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.clips-card-header .header-left-info {
  display: flex;
  align-items: center;
  gap: 8px;
  color: white;
  font-size: 16px;
  font-weight: 500;
}

.clips-card-header .header-icon {
  font-size: 18px;
}

.clips-card-header .clip-count {
  display: inline-flex;
  align-items: center;
  color: white;
  font-size: 13px;
  font-weight: 500;
  background: rgba(255, 255, 255, 0.2);
  padding: 4px 12px;
  border-radius: 6px;
}

.clips-card :deep(.el-card__body) {
  padding: 0;
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-height: 0;
}

.clips-container {
  display: flex;
  flex-direction: column;
  padding: 20px;
  flex: 1;
  overflow: hidden;
  min-height: 0;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px;
  color: #c0c4cc;
}

.empty-state p {
  margin: 12px 0 4px;
  font-size: 16px;
  color: #909399;
}

.empty-state span {
  font-size: 13px;
  color: #c0c4cc;
}

.clips-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  flex: 1;
  overflow-y: auto;
  padding: 4px 8px 0 0;
  min-height: 0;
}

.clips-list::-webkit-scrollbar {
  width: 6px;
}

.clips-list::-webkit-scrollbar-track {
  background: #f0f2f5;
  border-radius: 3px;
}

.clips-list::-webkit-scrollbar-thumb {
  background: #c0c4cc;
  border-radius: 3px;
}

.clips-list::-webkit-scrollbar-thumb:hover {
  background: #909399;
}

.clip-card {
  border-radius: 10px;
  border: 1px solid #e4e7ed;
  background: white;
  transition: transform 0.2s, box-shadow 0.2s, border-color 0.2s;
  flex-shrink: 0;
}

.clip-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(64, 158, 255, 0.12);
  border-color: #409eff;
}

.clip-card :deep(.el-card__body) {
  padding: 14px 18px;
}

.clip-header {
  margin-bottom: 10px;
  padding-bottom: 10px;
  border-bottom: 1px dashed #e4e7ed;
}

.clip-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.clip-type {
  display: flex;
  align-items: center;
  gap: 4px;
  background: linear-gradient(135deg, #e6f7ff 0%, #f0f9ff 100%);
  color: #409eff;
  padding: 5px 12px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 600;
  border: 1px solid #b3d8ff;
}

.clip-time {
  color: #909399;
  font-size: 12px;
}

.clip-content {
  margin-bottom: 12px;
}

.text-content {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  color: #303133;
  font-size: 14px;
  line-height: 1.7;
  padding: 10px 14px;
  background: linear-gradient(135deg, #f9fafc 0%, #f5f7fa 100%);
  border-radius: 8px;
  border: 1px solid #ebeef5;
}

.image-content {
  max-width: 100%;
  border-radius: 8px;
}

.clip-actions {
  display: flex;
  gap: 10px;
  padding-top: 10px;
  border-top: 1px dashed #e4e7ed;
}

.clip-actions :deep(.el-button) {
  display: flex;
  align-items: center;
  gap: 4px;
  border-radius: 8px;
  font-weight: 500;
}

.pagination-wrapper {
  padding: 16px 0 0;
  border-top: 1px solid #e4e7ed;
  display: flex;
  justify-content: center;
  flex-shrink: 0;
}

.pagination-wrapper :deep(.el-pagination) {
  --el-pagination-button-bg-color: #f0f2f5;
}

.pagination-wrapper :deep(.el-pager li.is-active) {
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
  color: white;
}

/* 响应式布局 */
@media (max-width: 900px) {
  .layout-container {
    flex-direction: column;
  }

  .create-section {
    flex: none;
    width: 100%;
    min-width: 100%;
  }

  .clips-section {
    width: 100%;
  }
}

/* 移动端适配 */
@media (max-width: 600px) {
  .header {
    padding: 12px 16px;
  }

  .header h1 {
    font-size: 18px;
  }

  .header-right {
    gap: 8px;
  }

  .user-info,
  .device-info {
    display: none;
  }

  .connection-status {
    padding: 4px 8px;
    font-size: 12px;
  }

  .main {
    padding: 16px;
  }

  .layout-container {
    gap: 16px;
  }

  .create-section {
    flex: none;
    min-width: 100%;
  }

  .create-card :deep(.el-card__body) {
    padding: 16px;
  }

  .create-btn {
    height: 40px;
    font-size: 14px;
  }

  .clips-card {
    max-height: none;
  }

  .clips-container {
    padding: 16px;
  }

  .clip-meta {
    flex-wrap: wrap;
    gap: 4px;
  }

  .clip-type {
    padding: 4px 8px;
    font-size: 12px;
  }

  .clip-time {
    font-size: 11px;
  }

  .text-content {
    padding: 8px 10px;
    font-size: 13px;
    line-height: 1.5;
  }

  .image-thumb {
    max-height: 150px;
  }

  .clip-actions {
    flex-wrap: wrap;
    gap: 8px;
    padding-top: 8px;
  }

  .clip-actions :deep(.el-button) {
    flex: 1;
    min-width: calc(50% - 4px);
    font-size: 12px;
    padding: 6px 12px;
  }

  .preview-dialog :deep(.el-dialog) {
    width: 95% !important;
    margin: 10px auto;
  }

  .preview-header {
    padding: 12px 16px;
    flex-wrap: wrap;
    gap: 8px;
  }

  .preview-header-left {
    font-size: 14px;
  }

  .preview-header-actions {
    flex-wrap: wrap;
    gap: 6px;
  }

  .preview-header-actions :deep(.el-button) {
    font-size: 12px;
    padding: 6px 10px;
  }

  .preview-image {
    max-height: 50vh;
    border-radius: 8px;
  }
}
</style>

<style>
/* 全局样式 - 强制覆盖 Element Plus dialog 背景 */
.preview-dialog.el-dialog {
  background: transparent !important;
  box-shadow: none !important;
}

.preview-dialog .el-dialog__header {
  display: none !important;
}

.preview-dialog .el-dialog__body {
  padding: 0 !important;
  background: transparent !important;
}
</style>