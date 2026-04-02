import { useEffect, useRef, useState } from 'react'
import { useSearchParams } from 'react-router-dom'

import { createSandboxTrafficSession } from '../lib/api/traffic'
import type { TrafficFlow } from '../lib/api/types'
import { listSandboxes } from '../lib/api/sandbox'
import type { Sandbox } from '../lib/api/types'

const MAX_FLOWS = 500

function formatBytes(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

function rowColorClass(status?: number, flowType?: string): string {
  if (flowType === 'error') return 'bg-red-900/30 text-red-300'
  if (!status) return ''
  if (status >= 500) return 'bg-red-900/20'
  if (status >= 400) return 'bg-orange-900/20'
  if (status >= 300) return 'bg-yellow-900/20'
  if (status >= 200) return 'bg-green-900/10'
  return ''
}

function statusBadgeClass(status?: number): string {
  if (!status) return 'badge badge-ghost badge-sm'
  if (status >= 500) return 'badge badge-error badge-sm'
  if (status >= 400) return 'badge badge-warning badge-sm'
  if (status >= 300) return 'badge badge-info badge-sm'
  if (status >= 200) return 'badge badge-success badge-sm'
  return 'badge badge-ghost badge-sm'
}

export default function TrafficPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const [sandboxes, setSandboxes] = useState<Sandbox[]>([])
  const [selectedName, setSelectedName] = useState('')
  const [flows, setFlows] = useState<TrafficFlow[]>([])
  const [connected, setConnected] = useState(false)
  const [filter, setFilter] = useState('')
  const [methodFilter, setMethodFilter] = useState('All')
  const [protocolFilter, setProtocolFilter] = useState('All')
  const [sandboxesLoading, setSandboxesLoading] = useState(false)

  const sessionRef = useRef<ReturnType<typeof createSandboxTrafficSession> | null>(null)
  const bottomRef = useRef<HTMLDivElement>(null)
  const tableRef = useRef<HTMLDivElement>(null)
  const autoScrollRef = useRef(true)
  const selectedNameRef = useRef('')
  selectedNameRef.current = selectedName

  useEffect(() => {
    setSandboxesLoading(true)
    listSandboxes()
      .then((list) => {
        setSandboxes(list)
        const fromQuery = searchParams.get('sandbox')?.trim() ?? ''
        const hasQueryTarget = Boolean(fromQuery) && list.some((s) => s.name === fromQuery)
        const next = hasQueryTarget ? fromQuery : list[0]?.name ?? ''
        setSelectedName(next)
        if (next) {
          const p = new URLSearchParams(searchParams)
          p.set('sandbox', next)
          setSearchParams(p, { replace: true })
        }
      })
      .catch(() => {})
      .finally(() => setSandboxesLoading(false))
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (!selectedName) return

    sessionRef.current?.close()
    setFlows([])
    setConnected(false)

    const session = createSandboxTrafficSession(selectedName, (event) => {
      if (event.type === 'open') setConnected(true)
      if (event.type === 'close') setConnected(false)
      if (event.type === 'flow') {
        setFlows((prev) => [...prev.slice(-(MAX_FLOWS - 1)), event.flow])
      }
    })

    sessionRef.current = session
    return () => {
      session.close()
      setConnected(false)
    }
  }, [selectedName])

  useEffect(() => {
    if (autoScrollRef.current && bottomRef.current) {
      bottomRef.current.scrollIntoView({ behavior: 'smooth' })
    }
  }, [flows])

  const handleTableScroll = () => {
    const el = tableRef.current
    if (!el) return
    const atBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 40
    autoScrollRef.current = atBottom
  }

  const handleSandboxChange = (name: string) => {
    setSelectedName(name)
    const p = new URLSearchParams(searchParams)
    if (name) {
      p.set('sandbox', name)
    } else {
      p.delete('sandbox')
    }
    setSearchParams(p, { replace: true })
  }

  const visibleFlows = flows.filter((f) => {
    if (filter && !f.url.toLowerCase().includes(filter.toLowerCase())) return false
    if (methodFilter !== 'All' && f.method !== methodFilter) return false
    if (protocolFilter !== 'All') {
      const isHttps = f.url.startsWith('https://')
      if (protocolFilter === 'HTTPS' && !isHttps) return false
      if (protocolFilter === 'HTTP' && isHttps) return false
    }
    return true
  })

  const methods = ['All', ...Array.from(new Set(flows.map((f) => f.method).filter(Boolean)))] as string[]

  return (
    <>
      <header className="card border border-base-300 bg-base-100 shadow-sm">
        <div className="card-body gap-3">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <div>
              <h2 className="text-2xl font-semibold">Traffic Monitor</h2>
              <p className="text-sm text-base-content/70">
                Live HTTP/HTTPS traffic via mitmproxy sidecar. Requires sandbox started with{' '}
                <code className="text-xs">metadata.mitm=true</code>.
              </p>
            </div>
            <div className="flex items-center gap-2">
              <span className={`badge ${connected ? 'badge-success' : 'badge-ghost'} gap-1`}>
                <span>{connected ? '●' : '○'}</span>
                {connected ? 'LIVE' : 'DISCONNECTED'}
              </span>
              <button
                className="btn btn-sm btn-outline"
                type="button"
                onClick={() => setFlows([])}
              >
                Clear
              </button>
            </div>
          </div>

          <div className="flex flex-wrap items-center gap-2">
            <label className="flex items-center gap-2">
              <span className="text-sm">Sandbox</span>
              <select
                className="select select-sm select-bordered"
                style={{ width: '320px' }}
                value={selectedName}
                onChange={(e) => handleSandboxChange(e.target.value)}
                disabled={sandboxesLoading || sandboxes.length === 0}
              >
                {sandboxes.length === 0 ? (
                  <option value="">No sandboxes</option>
                ) : (
                  sandboxes.map((s) => (
                    <option key={s.name ?? s.id} value={s.name ?? ''}>
                      {s.name ?? '-'}
                    </option>
                  ))
                )}
              </select>
            </label>

            <input
              className="input input-sm input-bordered"
              placeholder="Filter by URL..."
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
              style={{ width: '220px' }}
            />

            <label className="flex items-center gap-2">
              <span className="text-sm">Protocol</span>
              <select
                className="select select-sm select-bordered"
                value={protocolFilter}
                onChange={(e) => setProtocolFilter(e.target.value)}
              >
                <option>All</option>
                <option>HTTP</option>
                <option>HTTPS</option>
              </select>
            </label>

            <label className="flex items-center gap-2">
              <span className="text-sm">Method</span>
              <select
                className="select select-sm select-bordered"
                value={methodFilter}
                onChange={(e) => setMethodFilter(e.target.value)}
              >
                {methods.map((m) => (
                  <option key={m}>{m}</option>
                ))}
              </select>
            </label>
          </div>
        </div>
      </header>

      <section>
        <div className="card border border-base-300 bg-base-100 shadow-sm">
          <div className="card-body p-0">
            {!selectedName ? (
              <p className="p-6 text-sm text-base-content/60">Select a sandbox to begin monitoring traffic.</p>
            ) : (
              <div
                ref={tableRef}
                className="overflow-auto"
                style={{ maxHeight: 'calc(100vh - 320px)' }}
                onScroll={handleTableScroll}
              >
                <table className="table table-xs table-pin-rows w-full">
                  <thead>
                    <tr>
                      <th className="w-20">Method</th>
                      <th>URL</th>
                      <th className="w-16">Status</th>
                      <th className="w-16">ms</th>
                      <th className="w-20">Size</th>
                    </tr>
                  </thead>
                  <tbody>
                    {visibleFlows.length === 0 ? (
                      <tr>
                        <td colSpan={5} className="py-6 text-center text-sm text-base-content/50">
                          {connected ? 'Waiting for traffic...' : 'Not connected — sandbox may not have mitm=true metadata.'}
                        </td>
                      </tr>
                    ) : (
                      visibleFlows.map((f, i) => (
                        <tr key={i} className={rowColorClass(f.status, f.type)}>
                          <td>
                            {f.method ? (
                              <span className="badge badge-outline badge-xs font-mono">{f.method}</span>
                            ) : (
                              <span className="text-base-content/40">—</span>
                            )}
                          </td>
                          <td
                            className="max-w-xs overflow-hidden text-ellipsis whitespace-nowrap font-mono text-xs"
                            title={f.url}
                          >
                            {f.type === 'error' ? (
                              <span className="text-error">{f.message ?? f.url}</span>
                            ) : (
                              f.url
                            )}
                          </td>
                          <td>
                            {f.status != null ? (
                              <span className={statusBadgeClass(f.status)}>{f.status}</span>
                            ) : (
                              <span className="text-base-content/40">—</span>
                            )}
                          </td>
                          <td className="text-right font-mono text-xs">
                            {f.duration_ms != null ? f.duration_ms : <span className="text-base-content/40">—</span>}
                          </td>
                          <td className="text-right font-mono text-xs">
                            {f.res_size != null ? formatBytes(f.res_size) : <span className="text-base-content/40">—</span>}
                          </td>
                        </tr>
                      ))
                    )}
                  </tbody>
                </table>
                <div ref={bottomRef} />
              </div>
            )}
          </div>
        </div>
      </section>
    </>
  )
}
