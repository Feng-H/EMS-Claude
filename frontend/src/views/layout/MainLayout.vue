<template>
  <div class="main-layout">
    <!-- 侧边栏 -->
    <aside class="sidebar" :class="{ collapsed: sidebarCollapsed }">
      <!-- Logo 区域 -->
      <div class="sidebar-header">
        <div class="logo-container">
          <div class="logo-icon">
            <svg viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M32 8L12 20L8 32L12 44L32 52L52 44L56 32L52 32L32 8Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M8 32L16 36M56 32L48 36" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <circle cx="32" cy="28" r="4" fill="currentColor"/>
              <circle cx="32" cy="44" r="4" fill="currentColor"/>
            </svg>
          </div>
          <transition name="logo-text">
            <span v-show="!sidebarCollapsed" class="logo-text">EMS</span>
          </transition>
        </div>
        <!-- 折叠按钮 - 固定在右上角 -->
        <button class="sidebar-collapse-btn" @click="toggleSidebar" title="折叠/展开侧边栏">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 14px; height: 14px;">
            <polyline points="15 18 9 12 15 6"/>
          </svg>
        </button>
      </div>

      <!-- 导航菜单 -->
      <nav class="sidebar-nav">
        <div class="nav-section">
          <div v-if="!sidebarCollapsed" class="nav-section-title">基础功能</div>
          <router-link to="/dashboard" class="nav-item" :class="{ active: isActive('/dashboard') }">
            <div class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <rect x="3" y="3" width="7" height="7" rx="1"/>
                <rect x="14" y="3" width="7" height="7" rx="1"/>
                <rect x="3" y="14" width="7" height="7" rx="1"/>
                <rect x="14" y="14" width="7" height="7" rx="1"/>
              </svg>
            </div>
            <transition name="nav-text">
              <span v-show="!sidebarCollapsed" class="nav-text">首页</span>
            </transition>
            <div v-if="!sidebarCollapsed" class="nav-indicator"></div>
          </router-link>

          <router-link to="/equipment" class="nav-item" :class="{ active: isActive('/equipment') }">
            <div class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
                <polyline points="3.27 6.96 12 12.01 20.73 6.96"/>
                <line x1="12" y1="22.08" x2="12" y2="12"/>
              </svg>
            </div>
            <transition name="nav-text">
              <span v-show="!sidebarCollapsed" class="nav-text">设备台账</span>
            </transition>
            <div v-if="!sidebarCollapsed" class="nav-indicator"></div>
          </router-link>

          <div v-if="canManageOrg" class="nav-group" :class="{ active: isGroupActive('/organization') }">
            <button class="nav-item nav-group-toggle" @click="toggleGroup('organization')">
              <div class="nav-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <path d="M3 21h18"/>
                  <rect x="5" y="3" width="4" height="14" rx="1"/>
                  <rect x="15" y="8" width="4" height="9" rx="1"/>
                  <path d="M10 21v-4"/>
                </svg>
              </div>
              <transition name="nav-text">
                <span v-show="!sidebarCollapsed" class="nav-text">组织架构</span>
              </transition>
              <transition name="arrow">
                <svg v-show="!sidebarCollapsed" class="nav-arrow" :class="{ expanded: expandedGroups.organization }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="6 9 12 15 18 9"/>
                </svg>
              </transition>
            </button>
            <transition name="submenu">
              <div v-show="sidebarCollapsed || expandedGroups.organization" class="nav-submenu" :class="{ 'always-show': sidebarCollapsed }">
                <router-link to="/organization" class="nav-subitem" :class="{ active: isActive('/organization') && !isActive('/organization/users') }">
                  <span class="submenu-dot"></span>
                  <span class="submenu-text">组织管理</span>
                </router-link>
                <router-link to="/organization/users" class="nav-subitem" :class="{ active: isActive('/organization/users') }">
                  <span class="submenu-dot"></span>
                  <span class="submenu-text">人员管理</span>
                </router-link>
              </div>
            </transition>
          </div>

          <router-link v-if="!canManageOrg" to="/organization" class="nav-item" :class="{ active: isActive('/organization') }">
            <div class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M3 21h18"/>
                <rect x="5" y="3" width="4" height="14" rx="1"/>
                <rect x="15" y="8" width="4" height="9" rx="1"/>
                <path d="M10 21v-4"/>
              </svg>
            </div>
            <transition name="nav-text">
              <span v-show="!sidebarCollapsed" class="nav-text">组织架构</span>
            </transition>
            <div v-if="!sidebarCollapsed" class="nav-indicator"></div>
          </router-link>
        </div>

        <div class="nav-section">
          <div v-if="!sidebarCollapsed" class="nav-section-title">设备维护</div>

          <!-- 点检管理 -->
          <div class="nav-group" :class="{ active: isGroupActive('/inspection') }">
            <button class="nav-item nav-group-toggle" @click="toggleGroup('inspection')">
              <div class="nav-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                  <polyline points="22 4 12 14.01 9 11.01"/>
                </svg>
              </div>
              <transition name="nav-text">
                <span v-show="!sidebarCollapsed" class="nav-text">点检管理</span>
              </transition>
              <transition name="arrow">
                <svg v-show="!sidebarCollapsed" class="nav-arrow" :class="{ expanded: expandedGroups.inspection }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="6 9 12 15 18 9"/>
                </svg>
              </transition>
            </button>
            <transition name="submenu">
              <div v-show="sidebarCollapsed || expandedGroups.inspection" class="nav-submenu" :class="{ 'always-show': sidebarCollapsed }">
                <router-link to="/inspection/tasks" class="nav-subitem" :class="{ active: isActive('/inspection/tasks') }">
                  <span class="submenu-dot"></span>
                  <span class="submenu-text">点检任务</span>
                </router-link>
                <router-link v-if="canManageOrg" to="/inspection/templates" class="nav-subitem" :class="{ active: isActive('/inspection/templates') }">
                  <span class="submenu-dot"></span>
                  <span class="submenu-text">模板配置</span>
                </router-link>
                <router-link to="/inspection/execute" class="nav-subitem" :class="{ active: isActive('/inspection/execute') }">
                  <span class="submenu-dot"></span>
                  <span class="submenu-text">移动端点检</span>
                </router-link>
              </div>
            </transition>
          </div>

          <!-- 维修管理 -->
          <div class="nav-group" :class="{ active: isGroupActive('/repair') }">
            <button class="nav-item nav-group-toggle" @click="toggleGroup('repair')">
              <div class="nav-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/>
                </svg>
              </div>
              <transition name="nav-text">
                <span v-show="!sidebarCollapsed" class="nav-text">维修管理</span>
              </transition>
              <transition name="arrow">
                <svg v-show="!sidebarCollapsed" class="nav-arrow" :class="{ expanded: expandedGroups.repair }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="6 9 12 15 18 9"/>
                </svg>
              </transition>
            </button>
            <transition name="submenu">
              <div v-show="sidebarCollapsed || expandedGroups.repair" class="nav-submenu" :class="{ 'always-show': sidebarCollapsed }">
                <router-link to="/repair/orders" class="nav-subitem" :class="{ active: isActive('/repair/orders') }">
                  <span class="submenu-dot"></span>
                  <span class="submenu-text">维修工单</span>
                </router-link>
                <router-link to="/repair/create" class="nav-subitem" :class="{ active: isActive('/repair/create') }">
                  <span class="submenu-dot"></span>
                  <span class="submenu-text">报修申请</span>
                </router-link>
              </div>
            </transition>
          </div>

          <!-- 保养管理 -->
          <div class="nav-group" :class="{ active: isGroupActive('/maintenance') }">
            <button class="nav-item nav-group-toggle" @click="toggleGroup('maintenance')">
              <div class="nav-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                  <line x1="16" y1="2" x2="16" y2="6"/>
                  <line x1="8" y1="2" x2="8" y2="6"/>
                  <line x1="3" y1="10" x2="21" y2="10"/>
                </svg>
              </div>
              <transition name="nav-text">
                <span v-show="!sidebarCollapsed" class="nav-text">保养管理</span>
              </transition>
              <transition name="arrow">
                <svg v-show="!sidebarCollapsed" class="nav-arrow" :class="{ expanded: expandedGroups.maintenance }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="6 9 12 15 18 9"/>
                </svg>
              </transition>
            </button>
            <transition name="submenu">
              <div v-show="sidebarCollapsed || expandedGroups.maintenance" class="nav-submenu" :class="{ 'always-show': sidebarCollapsed }">
                <router-link v-if="canManageOrg" to="/maintenance/plans" class="nav-subitem" :class="{ active: isActive('/maintenance/plans') }">
                  <span class="submenu-dot"></span>
                  <span class="submenu-text">计划配置</span>
                </router-link>
                <router-link to="/maintenance/tasks" class="nav-subitem" :class="{ active: isActive('/maintenance/tasks') }">
                  <span class="submenu-dot"></span>
                  <span class="submenu-text">任务管理</span>
                </router-link>
                <router-link to="/maintenance/execute" class="nav-subitem" :class="{ active: isActive('/maintenance/execute') }">
                  <span class="submenu-dot"></span>
                  <span class="submenu-text">移动端保养</span>
                </router-link>
              </div>
            </transition>
          </div>

          <router-link to="/spareparts" class="nav-item" :class="{ active: isActive('/spareparts') }">
            <div class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
                <circle cx="12" cy="12" r="3"/>
              </svg>
            </div>
            <transition name="nav-text">
              <span v-show="!sidebarCollapsed" class="nav-text">备件管理</span>
            </transition>
            <div v-if="!sidebarCollapsed" class="nav-indicator"></div>
          </router-link>
        </div>

        <div class="nav-section">
          <div v-if="!sidebarCollapsed" class="nav-section-title">数据分析</div>
          <router-link to="/analytics" class="nav-item" :class="{ active: isActive('/analytics') }">
            <div class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <line x1="18" y1="20" x2="18" y2="10"/>
                <line x1="12" y1="20" x2="12" y2="4"/>
                <line x1="6" y1="20" x2="6" y2="14"/>
              </svg>
            </div>
            <transition name="nav-text">
              <span v-show="!sidebarCollapsed" class="nav-text">统计分析</span>
            </transition>
            <div v-if="!sidebarCollapsed" class="nav-indicator"></div>
          </router-link>
          <router-link to="/knowledge" class="nav-item" :class="{ active: isActive('/knowledge') }">
            <div class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/>
                <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"/>
              </svg>
            </div>
            <transition name="nav-text">
              <span v-show="!sidebarCollapsed" class="nav-text">知识库</span>
            </transition>
            <div v-if="!sidebarCollapsed" class="nav-indicator"></div>
          </router-link>
        </div>
      </nav>

      <!-- 侧边栏底部 -->
      <div class="sidebar-footer">
        <!-- Footer content if needed -->
      </div>
    </aside>

    <!-- 主内容区 -->
    <div class="main-content">
      <!-- 顶部栏 -->
      <header class="header">
        <div class="header-left">
          <button class="mobile-menu-btn" @click="toggleSidebar">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <line x1="3" y1="12" x2="21" y2="12"/>
              <line x1="3" y1="6" x2="21" y2="6"/>
              <line x1="3" y1="18" x2="21" y2="18"/>
            </svg>
          </button>
          <div class="breadcrumb">
            <span class="breadcrumb-current">{{ currentTitle }}</span>
          </div>
        </div>

        <div class="header-right">
          <!-- 主题切换 -->
          <button class="theme-toggle" @click="themeStore.toggleTheme()" :title="themeStore.theme === 'dark' ? '切换到浅色模式' : '切换到深色模式'">
            <svg v-if="themeStore.theme === 'dark'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <circle cx="12" cy="12" r="5"/>
              <line x1="12" y1="1" x2="12" y2="3"/>
              <line x1="12" y1="21" x2="12" y2="23"/>
              <path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>
            </svg>
            <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
            </svg>
          </button>

          <!-- 通知 -->
          <button class="header-icon-btn">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/>
              <path d="M13.73 21a2 2 0 0 1-3.46 0"/>
            </svg>
            <span class="notification-badge">3</span>
          </button>

          <!-- 用户信息 -->
          <div class="user-section">
            <div class="user-avatar">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                <circle cx="12" cy="7" r="4"/>
              </svg>
            </div>
            <div class="user-info">
              <span class="user-name">{{ authStore.userName }}</span>
              <span class="user-role">{{ roleText }}</span>
            </div>
          </div>

          <!-- 退出按钮 -->
          <button class="logout-btn" @click="handleLogout" title="退出登录">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
              <polyline points="16 17 21 12 16 7"/>
              <line x1="21" y1="12" x2="9" y2="12"/>
            </svg>
          </button>
        </div>
      </header>

      <!-- 内容区域 -->
      <main class="content">
        <router-view v-slot="{ Component }">
          <transition name="page" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>

    <!-- 移动端遮罩 -->
    <div
      v-if="!sidebarCollapsed"
      class="sidebar-overlay"
      @click="toggleSidebar"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { ElMessageBox } from 'element-plus'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const themeStore = useThemeStore()

const sidebarCollapsed = ref(false)
const expandedGroups = ref({
  organization: true,
  inspection: true,
  repair: false,
  maintenance: false,
})

const canManageOrg = computed(() => authStore.hasRole('admin', 'engineer'))

const roleTextMap: Record<string, string> = {
  admin: '系统管理员',
  supervisor: '设备主管',
  engineer: '设备工程师',
  maintenance: '维修工',
  operator: '操作工',
}

const roleText = computed(() => roleTextMap[authStore.userRole] || '未知')

const currentTitle = computed(() => route.meta?.title as string || '首页')

function isActive(path: string): boolean {
  return route.path === path || route.path.startsWith(path + '/')
}

function isGroupActive(path: string): boolean {
  return route.path.startsWith(path + '/')
}

function toggleGroup(group: keyof typeof expandedGroups) {
  if (sidebarCollapsed.value) return
  expandedGroups.value[group] = !expandedGroups.value[group]
}

function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

async function handleLogout() {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    authStore.logout()
    router.push('/login')
  } catch {
    // User cancelled
  }
}

// 响应式处理
function handleResize() {
  if (window.innerWidth < 1024) {
    sidebarCollapsed.value = true
  }
}

onMounted(() => {
  handleResize()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style>
/* Global overrides - Fix scrollbar */
html, body {
  margin: 0 !important;
  padding: 0 !important;
  border: none !important;
  outline: none !important;
  box-shadow: none !important;
  width: 100% !important;
  overflow-x: hidden !important;
  /* Important: Hide scrollbars */
  scrollbar-width: none !important;
  -ms-overflow-style: none !important;
}

#app {
  margin: 0 !important;
  padding: 0 !important;
  border: none !important;
  outline: none !important;
  box-shadow: none !important;
  width: 100% !important;
}

/* Hide all scrollbars */
*::-webkit-scrollbar {
  display: none !important;
  width: 0 !important;
  height: 0 !important;
}

* {
  scrollbar-width: none !important;
  -ms-overflow-style: none !important;
  box-sizing: border-box;
}

*:focus,
*:focus-visible,
*:focus-within {
  outline: none !important;
  box-shadow: none !important;
}
</style>

<style scoped>
/* 主布局 */
.main-layout {
  display: flex;
  min-height: 100vh;
  height: 100vh;
  background: var(--color-bg-primary);
  border: none !important;
  margin: 0;
  padding: 0;
  width: 100%;
  overflow: hidden;
}

/* 侧边栏 */
.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  width: 260px;
  height: 100vh;
  background: var(--color-bg-card);
  border-right: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  z-index: 100;
  transition: transform var(--transition-base), width var(--transition-base);
  box-sizing: border-box;
  overflow: hidden;
}

.sidebar.collapsed {
  width: 72px;
}

.sidebar-header {
  padding: 20px 16px;
  border-bottom: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: relative;
}

.logo-container {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  flex: 1;
}

/* 折叠按钮 - 固定大小 */
.sidebar-collapse-btn {
  position: absolute;
  right: 16px;
  top: 50%;
  transform: translateY(-50%);
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-tertiary);
  cursor: pointer;
  transition: all var(--transition-fast);
  flex-shrink: 0;
  /* Ensure SVG size is strict */
  padding: 0;
}

.sidebar-collapse-btn:hover {
  background: var(--color-bg-secondary);
  color: var(--color-text-secondary);
  border-color: var(--color-primary-dim);
}

.sidebar-collapse-btn svg {
  width: 18px;
  height: 18px;
  transition: transform var(--transition-fast);
}

.sidebar.collapsed .sidebar-collapse-btn svg {
  transform: rotate(180deg);
}

.logo-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--color-primary-dim), rgba(0, 255, 163, 0.1));
  border: 1.5px solid var(--color-primary-dim);
  border-radius: 10px;
  color: var(--color-primary);
  flex-shrink: 0;
}

.logo-icon svg {
  width: 24px;
  height: 24px;
}

.logo-text {
  font-size: 22px;
  font-weight: 800;
  letter-spacing: 2px;
  background: linear-gradient(135deg, var(--color-primary), var(--color-success));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* 侧边栏导航 */
.sidebar-nav {
  flex: 1;
  overflow-y: scroll; /* 始终显示滚动条以便我们隐藏它 */
  overflow-x: hidden;
  padding: 16px 8px;
  /* 强制隐藏滚动条 */
  scrollbar-width: none !important;
  -ms-overflow-style: none !important;
}

.sidebar-nav::-webkit-scrollbar {
  display: none !important;
  width: 0px !important;
  height: 0px !important;
  background: transparent !important;
}

.sidebar-nav {
  -ms-overflow-style: none !important;
  scrollbar-width: none !important;
}

.nav-section {
  margin-bottom: 24px;
}

.nav-section:last-child {
  margin-bottom: 0;
}

.nav-section-title {
  font-size: 11px;
  font-weight: 600;
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 1px;
  padding: 0 12px 8px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  margin: 2px 0;
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  text-decoration: none;
  transition: all var(--transition-fast);
  position: relative;
  cursor: pointer;
  background: transparent;
  border: none;
  width: 100%;
  text-align: left;
}

.nav-item:hover {
  background: var(--color-bg-secondary);
  color: var(--color-text-primary);
}

.nav-item.active {
  background: var(--color-primary-dim);
  color: var(--color-primary);
}

.nav-item.active .nav-icon {
  color: var(--color-primary);
}

.nav-indicator {
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 20px;
  background: var(--color-primary);
  border-radius: 0 2px 2px 0;
  opacity: 0;
  transition: opacity var(--transition-fast);
}

.nav-item.active .nav-indicator {
  opacity: 1;
}

.nav-icon {
  width: 22px;
  height: 22px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.nav-icon svg {
  width: 20px;
  height: 20px;
}

.nav-text {
  flex: 1;
  font-size: 14px;
  font-weight: 500;
}

.nav-arrow {
  width: 16px;
  height: 16px;
  transition: transform var(--transition-fast);
  color: var(--color-text-tertiary);
}

.nav-arrow.expanded {
  transform: rotate(180deg);
}

/* 子菜单 */
.nav-group {
  margin: 2px 0;
  position: relative;
}

.nav-group-toggle {
  border-radius: var(--radius-md);
}

.nav-submenu {
  padding-left: 46px;
  overflow: hidden;
}

.sidebar.collapsed .nav-submenu.always-show {
  position: absolute;
  left: 72px;
  top: 0;
  min-width: 180px;
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 8px;
  box-shadow: var(--shadow-lg);
  z-index: 50;
  opacity: 0;
  visibility: hidden;
  pointer-events: none;
  transition: opacity 0.2s, visibility 0.2s;
}

.sidebar.collapsed .nav-group:hover .nav-submenu.always-show {
  opacity: 1;
  visibility: visible;
  pointer-events: auto;
}

.nav-subitem {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  margin: 2px 0;
  border-radius: var(--radius-sm);
  color: var(--color-text-tertiary);
  text-decoration: none;
  font-size: 13px;
  transition: all var(--transition-fast);
}

.nav-subitem:hover {
  color: var(--color-text-secondary);
  background: var(--color-bg-secondary);
}

.nav-subitem.active {
  color: var(--color-primary);
}

.submenu-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
  opacity: 0.5;
}

.nav-subitem.active .submenu-dot {
  opacity: 1;
  box-shadow: 0 0 8px currentColor;
}

/* 主内容区 */
.main-content {
  flex: 1;
  margin-left: 260px;
  display: flex;
  flex-direction: column;
  height: 100vh;
  transition: margin-left var(--transition-base);
  box-sizing: border-box;
  background: var(--color-bg-primary);
  overflow: hidden;
}

.sidebar.collapsed + .main-content {
  margin-left: 72px;
}

/* 顶部栏 */
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  height: 64px;
  background: var(--color-bg-card);
  border-bottom: 1px solid var(--color-border);
  position: sticky;
  top: 0;
  z-index: 50;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.mobile-menu-btn {
  display: none;
  padding: 8px;
  background: transparent;
  border: none;
  color: var(--color-text-secondary);
  cursor: pointer;
  border-radius: var(--radius-sm);
}

.mobile-menu-btn:hover {
  background: var(--color-bg-secondary);
}

.breadcrumb {
  font-size: 14px;
}

.breadcrumb-current {
  color: var(--color-text-primary);
  font-weight: 500;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-icon-btn {
  position: relative;
  padding: 8px;
  background: transparent;
  border: none;
  color: var(--color-text-secondary);
  cursor: pointer;
  border-radius: var(--radius-sm);
  transition: all var(--transition-fast);
}

.header-icon-btn:hover {
  background: var(--color-bg-secondary);
  color: var(--color-text-primary);
}

.header-icon-btn svg {
  width: 20px;
  height: 20px;
}

/* 主题切换按钮 */
.theme-toggle {
  padding: 8px;
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all var(--transition-fast);
  display: flex;
  align-items: center;
  justify-content: center;
}

.theme-toggle:hover {
  background: var(--color-bg-secondary);
  color: var(--color-primary);
  border-color: var(--color-primary-dim);
}

.theme-toggle svg {
  width: 20px;
  height: 20px;
  transition: transform var(--transition-fast);
}

.theme-toggle:active svg {
  transform: scale(0.9);
}

.notification-badge {
  position: absolute;
  top: 4px;
  right: 4px;
  min-width: 16px;
  height: 16px;
  padding: 0 4px;
  background: var(--color-danger);
  color: white;
  font-size: 10px;
  font-weight: 600;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 用户信息 */
.user-section {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 6px 12px 6px 6px;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.user-section:hover {
  background: var(--color-bg-tertiary);
  border-color: var(--color-border-dim);
}

.user-avatar {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--color-primary), var(--color-success));
  border-radius: 8px;
  color: var(--color-bg-primary);
}

.user-avatar svg {
  width: 18px;
  height: 18px;
}

.user-info {
  display: flex;
  flex-direction: column;
}

.user-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-primary);
}

.user-role {
  font-size: 11px;
  color: var(--color-text-tertiary);
}

.logout-btn {
  padding: 8px;
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-tertiary);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.logout-btn:hover {
  background: var(--color-danger-dim);
  border-color: var(--color-danger-dim);
  color: var(--color-danger);
}

.logout-btn svg {
  width: 18px;
  height: 18px;
}

/* 内容区 */
.content {
  flex: 1;
  padding: 24px;
  overflow-y: scroll;
  overflow-x: hidden;
  /* 强制隐藏滚动条 */
  scrollbar-width: none !important;
  -ms-overflow-style: none !important;
}

.content::-webkit-scrollbar {
  display: none !important;
  width: 0px !important;
  height: 0px !important;
  background: transparent !important;
}

.content {
  -ms-overflow-style: none !important;
  scrollbar-width: none !important;
}

/* 遮罩层 */
.sidebar-overlay {
  display: none;
}

/* 过渡动画 */
.logo-text-enter-active,
.logo-text-leave-active,
.nav-text-enter-active,
.nav-text-leave-active {
  transition: all 0.2s ease;
}

.logo-text-enter-from,
.logo-text-leave-to,
.nav-text-enter-from,
.nav-text-leave-to {
  opacity: 0;
  transform: translateX(-10px);
}

.arrow-enter-active,
.arrow-leave-active {
  transition: all 0.2s ease;
}

.arrow-enter-from,
.arrow-leave-to {
  opacity: 0;
  transform: translateX(10px);
}

.submenu-enter-active,
.submenu-leave-active {
  transition: all 0.3s ease;
  overflow: hidden;
}

.submenu-enter-from,
.submenu-leave-to {
  opacity: 0;
  max-height: 0;
}

.submenu-enter-to,
.submenu-leave-from {
  opacity: 1;
  max-height: 200px;
}

/* 页面切换动画 */
.page-enter-active,
.page-leave-active {
  transition: all 0.3s ease;
}

.page-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.page-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* 响应式 */
@media (max-width: 1024px) {
  .sidebar {
    transform: translateX(-100%);
  }

  .sidebar:not(.collapsed) {
    transform: translateX(0);
  }

  .main-content {
    margin-left: 0;
  }

  .sidebar.collapsed + .main-content {
    margin-left: 0;
  }

  .mobile-menu-btn {
    display: block;
  }

  .sidebar-overlay {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    z-index: 99;
    opacity: 0;
    pointer-events: none;
    transition: opacity var(--transition-base);
  }

  .sidebar:not(.collapsed) + .main-content .sidebar-overlay {
    opacity: 1;
    pointer-events: auto;
  }

  .user-info {
    display: none;
  }

  .content {
    padding: 16px;
  }
}

@media (max-width: 640px) {
  .header {
    padding: 0 16px;
    height: 56px;
  }

  .breadcrumb {
    display: none;
  }

  .notification-badge {
    display: none;
  }
}
</style>
