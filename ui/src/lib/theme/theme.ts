const THEME_STORAGE_KEY = 'ui-theme'
const DEFAULT_THEME = 'light'

const THEMES = ['light', 'dark', 'cupcake', 'synthwave'] as const

export type ThemeName = (typeof THEMES)[number]

export function getAvailableThemes(): readonly ThemeName[] {
  return THEMES
}

export function getTheme(): ThemeName {
  const raw = localStorage.getItem(THEME_STORAGE_KEY)?.trim()
  if (raw && THEMES.includes(raw as ThemeName)) {
    return raw as ThemeName
  }
  return DEFAULT_THEME
}

export function applyTheme(theme: ThemeName): void {
  document.documentElement.setAttribute('data-theme', theme)
  localStorage.setItem(THEME_STORAGE_KEY, theme)
}
