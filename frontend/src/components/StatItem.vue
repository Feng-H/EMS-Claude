<template>
  <div class="stat-item">
    <el-icon class="stat-icon" :style="{ color }">
      <component :is="iconComponent" />
    </el-icon>
    <div class="stat-content">
      <div class="stat-label">{{ label }}</div>
      <div class="stat-value">{{ value }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

const props = defineProps<{
  icon: string
  label: string
  value: string | number
  color?: string
}>()

const iconComponent = computed(() => {
  return (ElementPlusIconsVue as any)[props.icon] || ElementPlusIconsVue.DataAnalysis
})
</script>

<style scoped>
.stat-item {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  font-size: 40px;
  opacity: 0.8;
}

.stat-content {
  flex: 1;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}
</style>
