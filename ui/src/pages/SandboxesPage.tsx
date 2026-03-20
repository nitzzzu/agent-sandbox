import { FormEvent, useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'

import { createSandbox, deleteSandbox, listSandboxes } from '../lib/api/sandbox'
import type { CreateSandboxRequest, Sandbox } from '../lib/api/types'

type CreateFormState = {
  name: string
  template: string
  image: string
  timeout: string
}

const initialCreateFormState: CreateFormState = {
  name: '',
  template: '',
  image: '',
  timeout: '',
}

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

function buildCreatePayload(form: CreateFormState): CreateSandboxRequest {
  const payload: CreateSandboxRequest = {}

  const name = form.name.trim()
  const template = form.template.trim()
  const image = form.image.trim()
  const timeout = form.timeout.trim()

  if (name) {
    payload.name = name
  }

  if (template) {
    payload.template = template
  }

  if (image) {
    payload.image = image
  }

  if (timeout) {
    const timeoutNumber = Number.parseInt(timeout, 10)
    if (!Number.isNaN(timeoutNumber)) {
      payload.timeout = timeoutNumber
    }
  }

  return payload
}

export default function SandboxesPage() {
  const navigate = useNavigate()
  const [sandboxes, setSandboxes] = useState<Sandbox[]>([])
  const [isListLoading, setIsListLoading] = useState(false)
  const [isCreateSubmitting, setIsCreateSubmitting] = useState(false)
  const [deletingName, setDeletingName] = useState<string | null>(null)
  const [confirmDeleteName, setConfirmDeleteName] = useState<string | null>(null)
  const [pageError, setPageError] = useState('')
  const [actionError, setActionError] = useState('')
  const [createErrorMessage, setCreateErrorMessage] = useState('')
  const [successMessage, setSuccessMessage] = useState('')
  const [formState, setFormState] = useState<CreateFormState>(initialCreateFormState)
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false)

  const refreshSandboxes = async (options?: { keepMessages?: boolean }) => {
    setIsListLoading(true)
    setPageError('')

    if (!options?.keepMessages) {
      setActionError('')
      setSuccessMessage('')
    }

    try {
      const list = await listSandboxes()
      setSandboxes(list)
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to load sandboxes'
      setPageError(message)
    } finally {
      setIsListLoading(false)
    }
  }

  useEffect(() => {
    void refreshSandboxes()
  }, [])

  const isDeleteDisabled = (sandboxName: string) => !sandboxName || Boolean(deletingName)

  const handleDeleteClick = (sandboxName?: string) => {
    if (!sandboxName) {
      return
    }
    setConfirmDeleteName(sandboxName)
  }

  const handleRefresh = async () => {
    await refreshSandboxes({ keepMessages: true })
  }

  const handleCreateSandbox = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    const timeoutText = formState.timeout.trim()
    if (timeoutText) {
      const parsed = Number.parseInt(timeoutText, 10)
      if (Number.isNaN(parsed)) {
        setCreateErrorMessage('Timeout must be an integer in seconds.')
        setSuccessMessage('')
        return
      }
    }

    setIsCreateSubmitting(true)
    setCreateErrorMessage('')
    setActionError('')
    setSuccessMessage('')

    try {
      await createSandbox(buildCreatePayload(formState))
      setFormState(initialCreateFormState)
      setIsCreateModalOpen(false)
      setSuccessMessage('Sandbox created successfully.')
      await refreshSandboxes({ keepMessages: true })
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to create sandbox'
      setCreateErrorMessage(message)
    } finally {
      setIsCreateSubmitting(false)
    }
  }

  const handleDeleteSandbox = async (name?: string) => {
    if (!name) {
      setActionError('Sandbox name is missing and cannot be deleted.')
      setSuccessMessage('')
      setConfirmDeleteName(null)
      return
    }

    setDeletingName(name)
    setActionError('')
    setSuccessMessage('')

    try {
      await deleteSandbox(name)
      setSuccessMessage(`Sandbox "${name}" deleted successfully.`)
      await refreshSandboxes({ keepMessages: true })
    } catch (error) {
      const message = error instanceof Error ? error.message : `Failed to delete sandbox "${name}"`
      setActionError(message)
    } finally {
      setDeletingName(null)
      setConfirmDeleteName(null)
    }
  }

  return (
    <>
      <header className="card border border-base-300 bg-base-100 shadow-sm">
        <div className="card-body gap-3">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div>
              <h2 className="text-2xl font-semibold">Sandbox Management </h2>
              <p className="text-sm text-base-content/70">Manage sandbox instances with list, creation, and deletion.</p>
            </div>
          </div>
        </div>
      </header>

      {(actionError || successMessage) && (
        <section className="space-y-2">
          {actionError && (
            <div className="alert alert-error">
              <span>{actionError}</span>
            </div>
          )}
          {successMessage && (
            <div className="alert alert-success">
              <span>{successMessage}</span>
            </div>
          )}
        </section>
      )}

      <section>
        <div className="card border border-base-300 bg-base-100 shadow-sm">
          <div className="card-body gap-4">
            <div className="flex items-center justify-between gap-2">
              <h3 className="card-title text-lg">Sandboxes List</h3>
              <div className="flex items-center gap-2">
                <button
                  className={`btn btn-sm btn-outline ${isListLoading ? 'btn-disabled' : ''}`}
                  type="button"
                  onClick={() => {
                    void handleRefresh()
                  }}
                  disabled={isListLoading}
                >
                  {isListLoading ? 'Refreshing...' : 'Refresh'}
                </button>
                <button
                  className="btn btn-sm btn-primary"
                  type="button"
                  onClick={() => {
                    setActionError('')
                    setCreateErrorMessage('')
                    setSuccessMessage('')
                    setIsCreateModalOpen(true)
                  }}
                >
                  Create Sandbox
                </button>
              </div>
            </div>

            {pageError && (
              <div className="alert alert-error">
                <span>{pageError}</span>
              </div>
            )}

              <div className="h-[calc(100vh-18rem)] overflow-x-auto rounded-box border border-base-300">
                  <table className="table  table-pin-rows table-zebra">
                <thead>
                  <tr>
                    <th>#</th>
                    <th>Name</th>
                    <th>Template</th>
                    <th>Image</th>
                    <th>Status</th>
                    <th>Timeout</th>
                    <th>Created</th>
                    <th className="text-center">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {sandboxes.length === 0 ? (
                    <tr>
                      <td className="text-center text-base-content/70" colSpan={8}>
                        {isListLoading ? 'Loading sandboxes...' : 'No sandboxes found.'}
                      </td>
                    </tr>
                  ) : (
                    sandboxes.map((sandbox, index) => {
                      const sandboxName = sandbox.name ?? ''
                      const isDeleting = deletingName === sandboxName

                      return (
                        <tr key={sandboxName || sandbox.id || `sandbox-${index}`}>
                          <td>{index + 1}</td>
                            <td className="font-medium">
                                <a className="link link-hover link-success" href={`/sandbox/${sandbox.name}/health`} target="_blank">{sandbox.name || '-'}</a>
                                <div className="text-xs text-base-content/40 font-normal">
                                    <a className="link link-hover" href={`/sandboxes/router/${sandbox.id}/${sandbox.port}/health`} target="_blank">{sandbox.id || '-'}</a>
                                </div>
                            </td>
                            <td>{sandbox.template || '-'}</td>
                          <td>{sandbox.image || '-'}</td>
                          <td>
                            {sandbox.status ? (
                              <span className={`badge badge-sm  ${sandbox.status === 'running' ? 'badge-success' : 'badge-warning'}`}>{sandbox.status}</span>
                            ) : (
                              '-'
                            )}
                          </td>
                          <td>{typeof sandbox.timeout === 'number' ? `${sandbox.timeout}s` : '-'}</td>
                          <td>{formatCreatedAt(sandbox.created_at)}</td>
                          <td className="text-right">
                            <div className="text-center">
                                <div className="mb-2 w-35">
                              <button
                                className="btn btn-xs btn-outline mr-2"
                                type="button"
                                disabled={!sandboxName || Boolean(deletingName)}
                                onClick={() => {
                                  navigate(`/logs?sandbox=${encodeURIComponent(sandboxName)}`)
                                }}
                              >
                                Logs
                              </button>
                              <button
                                className="btn btn-xs btn-outline"
                                type="button"
                                disabled={!sandboxName || Boolean(deletingName)}
                                onClick={() => {
                                  navigate(`/files?sandbox=${encodeURIComponent(sandboxName)}`)
                                }}
                              >
                                Files
                              </button>
                                </div>
                                <div>
                              <button
                                className="btn btn-xs  btn-warning  mr-1"
                                type="button"
                                disabled={!sandboxName || Boolean(deletingName)}
                                onClick={() => {
                                  navigate(`/terminal?sandbox=${encodeURIComponent(sandboxName)}`)
                                }}
                              >
                                Terminal
                              </button>
                              <button
                                className={`btn btn-xs btn-error ${isDeleting ? 'btn-disabled' : ''}`}
                                type="button"
                                disabled={isDeleteDisabled(sandboxName)}
                                onClick={() => {
                                  handleDeleteClick(sandbox.name)
                                }}
                              >
                                {isDeleting ? 'Deleting...' : 'Delete'}
                              </button>
                                </div>
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

      <dialog className={`modal ${isCreateModalOpen ? 'modal-open' : ''}`} open={isCreateModalOpen}>
        <div className="modal-box">
          <h3 className="text-lg font-bold">Create Sandbox</h3>
          <form className="mt-4 space-y-3" onSubmit={(event) => void handleCreateSandbox(event)}>
            <label className="form-control w-full">
              <div className="label">
                <span className="label-text">Name</span>
              </div>
              <input
                className="input input-sm input-bordered w-full"
                type="text"
                value={formState.name}
                onChange={(event) => {
                  setFormState((prev) => ({ ...prev, name: event.target.value }))
                }}
                placeholder="Optional"
              />
            </label>

            <label className="form-control w-full">
              <div className="label">
                <span className="label-text">Template</span>
              </div>
              <input
                className="input input-sm input-bordered w-full"
                type="text"
                value={formState.template}
                onChange={(event) => {
                  setFormState((prev) => ({ ...prev, template: event.target.value }))
                }}
                placeholder="Optional"
              />
            </label>

            <label className="form-control w-full">
              <div className="label">
                <span className="label-text">Image</span>
              </div>
              <input
                className="input input-sm input-bordered w-full"
                type="text"
                value={formState.image}
                onChange={(event) => {
                  setFormState((prev) => ({ ...prev, image: event.target.value }))
                }}
                placeholder="Optional"
              />
            </label>

            <label className="form-control w-full">
              <div className="label">
                <span className="label-text">Timeout (seconds)</span>
              </div>
              <input
                className="input input-sm input-bordered w-full"
                type="text"
                value={formState.timeout}
                onChange={(event) => {
                  setFormState((prev) => ({ ...prev, timeout: event.target.value }))
                }}
                placeholder="Optional"
              />
            </label>

            <div className="modal-action mt-4">
              <button
                className="btn btn-ghost"
                type="button"
                onClick={() => {
                  if (isCreateSubmitting) {
                    return
                  }
                  setCreateErrorMessage('')
                  setIsCreateModalOpen(false)
                }}
                disabled={isCreateSubmitting}
              >
                Cancel
              </button>
              <button className={`btn btn-primary ${isCreateSubmitting ? 'btn-disabled' : ''}`} type="submit" disabled={isCreateSubmitting}>
                {isCreateSubmitting ? 'Creating...' : 'Create Sandbox'}
              </button>
            </div>
          </form>
        </div>
        <button
          className="modal-backdrop"
          type="button"
          aria-label="Close"
          onClick={() => {
            if (isCreateSubmitting) {
              return
            }
            setCreateErrorMessage('')
            setIsCreateModalOpen(false)
          }}
        />
      </dialog>

      <dialog className={`modal ${createErrorMessage ? 'modal-open' : ''}`} open={Boolean(createErrorMessage)}>
        <div className="modal-box relative">
          <button
            className="btn btn-sm btn-circle btn-ghost absolute right-2 top-2"
            type="button"
            aria-label="Close"
            onClick={() => {
              setCreateErrorMessage('')
            }}
          >
            ✕
          </button>
          <h3 className="text-lg font-bold">Create Sandbox Failed</h3>
          <p className="py-4">{createErrorMessage}</p>
        </div>
        <button
          className="modal-backdrop"
          type="button"
          aria-label="Close"
          onClick={() => {
            setCreateErrorMessage('')
          }}
        />
      </dialog>

      <dialog className={`modal ${confirmDeleteName ? 'modal-open' : ''}`} open={Boolean(confirmDeleteName)}>
        <div className="modal-box">
          <h3 className="text-lg font-bold">Delete Sandbox</h3>
          <p className="py-4">
            Are you sure you want to delete sandbox <span className="font-semibold">"{confirmDeleteName}"</span>?
          </p>
          <div className="modal-action">
            <button
              className="btn btn-ghost"
              type="button"
              onClick={() => {
                if (deletingName) {
                  return
                }
                setConfirmDeleteName(null)
              }}
              disabled={Boolean(deletingName)}
            >
              Cancel
            </button>
            <button
              className={`btn btn-error ${deletingName ? 'btn-disabled' : ''}`}
              type="button"
              disabled={Boolean(deletingName)}
              onClick={() => {
                void handleDeleteSandbox(confirmDeleteName ?? undefined)
              }}
            >
              {deletingName ? 'Deleting...' : 'Delete'}
            </button>
          </div>
        </div>
        <button
          className="modal-backdrop"
          type="button"
          aria-label="Close"
          onClick={() => {
            if (deletingName) {
              return
            }
            setConfirmDeleteName(null)
          }}
        />
      </dialog>
    </>
  )
}
