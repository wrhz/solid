import { defineConfig } from 'vite'
import fg from 'fast-glob'
import path from 'path'
import { execSync } from 'child_process'

export default defineConfig({
  resolve: {
    "alias": {
      "@": "/lib"
    }
  },
  build: {
    lib: {
      entry: getEntries(
        'resource/lib/**/*.{js,ts}',
        'static/js/**/*{.js,ts}',
        'static/ts/**/*{.js,ts}'
      ),
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