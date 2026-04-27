<template>
  <div class="h5-bind-view">
    <mobile-header title="账号绑定" :show-back="false" />

    <div class="content">
      <div class="bind-card">
        <div class="icon-group">
          <van-icon name="comment-o" size="60" color="#00d100" />
          <van-icon name="exchange" size="30" color="#999" />
          <van-icon name="apps-o" size="60" color="#c96442" />
        </div>
        
        <h3 class="title">飞书智能助手绑定</h3>
        <p class="desc">
          绑定后，您可以直接在飞书中与 EMS 智能助手进行对话，实时查询设备状态与报表。
        </p>

        <div v-if="isLoggedIn" class="user-status">
          <van-cell-group inset>
            <van-cell title="当前登录账号" :value="authStore.user?.name" />
            <van-cell title="飞书 OpenID" :value="shortOpenID" />
          </van-cell-group>

          <van-button
            type="primary"
            block
            round
            :loading="binding"
            class="action-btn"
            @click="handleBind"
          >
            立即绑定
          </van-button>
        </div>

        <div v-else class="login-status">
          <p class="login-tip">请先登录 EMS 系统以完成绑定</p>
          <van-button
            type="primary"
            block
            round
            plain
            class="action-btn"
            @click="goToLogin"
          >
            去登录
          </van-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showToast, showDialog } from 'vant'
import { useAuthStore } from '@/stores/auth'
import { authApi } from '@/api/auth'
import MobileHeader from '@/components/mobile/MobileHeader.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const openID = computed(() => (route.query.openid as string) || '')
const shortOpenID = computed(() => {
  if (!openID.value) return '未知'
  return openID.value.substring(0, 8) + '...'
})

const isLoggedIn = computed(() => authStore.isLoggedIn)
const binding = ref(false)

const handleBind = async () => {
  if (!openID.value) {
    showToast('无效的绑定请求：缺少 OpenID')
    return
  }

  binding.value = true
  try {
    // 使用统一的 authApi 进行绑定
    await authApi.bindLark(openID.value)
    
    await showDialog({
      title: '绑定成功',
      message: '您的飞书账号已成功关联 EMS 系统，现在可以返回飞书进行对话了。'
    })
    
    router.push('/h5')
  } catch (error: any) {
    showToast(error.response?.data?.error || '绑定失败')
  } finally {
    binding.value = false
  }
}

const goToLogin = () => {
  router.push({
    path: '/login',
    query: { redirect: route.fullPath }
  })
}

onMounted(() => {
  if (!openID.value) {
    showToast('参数错误')
  }
})
</script>

<style scoped>
.h5-bind-view {
  min-height: 100vh;
  background: var(--color-bg-primary);
  padding-top: 46px;
}

.content {
  padding: 40px 20px;
}

.bind-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-high);
  padding: 30px 20px;
  text-align: center;
  box-shadow: var(--shadow-lg);
  border: 1px solid var(--color-border);
}

.icon-group {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 15px;
  margin-bottom: 30px;
}

.title {
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 12px;
  font-family: var(--font-serif);
}

.desc {
  font-size: 14px;
  color: var(--color-text-secondary);
  line-height: 1.6;
  margin-bottom: 40px;
}

.user-status, .login-status {
  margin-top: 20px;
}

.login-tip {
  color: var(--color-warning);
  font-size: 14px;
  margin-bottom: 20px;
}

.action-btn {
  margin-top: 30px;
  height: 50px;
  font-weight: 600;
}
</style>
