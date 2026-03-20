const TOKEN_STORAGE_KEY = 'api-token'

export function getAuthToken(): string {
  try {
    return (localStorage.getItem(TOKEN_STORAGE_KEY) || '').trim()
  } catch {
    return ''
  }
}

export function setAuthToken(token: string): void {
  localStorage.setItem(TOKEN_STORAGE_KEY, token.trim())
}

export function clearAuthToken(): void {
  localStorage.removeItem(TOKEN_STORAGE_KEY)
}

export function hasAuthToken(): boolean {
  return getAuthToken() !== ''
}
