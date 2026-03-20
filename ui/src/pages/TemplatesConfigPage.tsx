import { useEffect, useState } from 'react'

import { getTemplatesConfig, saveTemplatesConfig } from '../lib/api/config'
import type { Template } from '../lib/api/types'

type EditableTemplatePool = {
  size?: number | string
  readySize?: number | string
  probePort?: number | string
  warmupCmd?: string
  startupCmd?: string
}

type EditableTemplate = Omit<Template, 'port' | 'pool'> & {
  port?: number | string
  pool?: EditableTemplatePool
}

function parseOptionalInteger(value: unknown, fieldName: string): number | undefined {
  if (value === undefined || value === null || value === '') {
    return undefined
  }

  if (typeof value === 'number') {
    if (Number.isInteger(value)) {
      return value
    }
    throw new Error(`${fieldName} must be an integer.`)
  }

  if (typeof value === 'string') {
    const trimmed = value.trim()
    if (!trimmed) {
      return undefined
    }

    const parsed = Number.parseInt(trimmed, 10)
    if (!Number.isNaN(parsed)) {
      return parsed
    }
  }

  throw new Error(`${fieldName} must be an integer.`)
}

function toSaveTemplatesPayload(templates: EditableTemplate[]): Template[] {
  if (!Array.isArray(templates)) {
    throw new Error('Templates must be an array.')
  }

  return templates.map((template, index) => {
    const pool = template.pool

    return {
      name: template.name?.trim() || undefined,
      description: template.description?.trim() || undefined,
      image: template.image?.trim() || undefined,
      type: template.type?.trim() || undefined,
      pattern: template.pattern?.trim() || undefined,
      noStartupProbe: Boolean(template.noStartupProbe),
      port: parseOptionalInteger(template.port, `Template #${index + 1} port`),
      pool: {
        size: parseOptionalInteger(pool?.size, `Template #${index + 1} pool.size`),
        readySize: parseOptionalInteger(pool?.readySize, `Template #${index + 1} pool.readySize`),
        probePort: parseOptionalInteger(pool?.probePort, `Template #${index + 1} pool.probePort`),
        warmupCmd: pool?.warmupCmd?.trim() || undefined,
        startupCmd: pool?.startupCmd?.trim() || undefined,
      },
    }
  })
}

export default function TemplatesConfigPage() {
  const [templates, setTemplates] = useState<EditableTemplate[]>([])
  const [selectedTemplateIndex, setSelectedTemplateIndex] = useState<number | null>(null)
  const [editorMode, setEditorMode] = useState<'form' | 'raw'>('form')
  const [rawTemplatesText, setRawTemplatesText] = useState('[]')
  const [isTemplatesLoading, setIsTemplatesLoading] = useState(false)
  const [isTemplatesSaving, setIsTemplatesSaving] = useState(false)
  const [templatesLoadError, setTemplatesLoadError] = useState('')
  const [templatesParseError, setTemplatesParseError] = useState('')
  const [templatesSaveError, setTemplatesSaveError] = useState('')
  const [templatesSaveSuccess, setTemplatesSaveSuccess] = useState('')

  const loadTemplates = async (options?: { keepMessages?: boolean }) => {
    setIsTemplatesLoading(true)
    setTemplatesLoadError('')
    setTemplatesParseError('')

    if (!options?.keepMessages) {
      setTemplatesSaveError('')
      setTemplatesSaveSuccess('')
    }

    try {
      const data = await getTemplatesConfig()
      const raw = typeof data === 'string' ? data.trim() : ''

      if (!raw) {
        setTemplates([])
        setRawTemplatesText('[]')
        setSelectedTemplateIndex(null)
        return
      }

      let parsed: unknown
      try {
        parsed = JSON.parse(raw)
      } catch {
        throw new Error('Failed to parse templates config JSON returned by backend.')
      }

      if (!Array.isArray(parsed)) {
        throw new Error('Templates config JSON must be an array.')
      }

      const nextTemplates = parsed as EditableTemplate[]
      setTemplates(nextTemplates)
      setRawTemplatesText(JSON.stringify(nextTemplates, null, 2))
      setSelectedTemplateIndex((prev) => {
        if (nextTemplates.length === 0) {
          return null
        }

        if (prev === null || prev >= nextTemplates.length) {
          return 0
        }

        return prev
      })
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to load templates config'
      if (message.includes('parse') || message.includes('array')) {
        setTemplatesParseError(message)
      } else {
        setTemplatesLoadError(message)
      }
      setTemplates([])
      setRawTemplatesText('[]')
      setSelectedTemplateIndex(null)
    } finally {
      setIsTemplatesLoading(false)
    }
  }

  useEffect(() => {
    void loadTemplates()
  }, [])

  const handleReloadTemplates = async () => {
    await loadTemplates({ keepMessages: true })
  }

  const clearSaveMessages = () => {
    setTemplatesSaveError('')
    setTemplatesSaveSuccess('')
  }

  const handleSelectTemplate = (index: number) => {
    setSelectedTemplateIndex(index)
    clearSaveMessages()
  }

  const updateSelectedTemplate = (updater: (template: EditableTemplate) => EditableTemplate) => {
    setTemplates((previous) => {
      if (selectedTemplateIndex === null || selectedTemplateIndex < 0 || selectedTemplateIndex >= previous.length) {
        return previous
      }

      const next = [...previous]
      next[selectedTemplateIndex] = updater(next[selectedTemplateIndex] ?? {})
      setRawTemplatesText(JSON.stringify(next, null, 2))
      return next
    })
  }

  const handleSwitchToRawMode = () => {
    setRawTemplatesText(JSON.stringify(templates, null, 2))
    setTemplatesParseError('')
    clearSaveMessages()
    setEditorMode('raw')
  }

  const handleSwitchToFormMode = () => {
    try {
      const parsed = JSON.parse(rawTemplatesText)
      if (!Array.isArray(parsed)) {
        throw new Error('Raw JSON must be an array.')
      }

      const nextTemplates = parsed as EditableTemplate[]
      setTemplates(nextTemplates)
      setSelectedTemplateIndex(nextTemplates.length === 0 ? null : 0)
      setTemplatesParseError('')
      clearSaveMessages()
      setEditorMode('form')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to parse raw templates JSON'
      setTemplatesParseError(message)
    }
  }

  const handleSaveTemplates = async () => {
    setIsTemplatesSaving(true)
    setTemplatesSaveError('')
    setTemplatesSaveSuccess('')

    try {
      const sourceTemplates =
        editorMode === 'raw'
          ? (() => {
              let parsed: unknown
              try {
                parsed = JSON.parse(rawTemplatesText)
              } catch {
                throw new Error('Raw JSON is invalid and cannot be saved.')
              }

              if (!Array.isArray(parsed)) {
                throw new Error('Raw JSON must be an array.')
              }

              return parsed as EditableTemplate[]
            })()
          : templates

      const payload = toSaveTemplatesPayload(sourceTemplates)
      await saveTemplatesConfig(payload)
      setTemplatesSaveSuccess('Templates saved successfully.')
      await loadTemplates({ keepMessages: true })
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to save templates config'
      setTemplatesSaveError(message)
    } finally {
      setIsTemplatesSaving(false)
    }
  }

  const selectedTemplate = selectedTemplateIndex === null ? undefined : templates[selectedTemplateIndex]
  const isRawMode = editorMode === 'raw'

  return (
    <>
      <header className="card border border-base-300 bg-base-100 shadow-sm">
        <div className="card-body gap-3">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div>
              <h2 className="text-2xl font-semibold">Templates Config</h2>
              <p className="text-sm text-base-content/70">View, edit, and save templates configuration.</p>
            </div>
          </div>
        </div>
      </header>

      <section id="templates-config">
        <div className="card border border-base-300 bg-base-100 shadow-sm">
          <div className="card-body gap-4">
            <div className="flex items-center justify-between gap-2">
              <h3 className="card-title text-lg">Templates Config</h3>
                <div className="flex items-center gap-2">
                    <button className={`btn btn-sm btn-dash  ${!isRawMode ? 'btn-warning' : 'btn-success'}`} type="button"
                            onClick={handleSwitchToFormMode} disabled={!isRawMode}>
                        Form
                    </button>
                    <button className={`btn btn-sm btn-dash ${isRawMode ? 'btn-success' : 'btn-warning'}`} type="button"
                            onClick={handleSwitchToRawMode} disabled={isRawMode}>
                        Raw
                    </button>
                    <div className="divider lg:divider-horizontal"></div>
                    <button
                        className={`btn btn-sm btn-outline ${isTemplatesLoading ? 'btn-disabled' : ''}`}
                        type="button"
                        onClick={() => {
                            void handleReloadTemplates()
                        }}
                        disabled={isTemplatesLoading}
                    >
                        {isTemplatesLoading ? 'Reloading...' : 'Reload'}
                    </button>
                    <button
                        className={`btn btn-sm btn-primary ${isTemplatesSaving ? 'btn-disabled' : ''}`}
                        type="button"
                        onClick={() => {
                            void handleSaveTemplates()
                        }}
                        disabled={isTemplatesSaving || isTemplatesLoading}
                    >
                        {isTemplatesSaving ? 'Saving...' : 'Save Templates'}
                    </button>
                </div>
            </div>

              {(templatesLoadError || templatesParseError || templatesSaveError || templatesSaveSuccess) && (
              <div className="space-y-2">
                {templatesLoadError && (
                  <div className="alert alert-error">
                    <span>{templatesLoadError}</span>
                  </div>
                )}
                {templatesParseError && (
                  <div className="alert alert-error">
                    <span>{templatesParseError}</span>
                  </div>
                )}
                {templatesSaveError && (
                  <div className="alert alert-error">
                    <span>{templatesSaveError}</span>
                  </div>
                )}
                {templatesSaveSuccess && (
                  <div className="alert alert-success">
                    <span>{templatesSaveSuccess}</span>
                  </div>
                )}
              </div>
            )}

            {isRawMode ? (
              <label className="form-control w-full">
                <div className="label">
                  <span className="label-text">Templates JSON</span>
                </div>
                <textarea
                  className="textarea textarea-sm textarea-bordered min-h-[650px] w-full font-mono text-xs"
                  value={rawTemplatesText}
                  onChange={(event) => {
                    setRawTemplatesText(event.target.value)
                    setTemplatesParseError('')
                    clearSaveMessages()
                  }}
                />
              </label>
            ) : (
              <div className="grid gap-4 lg:grid-cols-[300px,1fr]">
                <div className="rounded-box border border-base-300 p-2">
                  <div className="mb-2 text-sm font-medium">Templates</div>
                  <div className="custom-scrollbar max-h-[420px] overflow-auto">
                    {templates.length === 0 ? (
                      <div className="px-2 py-3 text-sm text-base-content/70">{isTemplatesLoading ? 'Loading templates...' : 'No templates found.'}</div>
                    ) : (
                      <div className="menu w-full gap-1 p-0">
                        {templates.map((template, index) => {
                          const label = template.name?.trim() || `Template ${index + 1}`
                          const isSelected = index === selectedTemplateIndex

                          return (
                            <button
                              key={`${label}-${index}`}
                              type="button"
                              className={`btn btn-sm justify-start ${isSelected ? 'btn-primary' : 'btn-ghost'}`}
                              onClick={() => {
                                handleSelectTemplate(index)
                              }}
                            >
                              {label}
                            </button>
                          )
                        })}
                      </div>
                    )}
                  </div>
                </div>

                <div className="rounded-box border border-base-300 p-4">
                  {!selectedTemplate ? (
                    <div className="text-sm text-base-content/70">Select a template from the list to edit.</div>
                  ) : (
                    <div className="space-y-3">
                      <div className="grid gap-3 md:grid-cols-2">
                        <label className="form-control w-full">
                          <div className="label">
                            <span className="label-text">Name</span>
                          </div>
                          <input
                            className="input input-sm input-bordered w-full"
                            type="text"
                            value={selectedTemplate.name ?? ''}
                            onChange={(event) => {
                              const value = event.target.value
                              updateSelectedTemplate((prev) => ({ ...prev, name: value }))
                            }}
                          />
                        </label>

                        <label className="form-control w-full">
                          <div className="label">
                            <span className="label-text">Pattern</span>
                          </div>
                          <input
                            className="input input-sm input-bordered w-full"
                            type="text"
                            value={selectedTemplate.pattern ?? ''}
                            onChange={(event) => {
                              const value = event.target.value
                              updateSelectedTemplate((prev) => ({ ...prev, pattern: value }))
                            }}
                          />
                        </label>

                        <label className="form-control w-full md:col-span-2">
                          <div className="label">
                            <span className="label-text">Description</span>
                          </div>
                          <input
                            className="input input-sm input-bordered w-full"
                            type="text"
                            value={selectedTemplate.description ?? ''}
                            onChange={(event) => {
                              const value = event.target.value
                              updateSelectedTemplate((prev) => ({ ...prev, description: value }))
                            }}
                          />
                        </label>

                        <label className="form-control w-full md:col-span-2">
                          <div className="label">
                            <span className="label-text">Image</span>
                          </div>
                          <input
                            className="input input-sm input-bordered w-full"
                            type="text"
                            value={selectedTemplate.image ?? ''}
                            onChange={(event) => {
                              const value = event.target.value
                              updateSelectedTemplate((prev) => ({ ...prev, image: value }))
                            }}
                          />
                        </label>

                        <label className="form-control w-full">
                          <div className="label">
                            <span className="label-text">Type</span>
                          </div>
                          <input
                            className="input input-sm input-bordered w-full"
                            type="text"
                            value={selectedTemplate.type ?? ''}
                            onChange={(event) => {
                              const value = event.target.value
                              updateSelectedTemplate((prev) => ({ ...prev, type: value }))
                            }}
                          />
                        </label>

                        <label className="form-control w-full">
                          <div className="label">
                            <span className="label-text">Port</span>
                          </div>
                          <input
                            className="input input-sm input-bordered w-full"
                            type="text"
                            value={selectedTemplate.port ?? ''}
                            onChange={(event) => {
                              const value = event.target.value
                              updateSelectedTemplate((prev) => ({ ...prev, port: value }))
                            }}
                          />
                        </label>

                        <label className="form-control w-full">
                          <div className="label mr-1">
                            <span className="label-text">No Startup Probe</span>
                          </div>
                          <input
                            className="checkbox checkbox-sm"
                            type="checkbox"
                            checked={Boolean(selectedTemplate.noStartupProbe)}
                            onChange={(event) => {
                              const checked = event.target.checked
                              updateSelectedTemplate((prev) => ({ ...prev, noStartupProbe: checked }))
                            }}
                          />
                        </label>

                        <div className="rounded-box border border-base-300 bg-base-200/40 p-3 md:col-span-2">
                          <div className="mb-2 text-sm font-medium">Template Pool</div>
                          <div className="grid gap-3 md:grid-cols-3">
                            <label className="form-control w-full">
                              <div className="label">
                                <span className="label-text">Pool Size</span>
                              </div>
                              <input
                                className="input input-sm input-bordered w-full"
                                type="text"
                                value={selectedTemplate.pool?.size ?? ''}
                                onChange={(event) => {
                                  const value = event.target.value
                                  updateSelectedTemplate((prev) => ({
                                    ...prev,
                                    pool: {
                                      ...prev.pool,
                                      size: value,
                                    },
                                  }))
                                }}
                              />
                            </label>

                            <label className="form-control w-full">
                              <div className="label">
                                <span className="label-text">Pool Probe Port</span>
                              </div>
                              <input
                                className="input input-sm input-bordered w-full"
                                type="text"
                                value={selectedTemplate.pool?.probePort ?? ''}
                                onChange={(event) => {
                                  const value = event.target.value
                                  updateSelectedTemplate((prev) => ({
                                    ...prev,
                                    pool: {
                                      ...prev.pool,
                                      probePort: value,
                                    },
                                  }))
                                }}
                              />
                            </label>
                            <label className="form-control w-full md:col-span-3">
                              <div className="label">
                                <span className="label-text">Pool Warmup Command</span>
                              </div>
                              <input
                                className="input input-sm input-bordered w-full"
                                type="text"
                                value={selectedTemplate.pool?.warmupCmd ?? ''}
                                onChange={(event) => {
                                  const value = event.target.value
                                  updateSelectedTemplate((prev) => ({
                                    ...prev,
                                    pool: {
                                      ...prev.pool,
                                      warmupCmd: value,
                                    },
                                  }))
                                }}
                              />
                            </label>

                            <label className="form-control w-full md:col-span-3">
                              <div className="label">
                                <span className="label-text">Pool Startup Command</span>
                              </div>
                              <input
                                className="input input-sm input-bordered w-full"
                                type="text"
                                value={selectedTemplate.pool?.startupCmd ?? ''}
                                onChange={(event) => {
                                  const value = event.target.value
                                  updateSelectedTemplate((prev) => ({
                                    ...prev,
                                    pool: {
                                      ...prev.pool,
                                      startupCmd: value,
                                    },
                                  }))
                                }}
                              />
                            </label>
                          </div>
                        </div>
                      </div>
                    </div>
                  )}
                </div>
              </div>
            )}
          </div>
        </div>
      </section>
    </>
  )
}
