import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), vueDevTools()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    proxy: {
      '/api/v1/me': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },

      '/api/v1/pastes': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },

      '/api/v1/users': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
    },
  },
})
