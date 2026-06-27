import { defineConfig } from 'vite'
import fg from 'fast-glob'
import path from 'path'
import vue from '@vitejs/plugin-vue'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [vue(), react()],
  define: {
    'process.env.NODE_ENV': JSON.stringify(
      process.env.NODE_ENV || 'production'
    )
  },
  resolve: {
    "alias": {
      "@": "/resource/lib"
    }
  },
  build: {
    lib: {
      entry: getEntries('resource/**/*.{js,ts,jsx,tsx}'),
      formats: ['es']
    },
    rollupOptions: {
      output: {
        entryFileNames: '[name].js',
        chunkFileNames: 'chunks/[name]-[hash].js',
      },
      external: [/^\/lib\//]
    }
  }
})

function getEntries(...globPaths: string[]): Record<string, string> {
  const entries: Record<string, string> = {}

  for (const globPath of globPaths) {
    const files = fg.sync(globPath, { cwd: process.cwd() })
    for (const file of files) {
      const relative = file.replace('src/components/', '').replace(/\.\w+$/, '')
      entries[relative] = path.resolve(file)
    }
  }

  return entries
}