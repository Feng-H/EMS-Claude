<template>
  <div class="mobile-layout">
    <router-view v-slot="{ Component }">
      <keep-alive>
        <component :is="Component" />
      </keep-alive>
    </router-view>

    <!-- 底部导航栏 -->
    <van-tabbar v-model="active" route active-color="#c96442" inactive-color="#999">
      <van-tabbar-item replace to="/h5" icon="home-o">工作台</van-tabbar-item>
      <van-tabbar-item replace to="/h5/tasks" icon="todo-list-o">我的任务</van-tabbar-item>
      <van-tabbar-item replace to="/h5/profile" icon="user-o">我的</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const active = ref(0)
</script>

<style scoped>
.mobile-layout {
  min-height: 100vh;
  background: var(--color-bg-primary);
  padding-bottom: constant(safe-area-inset-bottom); /* 适配 iOS 底部安全区 */
  padding-bottom: env(safe-area-inset-bottom);
}

:deep(.van-tabbar) {
  background: rgba(250, 249, 245, 0.95);
  backdrop-filter: blur(10px);
  height: 60px;
  border-top: 1px solid var(--color-border);
}

:deep(.van-tabbar-item) {
  font-family: var(--font-serif);
  font-weight: 500;
}

:deep(.van-tabbar-item__icon) {
  font-size: 22px;
}
</style>
