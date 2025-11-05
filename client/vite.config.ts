import path from "path"
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import compression from 'vite-plugin-compression'

export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
    compression(),        // gzip (default)
    compression({
      algorithm: 'brotliCompress',
      ext: '.br'
    }),
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  build: {
    outDir: "../cmd/dashboard/dist",
    emptyOutDir: true,
  },
  // server: {
  //   hmr: false,
  // },  
})