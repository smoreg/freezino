import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { visualizer } from 'rollup-plugin-visualizer'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react(),
    // Bundle analyzer - generates stats.html after build
    visualizer({
      filename: './dist/stats.html',
      open: false,
      gzipSize: true,
      brotliSize: true,
    }),
  ],
  build: {
    // Optimize chunk size and code splitting
    rollupOptions: {
      output: {
        manualChunks: {
          // Separate vendor chunks for better caching
          'react-vendor': ['react', 'react-dom', 'react-router-dom'],
          'ui-vendor': ['framer-motion', 'react-hot-toast', 'react-confetti'],
          'i18n-vendor': ['i18next', 'react-i18next'],
          'game-components': [
            './src/components/games/Roulette.tsx',
            './src/components/games/Slots.tsx',
            './src/components/games/Blackjack.tsx',
            './src/components/games/Crash.tsx',
            './src/components/games/HiLo.tsx',
            './src/components/games/Wheel.tsx',
          ],
        },
      },
    },
    // Increase chunk size warning limit (default is 500kb)
    chunkSizeWarningLimit: 1000,
    // Enable minification (esbuild is faster than terser)
    minify: 'esbuild',
    // Source maps for debugging (disable in production if needed)
    sourcemap: false,
  },
  // Optimize dependencies
  optimizeDeps: {
    include: [
      'react',
      'react-dom',
      'react-router-dom',
      'framer-motion',
      'react-hot-toast',
      'i18next',
      'react-i18next',
    ],
  },
  server: {
    // Optimize dev server
    hmr: {
      overlay: true,
    },
  },
})
