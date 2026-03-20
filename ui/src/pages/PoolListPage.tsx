import { useEffect, useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'

import { deletePool, listPools } from '../lib/api/pool'
import type { PoolTemplate } from '../lib/api/types'

type PoolListLocationState = {
  successMessage?: string
}

export default function PoolListPage() {
  const location = useLocation()
  const navigate = useNavigate()

  const [pools, setPools] = useState<PoolTemplate[]>([])
  const [isPoolsLoading, setIsPoolsLoading] = useState(false)
  const [poolsLoadError, setPoolsLoadError] = useState('')
  const [poolActionError, setPoolActionError] = useState('')
  const [poolSuccessMessage, setPoolSuccessMessage] = useState('')
  const [deletingPoolName, setDeletingPoolName] = useState<string | null>(null)
  const [confirmDeletePoolName, setConfirmDeletePoolName] = useState<string | null>(null)

  const loadPools = async (options?: { keepMessages?: boolean }) => {
    setIsPoolsLoading(true)
    setPoolsLoadError('')

    if (!options?.keepMessages) {
      setPoolActionError('')
      setPoolSuccessMessage('')
    }

    try {
      const list = await listPools()
      setPools(list)
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to load pools'
      setPoolsLoadError(message)
      setPools([])
    } finally {
      setIsPoolsLoading(false)
    }
  }

  useEffect(() => {
    void loadPools()
  }, [])

  useEffect(() => {
    const state = location.state as PoolListLocationState | null
    if (!state?.successMessage) {
      return
    }

    setPoolActionError('')
    setPoolSuccessMessage(state.successMessage)
    navigate(location.pathname, { replace: true, state: null })
  }, [location.pathname, location.state, navigate])

  const handleReloadPools = async () => {
    await loadPools({ keepMessages: true })
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
      setPoolSuccessMessage(`Pool "${name}" deleted successfully.`)
      await loadPools({ keepMessages: true })
    } catch (error) {
      const message = error instanceof Error ? error.message : `Failed to delete pool "${name}"`
      setPoolActionError(message)
    } finally {
      setDeletingPoolName(null)
      setConfirmDeletePoolName(null)
    }
  }

  return (
    <>
      <header className="card border border-base-300 bg-base-100 shadow-sm">
        <div className="card-body gap-3">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div>
              <h2 className="text-2xl font-semibold">Pool Management</h2>
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
            <div className="flex items-center justify-between gap-2">
              <h3 className="card-title text-lg">Pool List</h3>
              <button
                className={`btn btn-sm btn-outline ${isPoolsLoading ? 'btn-disabled' : ''}`}
                type="button"
                onClick={() => {
                  void handleReloadPools()
                }}
                disabled={isPoolsLoading || Boolean(deletingPoolName)}
              >
                {isPoolsLoading ? 'Refreshing...' : 'Refresh'}
              </button>
            </div>

            {poolsLoadError && (
              <div className="alert alert-error">
                <span>{poolsLoadError}</span>
              </div>
            )}

              <div className="h-[calc(100vh-18rem)] overflow-x-auto rounded-box border border-base-300">
                  <table className="table  table-pin-rows table-zebra">
                <thead>
                  <tr>
                    <th>#</th>
                    <th>Name</th>
                    <th>Image</th>
                    <th className="text-center">Ready/Size</th>
                    <th className="text-center">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {pools.length === 0 ? (
                    <tr>
                      <td className="text-center text-base-content/70" colSpan={5}>
                        {isPoolsLoading ? 'Loading pools...' : 'No pools found.'}
                      </td>
                    </tr>
                  ) : (
                    pools.map((pool, index) => {
                      const poolName = pool.name?.trim() || ''
                      const isDeleting = deletingPoolName === poolName

                      return (
                        <tr key={poolName || `pool-${index}`}>
                          <td>{index + 1}</td>
                          <td className="font-medium">{poolName || '-'}</td>
                          <td className="font-medium ">{pool.image || '-'}</td>
                          <td className="text-center">
                              <div className="badge badge-secondary badge-sm">
                                  {typeof pool.pool?.readySize === 'number' ? pool.pool.readySize : '-'} / {typeof pool.pool?.size === 'number' ? pool.pool.size : '-'}
                              </div>
                          </td>
                          <td className="text-right">
                            <div className="inline-flex items-center gap-2">
                              <button
                                className="btn btn-xs btn-outline"
                                type="button"
                                disabled={!poolName || Boolean(deletingPoolName)}
                                onClick={() => {
                                  navigate(`/pool/${encodeURIComponent(poolName)}`, {
                                    state: {
                                      selectedPool: pool,
                                    },
                                  })
                                }}
                              >
                                View Details
                              </button>
                              <button
                                className={`btn btn-xs btn-error ${isDeleting ? 'btn-disabled' : ''}`}
                                type="button"
                                disabled={!poolName || Boolean(deletingPoolName)}
                                onClick={() => {
                                  setConfirmDeletePoolName(poolName)
                                }}
                              >
                                {isDeleting ? 'Deleting...' : 'Delete'}
                              </button>
                            </div>
                          </td>
                        </tr>
                      )
                    })
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
