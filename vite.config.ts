import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), tailwindcss()],
  server: {
    host: '0.0.0.0',
    proxy: {
      '/web_api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/web_api/, '/web_api'),
      },
    },
  },
  build: {
    outDir: 'internal/web/dist',
    assetsDir: 'assets',
    emptyOutDir: true,
  },
})
