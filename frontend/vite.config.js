


import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import process from 'process'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 3002,
    proxy: {
      '/api': {
        target: 'http://localhost:80',
        changeOrigin: true,
        secure: false,
        ws: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      }
    }
  },
  define: {
    'process.env': process.env,  // This defines process.env for use in the browser
  },
  resolve: {
    alias: {
      process: 'process/browser',  // This tells Vite to use the browser-friendly version of process
    }
  }
})
