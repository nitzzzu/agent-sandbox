import { ChangeEvent, useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { useLocation, useNavigate, useParams } from 'react-router-dom'

import { deletePool, listPoolSandboxes } from '../lib/api/pool'
import type { PoolSandbox, PoolTemplate } from '../lib/api/types'

type PoolSandboxSortKey = 'created-desc' | 'created-asc' | 'name-asc' | 'name-desc'

const refreshIntervalOptions = [2000, 5000, 10000, 30000]

function formatCreatedAt(value?: string): string {
  if (!value) {
    return '-'
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }

  return date.toLocaleString()
}

function getCreatedTimestamp(value?: string): number {
  if (!value) {
    return 0
  }

  const parsed = Date.parse(value)
  if (Number.isNaN(parsed)) {
    return 0
  }

  return parsed
}

function decodePoolName(value?: string): string {
  if (!value) {
    return ''
  }

  try {
    return decodeURIComponent(value)
  } catch {
    return value
  }
}

type PoolDetailLocationState = {
  selectedPool?: PoolTemplate
}

export default function PoolDetailPage() {
  const navigate = useNavigate()
  const location = useLocation()
  const { poolName: poolNameParam } = useParams()
  const poolName = decodePoolName(poolNameParam)
  const selectedPool = (location.state as PoolDetailLocationState | null)?.selectedPool ?? null

  const [poolSandboxes, setPoolSandboxes] = useState<PoolSandbox[]>([])
  const [poolTemplate, setPoolTemplate] = useState<PoolTemplate | null>(null)
  const [isPoolDetailLoading, setIsPoolDetailLoading] = useState(false)
  const [poolDetailError, setPoolDetailError] = useState('')
  const [poolMetaError, setPoolMetaError] = useState('')
  const [poolActionError, setPoolActionError] = useState('')
  const [poolSuccessMessage, setPoolSuccessMessage] = useState('')
  const [deletingPoolName, setDeletingPoolName] = useState<string | null>(null)
  const [confirmDeletePoolName, setConfirmDeletePoolName] = useState<string | null>(null)
  const [isAutoRefresh, setIsAutoRefresh] = useState(true)
  const [refreshIntervalMs, setRefreshIntervalMs] = useState(5000)
  const [sortKey, setSortKey] = useState<PoolSandboxSortKey>('created-desc')

  const poolDetailRequestInFlightRef = useRef(false)
  const deletingPoolNameRef = useRef<string | null>(null)

  deletingPoolNameRef.current = deletingPoolName

  const loadPoolSandboxes = useCallback(async (name: string, options?: { keepMessages?: boolean }) => {
    if (poolDetailRequestInFlightRef.current) {
      return
    }

    poolDetailRequestInFlightRef.current = true
    setIsPoolDetailLoading(true)
    setPoolDetailError('')

    if (!options?.keepMessages) {
      setPoolActionError('')
      setPoolSuccessMessage('')
    }

    try {
      const sandboxes = await listPoolSandboxes(name)
      setPoolSandboxes(sandboxes)
    } catch (error) {
      const message = error instanceof Error ? error.message : `Failed to load sandboxes for pool "${name}"`
      setPoolDetailError(message)
      setPoolSandboxes([])
    } finally {
      poolDetailRequestInFlightRef.current = false
      setIsPoolDetailLoading(false)
    }
  }, [])

  useEffect(() => {
    if (!poolName) {
      setPoolActionError('Pool name is missing and cannot be viewed.')
      setPoolSuccessMessage('')
      return
    }

    if (!selectedPool || (selectedPool.name?.trim() || '') !== poolName) {
      setPoolMetaError(`Failed to load metadata for pool "${poolName}".`)
      setPoolTemplate(null)
    } else {
      setPoolMetaError('')
      setPoolTemplate(selectedPool)
    }

    void loadPoolSandboxes(poolName)
  }, [poolName, selectedPool, loadPoolSandboxes])

  useEffect(() => {
    if (!isAutoRefresh || !poolName) {
      return
    }

    const timer = window.setInterval(() => {
      if (deletingPoolNameRef.current) {
        return
      }
      void loadPoolSandboxes(poolName, { keepMessages: true })
    }, refreshIntervalMs)

    return () => {
      window.clearInterval(timer)
    }
  }, [isAutoRefresh, refreshIntervalMs, poolName, loadPoolSandboxes])

  const sortedPoolSandboxes = useMemo(() => {
    const next = [...poolSandboxes]

    next.sort((a, b) => {
      if (sortKey === 'created-desc') {
        return getCreatedTimestamp(b.created_at) - getCreatedTimestamp(a.created_at)
      }

      if (sortKey === 'created-asc') {
        return getCreatedTimestamp(a.created_at) - getCreatedTimestamp(b.created_at)
      }

      const nameA = (a.name ?? '').toLocaleLowerCase()
      const nameB = (b.name ?? '').toLocaleLowerCase()
      if (sortKey === 'name-asc') {
        return nameA.localeCompare(nameB)
      }

      return nameB.localeCompare(nameA)
    })

    return next
  }, [poolSandboxes, sortKey])

  const handleRefreshPoolDetail = async () => {
    if (!poolName) {
      return
    }

    await loadPoolSandboxes(poolName, { keepMessages: true })
  }

  const handleDeletePool = async (name?: string) => {
    if (!name) {
      setPoolActionError('Pool name is missing and cannot be deleted.')
      setPoolSuccessMessage('')
      setConfirmDeletePoolName(null)
      return
    }

    setDeletingPoolName(name)
    setPoolActionError('')
    setPoolSuccessMessage('')

    try {
      await deletePool(name)
      navigate('..', {
        replace: true,
        state: {
          successMessage: `Pool "${name}" deleted successfully.`,
        },
      })
    } catch (error) {
      const message = error instanceof Error ? error.message : `Failed to delete pool "${name}"`
      setPoolActionError(message)
    } finally {
      setDeletingPoolName(null)
      setConfirmDeletePoolName(null)
    }
  }

  const handleRefreshIntervalChange = (event: ChangeEvent<HTMLSelectElement>) => {
    const parsed = Number.parseInt(event.target.value, 10)
    if (!Number.isNaN(parsed) && parsed > 0) {
      setRefreshIntervalMs(parsed)
    }
  }

  const handleSortChange = (event: ChangeEvent<HTMLSelectElement>) => {
    setSortKey(event.target.value as PoolSandboxSortKey)
  }

  return (
    <>
      <header className="card border border-base-300 bg-base-100 shadow-sm">
        <div className="card-body gap-3">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div>
              <h2 className="text-2xl font-semibold">Pool Management - Pool Detail</h2>
              <p className="text-sm text-base-content/70">Manage pools and inspect sandboxes grouped by template pool.</p>
            </div>
          </div>
        </div>
      </header>

      {(poolActionError || poolSuccessMessage) && (
        <section className="space-y-2">
          {poolActionError && (
            <div className="alert alert-error">
              <span>{poolActionError}</span>
            </div>
          )}
          {poolSuccessMessage && (
            <div className="alert alert-success">
              <span>{poolSuccessMessage}</span>
            </div>
          )}
        </section>
      )}

      <section>
        <div className="card border border-base-300 bg-base-100 shadow-sm">
          <div className="card-body gap-4">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <div>
                <div className="text-sm text-base-content/70">Pool Name</div>
                <div className="text-lg font-semibold">{poolName || '-'}</div>
              </div>
              <div className="flex flex-wrap items-center gap-2">
                <label className="label cursor-pointer gap-2 py-0">
                  <span className="label-text text-sm">Auto Refresh</span>
                  <input
                    className="toggle toggle-sm"
                    type="checkbox"
                    checked={isAutoRefresh}
                    onChange={() => {
                      setIsAutoRefresh((prev) => !prev)
                    }}
                    disabled={!poolName || Boolean(deletingPoolName)}
                  />
                </label>

                <label className="flex items-center gap-2">
                  <span className="text-sm">Interval</span>
                  <select className="select select-sm select-bordered" value={String(refreshIntervalMs)} onChange={handleRefreshIntervalChange} disabled={!isAutoRefresh || Boolean(deletingPoolName)}>
                    {refreshIntervalOptions.map((interval) => (
                      <option key={interval} value={interval}>
                        {interval / 1000}s
                      </option>
                    ))}
                  </select>
                </label>

                <label className="flex items-center gap-2">
                  <span className="text-sm">Sort</span>
                  <select className="select select-sm select-bordered" value={sortKey} onChange={handleSortChange}>
                    <option value="created-desc">Created (newest first)</option>
                    <option value="created-asc">Created (oldest first)</option>
                    <option value="name-asc">Name (A-Z)</option>
                    <option value="name-desc">Name (Z-A)</option>
                  </select>
                </label>

                <button className="btn btn-sm btn-ghost" type="button" onClick={() => navigate('/pool')} disabled={Boolean(deletingPoolName)}>
                  Back to List
                </button>
                <button
                  className={`btn btn-sm btn-outline ${isPoolDetailLoading ? 'btn-disabled' : ''}`}
                  type="button"
                  onClick={() => {
                    void handleRefreshPoolDetail()
                  }}
                  disabled={!poolName || isPoolDetailLoading || Boolean(deletingPoolName)}
                >
                  {isPoolDetailLoading ? 'Refreshing...' : 'Refresh'}
                </button>
                <button
                  className={`btn btn-sm btn-error ${deletingPoolName ? 'btn-disabled' : ''}`}
                  type="button"
                  disabled={!poolName || Boolean(deletingPoolName)}
                  onClick={() => {
                    setConfirmDeletePoolName(poolName || null)
                  }}
                >
                  {deletingPoolName === poolName ? 'Deleting...' : 'Delete Pool'}
                </button>
              </div>
            </div>

            {(poolDetailError || poolMetaError) && (
              <div className="space-y-2">
                {poolDetailError && (
                  <div className="alert alert-error">
                    <span>{poolDetailError}</span>
                  </div>
                )}
                {poolMetaError && (
                  <div className="alert alert-error">
                    <span>{poolMetaError}</span>
                  </div>
                )}
              </div>
            )}

            <div className="rounded-box border border-base-300 bg-base-200/40 p-4 space-y-4">
              <div className="mb-2 text-sm font-medium">Pool Details</div>

              <div>
                <div className="text-xs text-base-content/70">Image</div>
                <div className="text-sm break-all">{poolTemplate?.image || '-'}</div>
              </div>
              <div>
                <div className="text-xs text-base-content/70">Description</div>
                <div className="text-sm">{poolTemplate?.description || '-'}</div>
              </div>
              <div>
                <div className="text-xs text-base-content/70">Port</div>
                <div className="text-sm">{typeof poolTemplate?.port === 'number' ? poolTemplate.port : '-'}</div>
              </div>

              <div className="rounded-box border border-base-300 bg-base-100/70 p-3">
                <div className="mb-2 text-sm font-medium">Pool Config</div>
                <div className="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
                  <div>
                    <div className="text-xs text-base-content/70">Pool Size</div>
                    <div className="text-sm">{typeof poolTemplate?.pool?.size === 'number' ? poolTemplate.pool.size : '-'}</div>
                  </div>
                  <div>
                    <div className="text-xs text-base-content/70">Pool Ready Size</div>
                    <div className="text-sm">{typeof poolTemplate?.pool?.readySize === 'number' ? poolTemplate.pool.readySize : '-'}</div>
                  </div>
                  <div>
                    <div className="text-xs text-base-content/70">Pool Probe Port</div>
                    <div className="text-sm">{typeof poolTemplate?.pool?.probePort === 'number' ? poolTemplate.pool.probePort : '-'}</div>
                  </div>
                  <div>
                    <div className="text-xs text-base-content/70">Pool Warmup Command</div>
                    <div className="text-sm break-all">{poolTemplate?.pool?.warmupCmd || '-'}</div>
                  </div>
                  <div>
                    <div className="text-xs text-base-content/70">Pool Startup Command</div>
                    <div className="text-sm break-all">{poolTemplate?.pool?.startupCmd || '-'}</div>
                  </div>
                </div>
              </div>
            </div>

            <div className="custom-scrollbar overflow-x-auto rounded-box border border-base-300">
              <table className="table table-zebra">
                <thead>
                  <tr>
                    <th>#</th>
                    <th>Name</th>
                    <th>Template</th>
                    <th>Status</th>
                    <th>Timeout</th>
                    <th>Created</th>
                  </tr>
                </thead>
                <tbody>
                  {sortedPoolSandboxes.length === 0 ? (
                    <tr>
                      <td className="text-center text-base-content/70" colSpan={6}>
                        {isPoolDetailLoading ? 'Loading pool sandboxes...' : 'No sandboxes found in this pool.'}
                      </td>
                    </tr>
                  ) : (
                    sortedPoolSandboxes.map((sandbox, index) => (
                      <tr key={sandbox.name || sandbox.id || `pool-sandbox-${index}`}>
                        <td>{index + 1}</td>
                        <td className="font-medium">{sandbox.name || '-'}</td>
                        <td>{sandbox.template || '-'}</td>
                        <td>
                          {sandbox.status ? (
                            <span className={`badge ${sandbox.status === 'running' ? 'badge-success' : 'badge-warning'}`}>{sandbox.status}</span>
                          ) : (
                            '-'
                          )}
                        </td>
                        <td>{typeof sandbox.timeout === 'number' ? `${sandbox.timeout}s` : '-'}</td>
                        <td>{formatCreatedAt(sandbox.created_at)}</td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </section>

      <dialog className={`modal ${confirmDeletePoolName ? 'modal-open' : ''}`} open={Boolean(confirmDeletePoolName)}>
        <div className="modal-box">
          <h3 className="text-lg font-bold">Delete Pool</h3>
          <p className="py-4">
            Are you sure you want to delete pool <span className="font-semibold">"{confirmDeletePoolName}"</span> and all sandboxes in it?
          </p>
          <div className="modal-action">
            <button
              className="btn btn-ghost"
              type="button"
              onClick={() => {
                if (deletingPoolName) {
                  return
                }
                setConfirmDeletePoolName(null)
              }}
              disabled={Boolean(deletingPoolName)}
            >
              Cancel
            </button>
            <button
              className={`btn btn-error ${deletingPoolName ? 'btn-disabled' : ''}`}
              type="button"
              disabled={Boolean(deletingPoolName)}
              onClick={() => {
                void handleDeletePool(confirmDeletePoolName ?? undefined)
              }}
            >
              {deletingPoolName ? 'Deleting...' : 'Delete Pool'}
            </button>
          </div>
        </div>
        <button
          className="modal-backdrop"
          type="button"
          aria-label="Close"
          onClick={() => {
            if (deletingPoolName) {
              return
            }
            setConfirmDeletePoolName(null)
          }}
        />
      </dialog>
    </>
  )
}
