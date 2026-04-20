<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { resetAuthCheck } from '@/router'
import { Hide, View } from '@element-plus/icons-vue'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const deviceName = ref('')
const showPassword = ref(false)
const isLoading = ref(false)

onMounted(() => {
  // 从localStorage读取或生成默认设备名
  const savedDeviceName = localStorage.getItem('deviceName')
  if (savedDeviceName) {
    deviceName.value = savedDeviceName
  } else {
    // 根据平台生成默认名称
    const platform = navigator.platform || 'Unknown'
    const isMobile = /iPhone|iPad|Android/i.test(navigator.userAgent)
    deviceName.value = isMobile ? 'Mobile Device' : platform.includes('Win') ? 'Windows PC' : platform.includes('Mac') ? 'Mac' : 'Device'
  }
})

async function handleLogin() {
  if (!username.value || !password.value) {
    ElMessage.warning('Please enter username and password')
    return
  }

  if (!deviceName.value.trim()) {
    ElMessage.warning('Please enter device name')
    return
  }

  // 保存设备名到localStorage
  localStorage.setItem('deviceName', deviceName.value.trim())

  isLoading.value = true
  try {
    await authStore.login(username.value, password.value)
    resetAuthCheck()
    ElMessage.success('Login successful')
    router.push('/')
  } catch (error: any) {
    ElMessage.error(error)
  } finally {
    isLoading.value = false
  }
}

function goToRegister() {
  router.push('/register')
}
</script>

<template>
  <div class="login-container">
    <div class="login-card">
      <div class="logo-section">
        <div class="logo-icon">
          <svg viewBox="0 0 24 24" width="48" height="48" fill="#409eff">
            <path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-7 14H6v-2h6v2zm4-4H6v-2h10v2zm0-4H6V7h10v2z"/>
          </svg>
        </div>
        <h1>Clipper</h1>
        <p class="subtitle">Welcome back</p>
      </div>

      <el-form @submit.prevent="handleLogin" label-position="top">
        <el-form-item label="Username">
          <el-input
            v-model="username"
            placeholder="Enter your username"
            size="large"
            clearable
          />
        </el-form-item>

        <el-form-item label="Password">
          <el-input
            v-model="password"
            :type="showPassword ? 'text' : 'password'"
            placeholder="Enter your password"
            size="large"
            clearable
          >
            <template #suffix>
              <el-icon
                class="password-icon"
                @click="showPassword = !showPassword"
              >
                <View v-if="showPassword" />
                <Hide v-else />
              </el-icon>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item label="Device Name">
          <el-input
            v-model="deviceName"
            placeholder="Name this device (e.g. iPhone, MacBook)"
            size="large"
            clearable
          />
        </el-form-item>

        <el-button
          type="primary"
          size="large"
          :loading="isLoading"
          @click="handleLogin"
          class="submit-btn"
        >
          {{ isLoading ? 'Signing in...' : 'Sign In' }}
        </el-button>
      </el-form>

      <div class="divider"></div>

      <p class="register-link">
        Don't have an account?
        <a @click="goToRegister">Create one</a>
      </p>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #f0f7ff 0%, #e8f4fc 100%);
  padding: 20px;
}

.login-card {
  padding: 32px 40px 40px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 24px rgba(64, 158, 255, 0.15);
  width: 400px;
  max-width: 100%;
}

.logo-section {
  text-align: center;
  margin-bottom: 32px;
}

.logo-icon {
  margin-bottom: 12px;
}

.logo-section h1 {
  color: #409eff;
  margin: 0 0 8px 0;
  font-size: 28px;
  font-weight: 600;
}

.subtitle {
  color: #909399;
  font-size: 14px;
  margin: 0;
}

.el-form-item {
  margin-bottom: 20px;
}

.el-form-item :deep(.el-form-item__label) {
  color: #303133;
  font-weight: 500;
  font-size: 14px;
  padding-bottom: 8px;
}

.el-input :deep(.el-input__wrapper) {
  border-radius: 8px;
  padding: 4px 12px;
}

.el-input :deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #409eff inset;
}

.el-input :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px #409eff inset;
}

.password-icon {
  cursor: pointer;
  color: #909399;
  transition: color 0.2s;
}

.password-icon:hover {
  color: #409eff;
}

.submit-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  font-weight: 500;
  border-radius: 8px;
  margin-top: 8px;
}

.divider {
  height: 1px;
  background: #e4e7ed;
  margin: 24px 0;
}

.register-link {
  text-align: center;
  color: #909399;
  font-size: 14px;
  margin: 0;
}

.register-link a {
  color: #409eff;
  cursor: pointer;
  font-weight: 500;
  margin-left: 4px;
}

.register-link a:hover {
  text-decoration: underline;
}
</style>