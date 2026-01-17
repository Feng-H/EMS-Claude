<template>
  <div class="change-password-container">
    <!-- 背景动画效果 -->
    <div class="bg-effects">
      <div class="grid-overlay"></div>
      <div class="floating-orbs">
        <div class="orb orb-1"></div>
        <div class="orb orb-2"></div>
        <div class="orb orb-3"></div>
      </div>
    </div>

    <!-- 修改密码卡片 -->
    <div class="change-password-wrapper animate-fade-in">
      <div class="change-password-card">
        <div class="card-header">
          <div class="icon-wrapper">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 48px; height: 48px;">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
              <path d="M7 11V7a5 5 0 0 1 10 0v4" />
            </svg>
          </div>
          <h2>修改密码</h2>
          <p>首次登录需要修改默认密码</p>
        </div>

        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          @submit.prevent="handleSubmit"
          class="password-form"
        >
          <el-form-item prop="oldPassword">
            <div class="input-group">
              <div class="input-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4" />
                </svg>
              </div>
              <el-input
                v-model="form.oldPassword"
                type="password"
                placeholder="当前密码"
                size="large"
                show-password
              />
            </div>
          </el-form-item>

          <el-form-item prop="newPassword">
            <div class="input-group">
              <div class="input-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/>
                </svg>
              </div>
              <el-input
                v-model="form.newPassword"
                type="password"
                placeholder="新密码（至少6位）"
                size="large"
                show-password
              />
            </div>
          </el-form-item>

          <el-form-item prop="confirmPassword">
            <div class="input-group">
              <div class="input-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <path d="M9 12l2 2 4-4m6 2a9 9 0 1 1-18 0 9 9 0 0 1 18 0z"/>
                </svg>
              </div>
              <el-input
                v-model="form.confirmPassword"
                type="password"
                placeholder="确认新密码"
                size="large"
                show-password
                @keyup.enter="handleSubmit"
              />
            </div>
          </el-form-item>

          <el-form-item>
            <el-button
              type="primary"
              size="large"
              :loading="loading"
              native-type="submit"
              class="submit-btn"
            >
              <span v-if="!loading">修改密码</span>
              <span v-else>提交中...</span>
            </el-button>
          </el-form-item>
        </el-form>

        <div class="card-footer">
          <div class="password-tips">
            <p>密码要求：</p>
            <ul>
              <li>长度至少 6 个字符</li>
              <li>建议包含大小写字母、数字和特殊字符</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '@/api/auth'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'

const router = useRouter()

const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const validateConfirmPassword = (rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('请再次输入新密码'))
  } else if (value !== form.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules: FormRules = {
  oldPassword: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 个字符', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' },
  ],
}

async function handleSubmit() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      await authApi.changePassword({
        old_password: form.oldPassword,
        new_password: form.newPassword,
      })
      ElMessage.success('密码修改成功，请重新登录')
      // 跳转到首页
      router.push('/dashboard')
    } catch (error) {
      // Error handled by request interceptor
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.change-password-container {
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

/* 修改密码卡片 */
.change-password-wrapper {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 450px;
  padding: 20px;
}

.change-password-card {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: 24px;
  padding: 40px;
  position: relative;
  overflow: hidden;
}

.change-password-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--color-primary), var(--color-success));
}

.card-header {
  text-align: center;
  margin-bottom: 32px;
}

.icon-wrapper {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, rgba(0, 212, 255, 0.1), rgba(0, 255, 163, 0.1));
  border: 2px solid var(--color-primary-dim);
  border-radius: 20px;
  color: var(--color-primary);
  margin-bottom: 20px;
}

.card-header h2 {
  font-size: 24px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 8px 0;
}

.card-header p {
  font-size: 14px;
  color: var(--color-text-tertiary);
  margin: 0;
}

.password-form {
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

.submit-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  background: linear-gradient(135deg, var(--color-primary), #00B8E4);
  border: none;
  color: var(--color-bg-primary);
  transition: all var(--transition-base);
}

.submit-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px var(--color-primary-glow);
}

.card-footer {
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid var(--color-border);
}

.password-tips {
  text-align: left;
}

.password-tips p {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-secondary);
  margin: 0 0 8px 0;
}

.password-tips ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.password-tips li {
  font-size: 12px;
  color: var(--color-text-tertiary);
  padding-left: 16px;
  position: relative;
  margin-bottom: 4px;
}

.password-tips li::before {
  content: '•';
  position: absolute;
  left: 0;
  color: var(--color-primary);
}

/* 响应式 */
@media (max-width: 480px) {
  .change-password-wrapper {
    padding: 16px;
  }

  .change-password-card {
    padding: 24px;
    border-radius: 16px;
  }

  .orb-1, .orb-2, .orb-3 {
    display: none;
  }
}
</style>
