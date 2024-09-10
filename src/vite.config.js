import { defineConfig } from "vite";
import path from "path"

export default defineConfig({
  build: {
    outDir: path.join(__dirname, '..', 'templates', "static"),
    manifest: true,
    modulePreload: {
      polyfill: false,
    },
    rollupOptions: {
      input: '/src/main.ts'
    }
  }
})