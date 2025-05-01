import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  define:{'__VUE_PROD_HYDRATION_MISMATCH_DETAILS__': JSON.stringify(false),},
  plugins: [vue()],
  css: {
    postcss:'./postcss.config.js',
  },
  build: {
    rollupOptions: {
      external:['@wails/runtime']
    }
  }
})
