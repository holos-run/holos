import { defineConfig } from 'vite'
import tsconfigPaths from 'vite-tsconfig-paths'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react(), tsconfigPaths()],
  base: '/app/',
  build: {
    outDir: "dist/app"
  },
  server: {
    // https://vitejs.dev/config/server-options.html#server-host
    // host true listens on all addresses
    host: true,
    port: 3000,
    strictPort: true,
    proxy: {
      // Handle /app/ in Vite for hot module reloading.
      // Proxy all paths *except* /app/ to the backend.
      '^/(?!app/).*$': {
        target: 'http://localhost:4000'
      }
    }
  },
})
