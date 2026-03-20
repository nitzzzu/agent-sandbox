import { requestEnvelope } from './http'
import type { SandboxLogsData } from './types'

type GetSandboxLogsOptions = {
  tailLines?: number
}

export async function getSandboxLogs(name: string, options?: GetSandboxLogsOptions): Promise<SandboxLogsData> {
  const params = new URLSearchParams()

  if (typeof options?.tailLines === 'number' && Number.isFinite(options.tailLines) && options.tailLines > 0) {
    params.set('tailLines', String(Math.trunc(options.tailLines)))
  }

  const query = params.toString()
  const path = `/logs/sandbox/${encodeURIComponent(name)}${query ? `?${query}` : ''}`

  return requestEnvelope<SandboxLogsData>(path, {
    method: 'GET',
  })
}
