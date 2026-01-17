import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  // 监听所有网络接口，让手机可以访问
  server: {
    host: '0.0.0.0',
    port: 5173,
    strictPort: true,
    // API代理配置，让手机也能访问后端
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path
      }
    }
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  build: {
    rollupOptions: {
      onwarn(warning, warn) {
        // Ignore certain warnings during build
        if (warning.code === 'CIRCULAR_DEPENDENCY') return
        if (warning.code === 'EMPTY_BUNDLE') return
        warn(warning)
      }
    }
  }
})
