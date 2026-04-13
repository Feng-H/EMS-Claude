<template>
  <van-nav-bar
    :title="title"
    :left-text="leftText"
    :left-arrow="showBack"
    :border="border"
    :fixed="fixed"
    :z-index="zIndex"
    @click-left="onClickLeft"
    @click-right="onClickRight"
  >
    <template #left>
      <slot name="left">
        <van-icon v-if="showBack" name="arrow-left" size="18" />
      </slot>
    </template>
    <template #right>
      <slot name="right">
        <van-icon v-if="rightIcon" :name="rightIcon" size="18" />
        <span v-if="rightText">{{ rightText }}</span>
      </slot>
    </template>
  </van-nav-bar>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

interface Props {
  title?: string
  leftText?: string
  showBack?: boolean
  rightText?: string
  rightIcon?: string
  border?: boolean
  fixed?: boolean
  zIndex?: number
}

const props = withDefaults(defineProps<Props>(), {
  title: '',
  leftText: '',
  showBack: true,
  rightText: '',
  rightIcon: '',
  border: true,
  fixed: true,
  zIndex: 100
})

const emit = defineEmits<{
  clickLeft: []
  clickRight: []
}>()

const router = useRouter()

const onClickLeft = () => {
  emit('clickLeft')
  if (props.showBack) {
    router.back()
  }
}

const onClickRight = () => {
  emit('clickRight')
}
</script>

<style scoped>
:deep(.van-nav-bar) {
  background: rgba(250, 249, 245, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--color-border);
}

:deep(.van-nav-bar__title) {
  font-size: 17px;
  font-weight: 500;
  color: var(--color-text-primary);
  font-family: var(--font-serif);
}

:deep(.van-nav-bar__text) {
  color: var(--color-terracotta);
  font-size: 15px;
  font-weight: 500;
}

:deep(.van-icon) {
  color: var(--color-text-secondary);
}

@media (prefers-color-scheme: dark) {
  :deep(.van-nav-bar) {
    background: rgba(26, 26, 19, 0.95);
    border-bottom-color: var(--color-border);
  }

  :deep(.van-nav-bar__title) {
    color: var(--color-text-primary);
  }

  :deep(.van-nav-bar__text) {
    color: var(--color-terracotta);
  }

  :deep(.van-icon) {
    color: var(--color-text-tertiary);
  }
}
</style>
