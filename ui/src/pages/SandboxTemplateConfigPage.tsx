import { useEffect, useState } from 'react'

import { getSandboxTemplateConfig, saveSandboxTemplateConfig } from '../lib/api/config'

export default function SandboxTemplateConfigPage() {
  const [templateText, setTemplateText] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [isSaving, setIsSaving] = useState(false)
  const [loadError, setLoadError] = useState('')
  const [saveError, setSaveError] = useState('')
  const [saveSuccess, setSaveSuccess] = useState('')

  const loadTemplate = async (options?: { keepMessages?: boolean }) => {
    setIsLoading(true)
    setLoadError('')

    if (!options?.keepMessages) {
      setSaveError('')
      setSaveSuccess('')
    }

    try {
      const content = await getSandboxTemplateConfig()
      setTemplateText(content ?? '')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to load sandbox template config'
      setLoadError(message)
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    void loadTemplate()
  }, [])

  const handleReload = async () => {
    await loadTemplate({ keepMessages: true })
  }

  const handleSave = async () => {
    setIsSaving(true)
    setSaveError('')
    setSaveSuccess('')

    try {
      await saveSandboxTemplateConfig(templateText)
      setSaveSuccess('Sandbox template saved successfully.')
      await loadTemplate({ keepMessages: true })
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to save sandbox template config'
      setSaveError(message)
    } finally {
      setIsSaving(false)
    }
  }

  return (
    <>
      <header className="card border border-base-300 bg-base-100 shadow-sm">
        <div className="card-body gap-3">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div>
              <h2 className="text-2xl font-semibold">Sandbox-Template Config</h2>
              <p className="text-sm text-base-content/70">View, edit, and save sandbox ReplicaSet template.</p>
            </div>
          </div>
        </div>
      </header>

      <section id="sandbox-template-config">
        <div className="card border border-base-300 bg-base-100 shadow-sm">
          <div className="card-body gap-4">
            <div className="flex items-center justify-between gap-2">
              <h3 className="card-title text-lg">Sandbox-Template Config</h3>
              <div className="flex items-center gap-2">
                <button
                  className={`btn btn-sm btn-outline ${isLoading ? 'btn-disabled' : ''}`}
                  type="button"
                  onClick={() => {
                    void handleReload()
                  }}
                  disabled={isLoading}
                >
                  {isLoading ? 'Reloading...' : 'Reload'}
                </button>
                <button
                  className={`btn btn-sm btn-primary ${isSaving ? 'btn-disabled' : ''}`}
                  type="button"
                  onClick={() => {
                    void handleSave()
                  }}
                  disabled={isSaving || isLoading}
                >
                  {isSaving ? 'Saving...' : 'Save Template'}
                </button>
              </div>
            </div>

            {(loadError || saveError || saveSuccess) && (
              <div className="space-y-2">
                {loadError && (
                  <div className="alert alert-error">
                    <span>{loadError}</span>
                  </div>
                )}
                {saveError && (
                  <div className="alert alert-error">
                    <span>{saveError}</span>
                  </div>
                )}
                {saveSuccess && (
                  <div className="alert alert-success">
                    <span>{saveSuccess}</span>
                  </div>
                )}
              </div>
            )}

            <label className="form-control w-full">
              <div className="label">
                <span className="label-text">ReplicaSet Template (Go Template + YAML)</span>
              </div>
              <textarea
                className="textarea textarea-sm textarea-bordered min-h-[700px] w-full font-mono text-xs"
                value={templateText}
                onChange={(event) => {
                  setTemplateText(event.target.value)
                  setSaveError('')
                  setSaveSuccess('')
                }}
              />
            </label>
          </div>
        </div>
      </section>
    </>
  )
}
