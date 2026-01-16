import { defineStore } from 'pinia'
import { ref, watch, computed } from 'vue'

export type Theme = 'light' | 'dark'

const STORAGE_KEY = 'ems-theme'
const DARK_CLASS = 'dark'

export const useThemeStore = defineStore('theme', () => {
  // 从 localStorage 读取主题，默认为 dark
  const stored = localStorage.getItem(STORAGE_KEY) as Theme | null
  const theme = ref<Theme>(stored || 'dark')

  // 应用主题到 HTML 根元素和 body
  function applyTheme(themeValue: Theme) {
    const html = document.documentElement
    const body = document.body

    // 设置 data-theme 属性
    html.setAttribute('data-theme', themeValue)

    // Element Plus 使用 'dark' 类来识别深色模式
    if (themeValue === 'dark') {
      html.classList.add(DARK_CLASS)
      body.classList.add(DARK_CLASS)
    } else {
      html.classList.remove(DARK_CLASS)
      body.classList.remove(DARK_CLASS)
    }
  }

  // 初始化时应用主题
  applyTheme(theme.value)

  // 监听主题变化
  watch(theme, (newTheme) => {
    applyTheme(newTheme)
    localStorage.setItem(STORAGE_KEY, newTheme)
  })

  // 切换主题
  function toggleTheme() {
    theme.value = theme.value === 'dark' ? 'light' : 'dark'
  }

  // 设置主题
  function setTheme(newTheme: Theme) {
    theme.value = newTheme
  }

  return {
    theme,
    isDark: computed(() => theme.value === 'dark'),
    toggleTheme,
    setTheme,
  }
})
