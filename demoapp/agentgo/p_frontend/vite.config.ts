import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

const hmrHost = process.env.VITE_HMR_HOST
const hmrProtocol = process.env.VITE_HMR_PROTOCOL === 'wss' ? 'wss' : 'ws'
const hmrClientPort = Number(process.env.VITE_HMR_CLIENT_PORT)
const isWindowsMountedWorkspace = process.cwd().startsWith('/mnt/')
const usePolling = process.env.VITE_USE_POLLING
  ? process.env.VITE_USE_POLLING === 'true'
  : Boolean(process.env.WSL_DISTRO_NAME) || isWindowsMountedWorkspace

export default defineConfig({
  plugins: [vue()],
  server: {
    host: '0.0.0.0',
    port: 5173,
    strictPort: true,
    hmr: hmrHost
      ? {
          host: hmrHost,
          protocol: hmrProtocol,
          ...(Number.isInteger(hmrClientPort) && hmrClientPort > 0 ? { clientPort: hmrClientPort } : {}),
        }
      : undefined,
    // Windows-mounted files in WSL2 often do not emit stable fs events. Polling
    // keeps Vue HMR reliable there, while native Linux/macOS keeps event watching.
    watch: {
      usePolling,
      interval: usePolling ? 250 : undefined,
      awaitWriteFinish: { stabilityThreshold: 150, pollInterval: 50 },
    },
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
})
