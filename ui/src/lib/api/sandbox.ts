import { requestEnvelope } from './http'
import type { CreateSandboxRequest, Sandbox } from './types'

export async function listSandboxes(): Promise<Sandbox[]> {
  const data = await requestEnvelope<Sandbox[] | undefined>('/sandbox', {
    method: 'GET',
  })

  return Array.isArray(data) ? data : []
}

export async function createSandbox(payload: CreateSandboxRequest = {}): Promise<Sandbox> {
  return requestEnvelope<Sandbox>('/sandbox', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  })
}

export async function deleteSandbox(name: string): Promise<string> {
  return requestEnvelope<string>(`/sandbox/${encodeURIComponent(name)}`, {
    method: 'DELETE',
  })
}
