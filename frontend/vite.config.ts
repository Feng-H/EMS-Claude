import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
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
