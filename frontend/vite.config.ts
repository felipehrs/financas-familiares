import path from 'path'
import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import { VitePWA } from 'vite-plugin-pwa'

export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
    ...(process.env.VITEST
      ? []
      : [
          VitePWA({
            registerType: 'autoUpdate',
            injectRegister: 'auto',
            workbox: {
              globPatterns: ['**/*.{js,css,html,svg,woff2}'],
              navigateFallback: '/index.html',
              navigateFallbackDenylist: [/^\/api\//],
              runtimeCaching: [],
            },
            manifest: {
              name: 'Finanças Familiares',
              short_name: 'FinFam',
              description: 'Gestão financeira da família',
              start_url: '/',
              display: 'standalone',
              theme_color: '#16a34a',
              background_color: '#ffffff',
              icons: [
                {
                  src: '/icons/icon.svg',
                  sizes: '192x192',
                  type: 'image/svg+xml',
                  purpose: 'any',
                },
                {
                  src: '/icons/icon.svg',
                  sizes: '512x512',
                  type: 'image/svg+xml',
                  purpose: 'maskable',
                },
              ],
            },
          }),
        ]),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: ['./src/test/setup.ts'],
    passWithNoTests: true,
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json-summary', 'html'],
      reportsDirectory: './coverage',
      exclude: [
        'src/components/ui/**',
        'src/test/**',
        '**/*.d.ts',
        'src/main.tsx',
        'src/vite-env.d.ts',
      ],
    },
  },
})
