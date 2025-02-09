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
    outDir: 'plugins/serve/dist',
    assetsDir: 'assets',
    emptyOutDir: false,
    rollupOptions: {
      input: {
        main: 'index.html'
      },
      output: {
        // 清理除了 placeholder.html 以外的文件
        assetFileNames: (assetInfo) => {
          if (assetInfo.name === 'placeholder.html') {
            return assetInfo.name;
          }
          return 'assets/[name]-[hash][extname]';
        }
      }
    }
  },
})
