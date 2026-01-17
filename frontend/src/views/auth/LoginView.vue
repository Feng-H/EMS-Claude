<template>
  <div class="login-container">
    <!-- 背景动画效果 -->
    <div class="bg-effects">
      <div class="grid-overlay"></div>
      <div class="floating-orbs">
        <div class="orb orb-1"></div>
        <div class="orb orb-2"></div>
        <div class="orb orb-3"></div>
      </div>
      <div class="scanline"></div>
    </div>

    <!-- 登录卡片 -->
    <div class="login-wrapper animate-fade-in">
      <!-- 左侧：品牌信息 -->
      <div class="brand-section">
        <div class="brand-logo">
          <div class="logo-icon">
            <svg viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M32 8L12 20L8 32L12 44L32 52L52 44L56 32L52 32L32 8Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M8 32L16 36M56 32L48 36" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <circle cx="32" cy="28" r="4" fill="currentColor"/>
              <circle cx="32" cy="44" r="4" fill="currentColor"/>
            </svg>
          </div>
          <h1 class="brand-title">
            <span class="title-main">EMS</span>
            <span class="title-sub">设备管理系统</span>
          </h1>
        </div>
        <div class="brand-info">
          <div class="info-item">
            <span class="info-label">Equipment</span>
            <span class="info-value">Management System</span>
          </div>
          <div class="info-item">
            <span class="info-label">智能监控</span>
            <span class="info-value">实时管理 · 数据驱动</span>
          </div>
        </div>
        <div class="brand-stats">
          <div class="stat-item">
            <div class="stat-value">50K+</div>
            <div class="stat-label">设备数量</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">99.9%</div>
            <div class="stat-label">正常运行</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">24/7</div>
            <div class="stat-label">实时监控</div>
          </div>
        </div>
      </div>

      <!-- 右侧：登录/申请表单 -->
      <div class="form-section">
        <div class="form-card">
          <div class="form-header">
            <h2>{{ isLoginMode ? '欢迎回来' : '申请账号' }}</h2>
            <p>{{ isLoginMode ? '请登录您的账号以继续' : '填写信息申请系统账号' }}</p>
          </div>

          <!-- 登录表单 -->
          <el-form
            v-if="isLoginMode"
            ref="formRef"
            :model="loginForm"
            :rules="loginRules"
            @submit.prevent="handleLogin"
            class="login-form"
          >
            <el-form-item prop="username">
              <div class="input-group">
                <div class="input-icon">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path d="M20 21v-2a4 4 0 0 0-4-4V8a4 4 0 0 0-4-4h-5.586a4 4 0 0 0-2.828 1.172l-2.414 2.414a4 4 0 0 0-5.656 0 4 4 0 0 0 0 5.656l2.414 2.414a4 4 0 0 0 1.172 1.172V8a4 4 0 0 0 4 4v2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </div>
                <el-input
                  v-model="loginForm.username"
                  placeholder="用户名"
                  size="large"
                  @keyup.enter="handleLogin"
                />
              </div>
            </el-form-item>

            <el-form-item prop="password">
              <div class="input-group">
                <div class="input-icon">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                    <path d="M7 11V7a5 5 0 0 1 10 0v4" />
                  </svg>
                </div>
                <el-input
                  v-model="loginForm.password"
                  type="password"
                  placeholder="密码"
                  size="large"
                  show-password
                  @keyup.enter="handleLogin"
                />
              </div>
            </el-form-item>

            <el-form-item>
              <el-button
                type="primary"
                size="large"
                :loading="loading"
                native-type="submit"
                class="login-btn"
              >
                <span v-if="!loading">登录</span>
                <span v-else>登录中...</span>
              </el-button>
            </el-form-item>
          </el-form>

          <!-- 申请表单 -->
          <el-form
            v-else
            ref="formRef"
            :model="applyForm"
            :rules="applyRules"
            @submit.prevent="handleApply"
            class="apply-form"
          >
            <el-form-item prop="username">
              <div class="input-group">
                <div class="input-icon">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path d="M20 21v-2a4 4 0 0 0-4-4V8a4 4 0 0 0-4-4h-5.586a4 4 0 0 0-2.828 1.172l-2.414 2.414a4 4 0 0 0-5.656 0 4 4 0 0 0 0 5.656l2.414 2.414a4 4 0 0 0 1.172 1.172V8a4 4 0 0 0 4 4v2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </div>
                <el-input
                  v-model="applyForm.username"
                  placeholder="用户名（3-50个字符）"
                  size="large"
                />
              </div>
            </el-form-item>

            <el-form-item prop="name">
              <div class="input-group">
                <div class="input-icon">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                    <circle cx="12" cy="7" r="4"/>
                  </svg>
                </div>
                <el-input
                  v-model="applyForm.name"
                  placeholder="真实姓名"
                  size="large"
                />
              </div>
            </el-form-item>

            <el-form-item prop="role">
              <div class="input-group">
                <div class="input-icon">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                  </svg>
                </div>
                <el-select
                  v-model="applyForm.role"
                  placeholder="选择角色"
                  size="large"
                  style="width: 100%;"
                >
                  <el-option label="操作工" value="operator" />
                  <el-option label="维修工" value="maintenance" />
                  <el-option label="设备工程师" value="engineer" />
                  <el-option label="设备主管" value="supervisor" />
                </el-select>
              </div>
            </el-form-item>

            <el-form-item prop="phone">
              <div class="input-group">
                <div class="input-icon">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/>
                  </svg>
                </div>
                <el-input
                  v-model="applyForm.phone"
                  placeholder="联系电话（可选）"
                  size="large"
                />
              </div>
            </el-form-item>

            <el-form-item prop="password">
              <div class="input-group">
                <div class="input-icon">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                    <path d="M7 11V7a5 5 0 0 1 10 0v4" />
                  </svg>
                </div>
                <el-input
                  v-model="applyForm.password"
                  type="password"
                  placeholder="初始密码（至少6位）"
                  size="large"
                  show-password
                />
              </div>
            </el-form-item>

            <el-form-item>
              <el-button
                type="primary"
                size="large"
                :loading="loading"
                native-type="submit"
                class="login-btn"
              >
                <span v-if="!loading">提交申请</span>
                <span v-else>提交中...</span>
              </el-button>
            </el-form-item>
          </el-form>

          <div class="form-footer">
            <div class="divider">
              <span>{{ isLoginMode ? '测试环境' : '申请提示' }}</span>
            </div>
            <div v-if="isLoginMode" class="test-account">
              <span>账号: <strong>admin</strong></span>
              <span>密码: <strong>password123</strong></span>
            </div>
            <div v-else class="test-account">
              <span>提交后需要等待管理员审核</span>
            </div>
            <div class="apply-link">
              <span @click="toggleMode">{{ isLoginMode ? '没有账号？立即申请' : '已有账号？返回登录' }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { authApi, type ApplyAccountRequest } from '@/api/auth'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const formRef = ref<FormInstance>()
const loading = ref(false)
const isLoginMode = ref(true)

// 登录表单
const loginForm = reactive({
  username: '',
  password: '',
})

// 申请表单
const applyForm = reactive({
  username: '',
  password: '',
  name: '',
  role: 'operator',
  phone: '',
})

const loginRules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

const applyRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '用户名长度在 3 到 50 个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 个字符', trigger: 'blur' },
  ],
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' },
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' },
  ],
}

const form = loginForm
const rules = loginRules

function toggleMode() {
  isLoginMode.value = !isLoginMode.value
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

async function handleLogin() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      const response = await authStore.login(loginForm)
      ElMessage.success('登录成功')

      // 检查是否需要修改密码
      if (response.must_change_password) {
        router.push('/change-password')
        return
      }

      const redirect = (route.query.redirect as string) || '/dashboard'
      router.push(redirect)
    } catch (error) {
      // Error is handled by request interceptor
    } finally {
      loading.value = false
    }
  })
}

async function handleApply() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      const data: ApplyAccountRequest = {
        username: applyForm.username,
        password: applyForm.password,
        name: applyForm.name,
        role: applyForm.role,
        phone: applyForm.phone || undefined,
      }
      await authApi.applyAccount(data)
      ElMessage.success('申请已提交，请等待管理员审核')
      // 重置表单并切换回登录模式
      Object.assign(applyForm, {
        username: '',
        password: '',
        name: '',
        role: 'operator',
        phone: '',
      })
      toggleMode()
    } catch (error) {
      // Error is handled by request interceptor
    } finally {
      loading.value = false
    }
  })
}

onMounted(() => {
  // 预填充测试账号（开发环境）
  if (import.meta.env.DEV) {
    loginForm.username = 'admin'
    loginForm.password = 'password123'
  }
})
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: var(--color-bg-primary);
}

/* 背景效果 */
.bg-effects {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.grid-overlay {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(0, 212, 255, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(0, 212, 255, 0.03) 1px, transparent 1px);
  background-size: 60px 60px;
  mask-image: radial-gradient(ellipse 80% 50% at 50% 50%, black, transparent);
  -webkit-mask-image: radial-gradient(ellipse 80% 50% at 50% 50%, black, transparent);
}

.floating-orbs {
  position: absolute;
  inset: 0;
  overflow: hidden;
}

.orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  opacity: 0.4;
  animation: float 8s ease-in-out infinite;
}

.orb-1 {
  width: 400px;
  height: 400px;
  background: var(--color-primary);
  top: -200px;
  right: -100px;
  animation-delay: 0s;
}

.orb-2 {
  width: 300px;
  height: 300px;
  background: var(--color-success);
  bottom: -150px;
  left: -50px;
  animation-delay: -3s;
}

.orb-3 {
  width: 200px;
  height: 200px;
  background: var(--color-warning);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation-delay: -5s;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
   33% {
    transform: translate(30px, -50px) scale(1.1);
  }
  66% {
    transform: translate(-20px, 20px) scale(0.9);
  }
}

.scanline {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, var(--color-primary), transparent);
  opacity: 0.5;
  animation: scan 4s linear infinite;
}

@keyframes scan {
  0% {
    top: 0;
    opacity: 0;
  }
  10% {
    opacity: 1;
  }
  90% {
    opacity: 1;
  }
  100% {
    top: 100%;
    opacity: 0;
  }
}

/* 登录主体 */
.login-wrapper {
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: 60px;
  align-items: center;
  max-width: 1100px;
  width: 100%;
  padding: 20px;
  position: relative;
  z-index: 1;
}

/* 品牌部分 */
.brand-section {
  display: flex;
  flex-direction: column;
  gap: 40px;
}

.brand-logo {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.logo-icon {
  width: 120px;
  height: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, rgba(0, 212, 255, 0.1), rgba(0, 255, 163, 0.1));
  border: 2px solid var(--color-primary-dim);
  border-radius: 24px;
  color: var(--color-primary);
  animation: pulse-glow 3s ease-in-out infinite;
}

.logo-icon svg {
  width: 64px;
  height: 64px;
}

.brand-title {
  text-align: center;
}

.title-main {
  display: block;
  font-size: 48px;
  font-weight: 800;
  letter-spacing: 8px;
  background: linear-gradient(135deg, var(--color-primary), var(--color-success));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 4px;
}

.title-sub {
  display: block;
  font-size: 18px;
  color: var(--color-text-secondary);
  font-weight: 400;
  letter-spacing: 4px;
}

.brand-info {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: 12px 20px;
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: 8px;
}

.info-label {
  color: var(--color-text-tertiary);
  font-size: 12px;
  letter-spacing: 2px;
  text-transform: uppercase;
}

.info-value {
  color: var(--color-text-secondary);
  font-size: 14px;
}

.brand-stats {
  display: flex;
  justify-content: space-around;
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-family: var(--font-numbers);
  font-size: 24px;
  font-weight: 700;
  color: var(--color-primary);
  line-height: 1.2;
}

.stat-label {
  font-size: 12px;
  color: var(--color-text-tertiary);
  margin-top: 4px;
}

/* 表单部分 */
.form-section {
  width: 100%;
}

.form-card {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: 24px;
  padding: 40px;
  position: relative;
  overflow: hidden;
}

.form-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--color-primary), var(--color-success));
}

.form-header {
  text-align: center;
  margin-bottom: 32px;
}

.form-header h2 {
  font-size: 24px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 8px;
}

.form-header p {
  font-size: 14px;
  color: var(--color-text-tertiary);
}

.login-form {
  margin-bottom: 24px;
}

.input-group {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 16px;
  color: var(--color-text-tertiary);
  z-index: 1;
  pointer-events: none;
}

.input-group .el-input {
  flex: 1;
}

.input-group :deep(.el-input__wrapper) {
  padding-left: 48px;
  background: var(--color-bg-secondary);
  border-color: var(--color-border);
}

.input-group :deep(.el-input__wrapper:hover) {
  border-color: var(--color-primary-dim);
}

.input-group :deep(.el-input__wrapper.is-focus) {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px var(--color-primary-dim);
}

.login-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  background: linear-gradient(135deg, var(--color-primary), #00B8E4);
  border: none;
  color: var(--color-bg-primary);
  transition: all var(--transition-base);
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px var(--color-primary-glow);
}

.form-footer {
  margin-top: 24px;
}

.divider {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-bottom: 16px;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--color-divider);
}

.divider span {
  font-size: 12px;
  color: var(--color-text-tertiary);
  padding: 0 8px;
}

.test-account {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: var(--color-bg-secondary);
  border-radius: 8px;
  font-size: 13px;
  color: var(--color-text-secondary);
}

.test-account strong {
  color: var(--color-primary);
}

.apply-link {
  margin-top: 12px;
  text-align: center;
}

.apply-link span {
  font-size: 13px;
  color: var(--color-primary);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.apply-link span:hover {
  color: var(--color-success);
  text-decoration: underline;
}

/* 响应式 */
@media (max-width: 968px) {
  .login-wrapper {
    grid-template-columns: 1fr;
    gap: 40px;
  }

  .brand-section {
    text-align: center;
  }

  .brand-info {
    align-items: center;
  }

  .info-item {
    flex-direction: column;
    gap: 4px;
  }

  .brand-stats {
    justify-content: center;
    gap: 24px;
  }
}

@media (max-width: 480px) {
  .login-wrapper {
    padding: 16px;
    gap: 24px;
  }

  .form-card {
    padding: 24px;
    border-radius: 16px;
  }

  .title-main {
    font-size: 36px;
  }

  .orb-1, .orb-2, .orb-3 {
    display: none;
  }
}
</style>
