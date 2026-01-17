<template>
  <div class="mobile-action-bar" :class="{ fixed: fixed }">
    <slot>
      <van-button
        v-for="(action, index) in actions"
        :key="index"
        :type="action.type || 'primary'"
        :size="action.size || 'large'"
        :text="action.text"
        :disabled="action.disabled"
        :loading="action.loading"
        :block="actions.length === 1"
        :class="{ 'flex-1': actions.length > 1 }"
        @click="handleAction(action, index)"
      />
    </slot>
  </div>
</template>

<script setup lang="ts">
import { Button as VanButton } from 'vant'

export interface ActionButton {
  text: string
  type?: 'primary' | 'success' | 'warning' | 'danger' | 'default'
  size?: 'large' | 'normal' | 'small' | 'mini'
  disabled?: boolean
  loading?: boolean
  onClick?: () => void | Promise<void>
}

interface Props {
  actions?: ActionButton[]
  fixed?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  actions: () => [],
  fixed: true
})

const emit = defineEmits<{
  action: [action: ActionButton, index: number]
}>()

const handleAction = async (action: ActionButton, index: number) => {
  emit('action', action, index)
  if (action.onClick) {
    await action.onClick()
  }
}
</script>

<style scoped>
.mobile-action-bar {
  padding: 12px 16px;
  background-color: #fff;
  border-top: 1px solid #ebedf0;
  display: flex;
  gap: 12px;
}

.mobile-action-bar.fixed {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 99;
  padding-bottom: calc(12px + env(safe-area-inset-bottom));
}

.flex-1 {
  flex: 1;
}

:deep(.van-button) {
  height: 44px;
  border-radius: 22px;
  font-size: 16px;
  font-weight: 500;
}

:deep(.van-button--primary) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
}

:deep(.van-button--success) {
  background: linear-gradient(135deg, #84fab0 0%, #8fd3f4 100%);
  border: none;
}

:deep(.van-button--danger) {
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a6f 100%);
  border: none;
}

:deep(.van-button--warning) {
  background: linear-gradient(135deg, #f9d423 0%, #ff4e50 100%);
  border: none;
}

:deep(.van-button--disabled) {
  opacity: 0.5;
}
</style>
