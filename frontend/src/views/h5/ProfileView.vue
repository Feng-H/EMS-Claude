<template>
  <div class="h5-profile-view">
    <mobile-header title="个人中心" :show-back="false" />

    <div class="content">
      <!-- 用户基本信息 -->
      <div class="profile-header">
        <div class="avatar-wrapper">
          <van-icon name="user-circle-o" size="80" color="#faf9f5" />
        </div>
        <div class="user-meta">
          <div class="user-name">{{ authStore.user?.name || '未登录' }}</div>
          <div class="user-role-badge">
            <van-tag round type="primary" color="rgba(250, 249, 245, 0.2)">
              {{ getRoleText(authStore.user?.role) }}
            </van-tag>
          </div>
        </div>
      </div>

      <!-- 功能列表 -->
      <van-cell-group inset class="menu-list">
        <van-cell title="所属工厂" :value="authStore.user?.factory_name || '未分配'" icon="hotel-o" />
        <van-cell title="联系电话" :value="authStore.user?.phone || '未填写'" icon="phone-o" />
      </van-cell-group>

      <van-cell-group inset class="menu-list">
        <van-cell title="修改密码" is-link to="/change-password" icon="shield-o" />
        <van-cell title="帮助与反馈" is-link icon="question-o" />
        <van-cell title="关于系统" is-link icon="info-o" />
      </van-cell-group>

      <!-- 退出登录按钮 -->
      <div class="logout-wrapper">
        <van-button
          type="danger"
          block
          round
          plain
          class="logout-btn"
          @click="handleLogout"
        >
          退出当前账号
        </van-button>
      </div>

      <div class="version">Version 1.0.0 (Demo)</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { showConfirmDialog, showToast } from 'vant'
import { useAuthStore } from '@/stores/auth'
import MobileHeader from '@/components/mobile/MobileHeader.vue'

const router = useRouter()
const authStore = useAuthStore()

const getRoleText = (role?: string) => {
  const roles: Record<string, string> = {
    admin: '系统管理员',
    supervisor: '设备主管',
    engineer: '设备工程师',
    maintenance: '维修工',
    operator: '操作工'
  }
  return roles[role || ''] || '未知角色'
}

const handleLogout = async () => {
  try {
    await showConfirmDialog({
      title: '退出确认',
      message: '确定要退出登录并返回登录页面吗？'
    })
    authStore.logout()
    showToast('已安全退出')
    router.push('/login')
  } catch {
    // 用户取消
  }
}
</script>

<style scoped>
.h5-profile-view {
  min-height: 100vh;
  background: var(--color-bg-primary);
  padding-top: 46px;
  padding-bottom: 60px;
}

.content {
  padding-top: 0;
}

.profile-header {
  background: var(--color-terracotta);
  padding: 40px 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: var(--space-xl);
  border-radius: 0 0 30px 30px;
  box-shadow: 0 4px 12px rgba(201, 100, 66, 0.2);
}

.avatar-wrapper {
  margin-bottom: var(--space-md);
  filter: drop-shadow(0 4px 8px rgba(0, 0, 0, 0.1));
}

.user-meta {
  text-align: center;
}

.user-name {
  font-size: 22px;
  font-weight: 600;
  color: #faf9f5;
  margin-bottom: 8px;
  font-family: var(--font-serif);
}

.user-role-badge {
  display: inline-block;
}

.menu-list {
  margin-top: var(--space-lg);
  box-shadow: var(--shadow-sm);
}

.logout-wrapper {
  margin: 40px 20px 20px;
}

.logout-btn {
  height: 48px;
  border-color: var(--color-danger);
  color: var(--color-danger);
  font-weight: 500;
}

.version {
  text-align: center;
  font-size: 12px;
  color: var(--color-text-muted);
  margin-top: 20px;
}

/* 暗色模式 */
@media (prefers-color-scheme: dark) {
  .profile-header {
    background: var(--color-terracotta);
    opacity: 0.9;
  }
}
</style>
