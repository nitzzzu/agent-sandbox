import { requestEnvelope } from './http'
import type { SandboxEventsData } from './types'

type ListSandboxEventsOptions = {
  sandbox?: string
  limit?: number
}

export async function listSandboxEvents(options?: ListSandboxEventsOptions): Promise<SandboxEventsData> {
  const params = new URLSearchParams()

  if (typeof options?.sandbox === 'string') {
    const sandbox = options.sandbox.trim()
    if (sandbox) {
      params.set('sandbox', sandbox)
    }
  }

  if (typeof options?.limit === 'number' && Number.isFinite(options.limit) && options.limit > 0) {
    params.set('limit', String(Math.trunc(options.limit)))
  }

  const query = params.toString()
  const path = `/events${query ? `?${query}` : ''}`

  return requestEnvelope<SandboxEventsData>(path, {
    method: 'GET',
  })
}
