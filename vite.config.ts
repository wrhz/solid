import { defineConfig } from 'vite'
import fg from 'fast-glob'
import path from 'path'
import vue from '@vitejs/plugin-vue'
import react from '@vitejs/plugin-react'

const isWatch = process.argv.includes('--watch')

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
        chunkFileNames: 'chunks/[name]-[hash].js'
      },
      external: [/^\/lib\//]
    },
    ...(isWatch
      ? {
          watch: {
            include: ['resource/src/**', 'resource/styles/**'],
            exclude: ['node_modules/**'],
            buildDelay: 300,
          }
        }
      : {}),
    cssCodeSplit: true
  }
})

function getEntries(...globPaths: string[]): Record<string, string> {
  const entries: Record<string, string> = {}

  for (const globPath of globPaths) {
    const files = fg.sync(globPath, { cwd: process.cwd() })
    for (const file of files) {
      const relative = file.replace(/\.\w+$/, '')
      entries[relative] = path.resolve(file)
    }
  }

  return entries
}