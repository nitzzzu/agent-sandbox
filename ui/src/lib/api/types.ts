export type ApiEnvelope<T> = {
  code: string
  data?: T
  error?: string
}

export type Sandbox = {
  name?: string
  template?: string
  image?: string
  status?: string
  id?: string
  port?: string
  app?: string
  timeout?: number
  created_at?: string
  metadata?: Record<string, string>
}

export type CreateSandboxRequest = {
  name?: string
  template?: string
  image?: string
  timeout?: number
}

export type TemplatePool = {
  size?: number
  readySize?: number
  probePort?: number
  warmupCmd?: string
  startupCmd?: string
}

export type Template = {
  name?: string
  pattern?: string
  image?: string
  port?: number
  type?: string
  noStartupProbe?: boolean
  pool?: TemplatePool
  description?: string
}

export type PoolTemplate = Template

export type PoolSandbox = Sandbox

export type SandboxLogsData = {
  sandbox: string
  pod: string
  container: string
  logs: string
  fetchedAt: string
}

export type SandboxTerminalRequest = {
  command: string
}

export type SandboxTerminalResult = {
  sandbox: string
  pod: string
  container: string
  command: string
  output: string
  errorOutput: string
  executedAt: string
}

export type SandboxTerminalWSClientMessage =
  | {
      type: 'init'
      cols: number
      rows: number
    }
  | {
      type: 'input'
      data: string
    }
  | {
      type: 'resize'
      cols: number
      rows: number
    }
  | {
      type: 'close'
    }

export type SandboxTerminalWSServerMessage = {
  type: 'ready' | 'output' | 'error' | 'exit' | 'closed'
  data?: string
  code?: number
}

export type TerminalConnectionState = 'disconnected' | 'connecting' | 'connected'

export type TerminalSessionEvent =
  | { type: 'open' }
  | { type: 'message'; message: unknown }
  | { type: 'error'; message: string }
  | { type: 'close' }

export type SandboxFileEntry = {
  name: string
  path: string
  isDir: boolean
  size: number
}

export type SandboxFilesListData = {
  sandbox: string
  pod: string
  container: string
  path: string
  entries: SandboxFileEntry[]
  fetchedAt: string
}

export type SandboxFileUploadData = {
  sandbox: string
  pod: string
  container: string
  path: string
  fileName: string
  uploadedAt: string
}

export type SandboxFileDeleteData = {
  sandbox: string
  pod: string
  container: string
  path: string
  deletedAt: string
}

export type SandboxFileDownloadResult = {
  blob: Blob
  fileName: string
}
