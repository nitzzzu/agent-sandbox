import type { Config } from 'tailwindcss'
import daisyui from 'daisyui'

type DaisyUiConfig = {
  daisyui?: {
    themes?: string[]
  }
}

export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {},
  },
  plugins: [daisyui],
  daisyui: {
    themes: ['light', 'dark'],
  },
} satisfies Config & DaisyUiConfig
