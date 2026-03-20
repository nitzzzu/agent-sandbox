import { getAuthToken } from '../auth/token'
import { API_BASE_URL } from '../config/api'
import type { ApiEnvelope } from './types'

function toErrorMessage(error: unknown): string {
  if (error instanceof Error && error.message) {
    return error.message
  }
  return 'Unexpected error'
}

export async function requestEnvelope<T>(path: string, init?: RequestInit): Promise<T> {
  let response: Response
  const token = getAuthToken()
  const headers = new Headers(init?.headers)
  if (token) {
    headers.set('X-Api-Key', token)
  }

  try {
    response = await fetch(`${API_BASE_URL}${path}`, {
      ...init,
      headers,
    })
  } catch (error) {
    throw new Error(`Network error: ${toErrorMessage(error)}`)
  }

  let envelope: ApiEnvelope<T>
  try {
    envelope = (await response.json()) as ApiEnvelope<T>
  } catch {
    if (!response.ok) {
      throw new Error(`Request failed with status ${response.status}`)
    }
    throw new Error('Invalid JSON response')
  }

  if (!response.ok || envelope.code !== '0') {
    throw new Error(envelope.error || `Request failed (code: ${envelope.code || String(response.status)})`)
  }

  return envelope.data as T
}
