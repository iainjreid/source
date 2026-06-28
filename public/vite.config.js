import { defineConfig } from 'vite';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  root: './public',
  plugins: [
    tailwindcss(),
  ],
  build: {
    minify: true,
    license: true,
  },
  esbuild: {
    jsx: 'automatic',
    legalComments: 'none',
  },
});
