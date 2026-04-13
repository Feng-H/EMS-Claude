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
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid #E8EDF3;
}

:deep(.van-nav-bar__title) {
  font-size: 17px;
  font-weight: 600;
  color: #1A202C;
}

:deep(.van-nav-bar__text) {
  color: #667eea;
  font-size: 15px;
  font-weight: 500;
}

:deep(.van-icon) {
  color: #64748B;
}

@media (prefers-color-scheme: dark) {
  :deep(.van-nav-bar) {
    background: rgba(26, 35, 50, 0.95);
    border-bottom-color: #2D3748;
  }

  :deep(.van-nav-bar__title) {
    color: #F0F3FF;
  }

  :deep(.van-nav-bar__text) {
    color: #667eea;
  }

  :deep(.van-icon) {
    color: #9CA3AF;
  }
}
</style>
