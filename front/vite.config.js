import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    proxy: {
      '/api/v1/auth': {
        target: 'http://localhost:8090',
        changeOrigin: true,
      },
      '/api/v1/task': {
        target: 'http://localhost:8091',
        changeOrigin: true,
      },
    }
  },
  plugins: [vue()]
})
