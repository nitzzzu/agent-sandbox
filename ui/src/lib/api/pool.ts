import { requestEnvelope } from './http'
import type { PoolSandbox, PoolTemplate } from './types'

export async function listPools(): Promise<PoolTemplate[]> {
  const data = await requestEnvelope<PoolTemplate[] | undefined>('/pool', {
    method: 'GET',
  })

  return Array.isArray(data) ? data : []
}

export async function listPoolSandboxes(name: string): Promise<PoolSandbox[]> {
  const data = await requestEnvelope<PoolSandbox[] | undefined>(`/pool/sandbox/${encodeURIComponent(name)}`, {
    method: 'GET',
  })

  return Array.isArray(data) ? data : []
}

export async function deletePool(name: string): Promise<string> {
  return requestEnvelope<string>(`/pool/${encodeURIComponent(name)}`, {
    method: 'DELETE',
  })
}
