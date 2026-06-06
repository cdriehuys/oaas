import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { heyApiPlugin } from '@hey-api/vite-plugin';
import tailwindCSS from '@tailwindcss/vite';

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindCSS(), heyApiPlugin({
    config: {
      input: '../openapi.yaml',
      output: './src/client',
      plugins: [
        {
          name: '@hey-api/client-fetch',
          runtimeConfigPath: './src/configure-client.ts',
        }
      ]
    }
  })],
  server: {
    proxy: {
      '/api': 'http://localhost:8000'
    }
  }
})
