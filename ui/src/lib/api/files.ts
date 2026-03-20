import { getAuthToken } from '../auth/token'
import { API_BASE_URL } from '../config/api'
import { requestEnvelope } from './http'
import type { SandboxFileDeleteData, SandboxFileDownloadResult, SandboxFileUploadData, SandboxFilesListData } from './types'

export async function listSandboxFiles(name: string, filePath = '/'): Promise<SandboxFilesListData> {
  const params = new URLSearchParams()
  params.set('path', filePath)

  return requestEnvelope<SandboxFilesListData>(`/sandbox/files/${encodeURIComponent(name)}?${params.toString()}`, {
    method: 'GET',
  })
}

export async function uploadSandboxFile(name: string, filePath: string, file: File): Promise<SandboxFileUploadData> {
  const params = new URLSearchParams()
  params.set('path', filePath)

  const formData = new FormData()
  formData.append('file', file)

  return requestEnvelope<SandboxFileUploadData>(`/sandbox/files/${encodeURIComponent(name)}/upload?${params.toString()}`, {
    method: 'POST',
    body: formData,
  })
}

export async function deleteSandboxFile(name: string, filePath: string): Promise<SandboxFileDeleteData> {
  const params = new URLSearchParams()
  params.set('path', filePath)

  return requestEnvelope<SandboxFileDeleteData>(`/sandbox/files/${encodeURIComponent(name)}?${params.toString()}`, {
    method: 'DELETE',
  })
}

function getFileNameFromDisposition(disposition: string | null): string {
  if (!disposition) {
    return 'download.bin'
  }

  const utfMatch = disposition.match(/filename\*=UTF-8''([^;]+)/i)
  if (utfMatch?.[1]) {
    return decodeURIComponent(utfMatch[1])
  }

  const normalMatch = disposition.match(/filename="?([^";]+)"?/i)
  if (normalMatch?.[1]) {
    return normalMatch[1]
  }

  return 'download.bin'
}

export async function downloadSandboxFile(name: string, filePath: string): Promise<SandboxFileDownloadResult> {
  const params = new URLSearchParams()
  params.set('path', filePath)

  const token = getAuthToken()
  const headers = new Headers()
  if (token) {
    headers.set('X-Api-Key', token)
  }

  const response = await fetch(`${API_BASE_URL}/sandbox/files/${encodeURIComponent(name)}/download?${params.toString()}`, {
    headers,
  })
  if (!response.ok) {
    let message = `Download failed with status ${response.status}`
    try {
      const payload = (await response.json()) as { error?: string }
      if (payload.error) {
        message = payload.error
      }
    } catch {
      // ignore non-json response
    }
    throw new Error(message)
  }

  return {
    blob: await response.blob(),
    fileName: getFileNameFromDisposition(response.headers.get('Content-Disposition')),
  }
}

export function triggerBrowserDownload(fileName: string, blob: Blob) {
  const url = URL.createObjectURL(blob)
  const anchor = document.createElement('a')
  anchor.href = url
  anchor.download = fileName
  document.body.appendChild(anchor)
  anchor.click()
  anchor.remove()
  URL.revokeObjectURL(url)
}
