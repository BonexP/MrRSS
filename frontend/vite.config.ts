import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  css: {
    preprocessorOptions: {
      scss: {}
    }
  },
  define: {
    __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: true
  }
})
