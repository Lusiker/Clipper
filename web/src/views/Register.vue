<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { resetAuthCheck } from '@/router'
import { Hide, View } from '@element-plus/icons-vue'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const showPassword = ref(false)
const showConfirmPassword = ref(false)
const isLoading = ref(false)
const usernameError = ref(false)

const passwordMismatch = computed(() => {
  if (confirmPassword.value === '') return false
  return password.value !== confirmPassword.value
})

async function handleRegister() {
  usernameError.value = false

  if (!username.value || !password.value || !confirmPassword.value) {
    ElMessage.warning('Please fill all fields')
    return
  }

  if (username.value.length < 3) {
    ElMessage.warning('Username must be at least 3 characters')
    return
  }

  if (password.value.length < 6) {
    ElMessage.warning('Password must be at least 6 characters')
    return
  }

  if (password.value !== confirmPassword.value) {
    ElMessage.warning('Passwords do not match')
    return
  }

  isLoading.value = true
  try {
    await authStore.register(username.value, password.value)
    resetAuthCheck()
    ElMessage.success('Registration successful')
    router.push('/')
  } catch (error: any) {
    if (error.includes('username already exists')) {
      usernameError.value = true
      ElMessage.error('Username already taken')
    } else {
      ElMessage.error(error)
    }
  } finally {
    isLoading.value = false
  }
}

function goToLogin() {
  router.push('/login')
}
</script>

<template>
  <div class="register-container">
    <div class="register-card">
      <div class="logo-section">
        <div class="logo-icon">
          <svg viewBox="0 0 24 24" width="48" height="48" fill="#409eff">
            <path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-7 14H6v-2h6v2zm4-4H6v-2h10v2zm0-4H6V7h10v2z"/>
          </svg>
        </div>
        <h1>Clipper</h1>
        <p class="subtitle">Create your account</p>
      </div>

      <el-form @submit.prevent="handleRegister" label-position="top">
        <el-form-item label="Username" :class="{ 'has-error': usernameError }">
          <el-input
            v-model="username"
            placeholder="Enter your username (min 3 characters)"
            size="large"
            clearable
            :status="usernameError ? 'error' : ''"
          />
          <div v-if="usernameError" class="error-message">
            Username already taken
          </div>
        </el-form-item>

        <el-form-item label="Password">
          <el-input
            v-model="password"
            :type="showPassword ? 'text' : 'password'"
            placeholder="Enter password (min 6 characters)"
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

        <el-form-item label="Confirm Password" :class="{ 'has-error': passwordMismatch }">
          <el-input
            v-model="confirmPassword"
            :type="showConfirmPassword ? 'text' : 'password'"
            placeholder="Confirm your password"
            size="large"
            :status="passwordMismatch ? 'error' : ''"
          >
            <template #suffix>
              <el-icon
                class="password-icon"
                @click="showConfirmPassword = !showConfirmPassword"
              >
                <View v-if="showConfirmPassword" />
                <Hide v-else />
              </el-icon>
            </template>
          </el-input>
          <div v-if="passwordMismatch" class="error-message">
            Passwords do not match
          </div>
        </el-form-item>

        <el-button
          type="primary"
          size="large"
          :loading="isLoading"
          @click="handleRegister"
          class="submit-btn"
        >
          {{ isLoading ? 'Creating Account...' : 'Create Account' }}
        </el-button>
      </el-form>

      <div class="divider"></div>

      <p class="login-link">
        Already have an account?
        <a @click="goToLogin">Sign in</a>
      </p>
    </div>
  </div>
</template>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #f0f7ff 0%, #e8f4fc 100%);
  padding: 20px;
}

.register-card {
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

.has-error :deep(.el-input__wrapper) {
  box-shadow: 0 0 0 1px #f56c6c inset !important;
}

.error-message {
  color: #f56c6c;
  font-size: 12px;
  margin-top: 4px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.error-message::before {
  content: '⚠';
}

.submit-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  font-weight: 500;
  border-radius: 8px;
  margin-top: 8px;
}

.submit-btn :deep(.el-button) {
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
  border: none;
}

.divider {
  height: 1px;
  background: #e4e7ed;
  margin: 24px 0;
}

.login-link {
  text-align: center;
  color: #909399;
  font-size: 14px;
  margin: 0;
}

.login-link a {
  color: #409eff;
  cursor: pointer;
  font-weight: 500;
  margin-left: 4px;
}

.login-link a:hover {
  text-decoration: underline;
}
</style>