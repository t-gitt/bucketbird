import { useState } from 'react'
import type { FormEvent } from 'react'

import Button from '../ui/Button'

type CredentialFormProps = {
  credential?: {
    id: string
    name: string
    provider: string
    region: string
    endpoint: string
    useSSL: boolean
  }
  onSubmit: (data: {
    name: string
    provider: string
    region: string
    endpoint: string
    accessKey: string
    secretKey: string
    useSSL: boolean
  }) => void
  onCancel: () => void
  isLoading?: boolean
  error?: string | null
}

const CredentialForm = ({ credential, onSubmit, onCancel, isLoading, error }: CredentialFormProps) => {
  const [name, setName] = useState(credential?.name || '')
  const [provider, setProvider] = useState(credential?.provider || 'MinIO')
  const [region, setRegion] = useState(credential?.region || 'us-east-1')
  const [endpoint, setEndpoint] = useState(credential?.endpoint || '')
  const [accessKey, setAccessKey] = useState('')
  const [secretKey, setSecretKey] = useState('')
  const [useSSL, setUseSSL] = useState(credential?.useSSL ?? true)

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    onSubmit({
      name,
      provider,
      region,
      endpoint,
      accessKey,
      secretKey,
      useSSL,
    })
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-6">
      <div>
        <h3 className="text-lg font-semibold text-slate-900 dark:text-white">
          {credential ? 'Edit credentials' : 'Add new credentials'}
        </h3>
        <p className="text-sm text-slate-500 dark:text-slate-400">
          {credential
            ? `Update the details for "${credential.name}".`
            : 'Connect a bucket-capable storage provider.'}
        </p>
      </div>

      {error && (
        <div className="rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-600 dark:border-red-500/50 dark:bg-red-500/10 dark:text-red-200">
          {error}
        </div>
      )}

      <label className="flex flex-col gap-2">
        <span className="text-sm font-medium text-slate-800 dark:text-slate-200">
          Credential name <span className="text-red-500">*</span>
        </span>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="e.g., My MinIO Server"
          required
          className="form-input rounded-lg border border-slate-300 bg-background-light px-3.5 py-2.5 text-slate-900 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20 dark:border-slate-700 dark:bg-background-dark dark:text-white"
        />
      </label>

      <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
        <label className="flex flex-col gap-2">
          <span className="text-sm font-medium text-slate-800 dark:text-slate-200">
            Provider <span className="text-red-500">*</span>
          </span>
          <select
            value={provider}
            onChange={(e) => setProvider(e.target.value)}
            required
            className="form-select rounded-lg border border-slate-300 bg-background-light px-3.5 py-2.5 text-sm text-slate-900 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20 dark:border-slate-700 dark:bg-background-dark dark:text-white"
          >
            <option>MinIO</option>
            <option>AWS S3</option>
            <option>Wasabi</option>
            <option>Backblaze B2</option>
            <option>DigitalOcean Spaces</option>
            <option>Other S3-compatible</option>
          </select>
        </label>

        <label className="flex flex-col gap-2">
          <span className="text-sm font-medium text-slate-800 dark:text-slate-200">
            Default region <span className="text-red-500">*</span>
          </span>
          <input
            type="text"
            value={region}
            onChange={(e) => setRegion(e.target.value)}
            placeholder="e.g., us-east-1"
            required
            className="form-input rounded-lg border border-slate-300 bg-background-light px-3.5 py-2.5 text-slate-900 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20 dark:border-slate-700 dark:bg-background-dark dark:text-white"
          />
        </label>
      </div>

      <label className="flex flex-col gap-2">
        <span className="text-sm font-medium text-slate-800 dark:text-slate-200">
          Endpoint URL <span className="text-red-500">*</span>
        </span>
        <input
          type="text"
          value={endpoint}
          onChange={(e) => setEndpoint(e.target.value)}
          placeholder="e.g., https://s3.amazonaws.com or http://localhost:9000"
          required
          className="form-input rounded-lg border border-slate-300 bg-background-light px-3.5 py-2.5 text-slate-900 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20 dark:border-slate-700 dark:bg-background-dark dark:text-white"
        />
      </label>

      <label className="flex items-center justify-between gap-4 rounded-lg border border-slate-200 bg-slate-50 px-4 py-3 dark:border-slate-700 dark:bg-slate-800/40">
        <div>
          <span className="text-sm font-medium text-slate-800 dark:text-slate-200">
            Use secure connection (HTTPS)
          </span>
          <p className="text-xs text-slate-500 dark:text-slate-400">
            Disable only when connecting to local or self-hosted endpoints without TLS.
          </p>
        </div>
        <input
          type="checkbox"
          checked={useSSL}
          onChange={(event) => setUseSSL(event.target.checked)}
          className="h-5 w-5 accent-primary"
        />
      </label>

      <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
        <label className="flex flex-col gap-2">
          <span className="text-sm font-medium text-slate-800 dark:text-slate-200">
            Access key ID <span className="text-red-500">*</span>
          </span>
          <input
            type="text"
            value={accessKey}
            onChange={(e) => setAccessKey(e.target.value)}
            placeholder={credential ? '••••••••••••••••••••' : 'Enter access key'}
            required={!credential}
            className="form-input rounded-lg border border-slate-300 bg-background-light px-3.5 py-2.5 text-slate-900 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20 dark:border-slate-700 dark:bg-background-dark dark:text-white"
          />
          {credential && !accessKey && (
            <span className="text-xs text-slate-500 dark:text-slate-400">
              Leave blank to keep current access key
            </span>
          )}
        </label>

        <label className="flex flex-col gap-2">
          <span className="text-sm font-medium text-slate-800 dark:text-slate-200">
            Secret access key <span className="text-red-500">*</span>
          </span>
          <input
            type="password"
            value={secretKey}
            onChange={(e) => setSecretKey(e.target.value)}
            placeholder={credential ? '••••••••••••••••••••' : 'Enter secret key'}
            required={!credential}
            className="form-input rounded-lg border border-slate-300 bg-background-light px-3.5 py-2.5 text-slate-900 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20 dark:border-slate-700 dark:bg-background-dark dark:text-white"
          />
          {credential && !secretKey && (
            <span className="text-xs text-slate-500 dark:text-slate-400">
              Leave blank to keep current secret key
            </span>
          )}
        </label>
      </div>

      <div className="flex justify-end gap-3">
        <Button
          type="button"
          variant="outline"
          onClick={onCancel}
          disabled={isLoading}
          className="border-slate-300 text-slate-700 dark:border-slate-700 dark:text-white"
        >
          Cancel
        </Button>
        <Button type="submit" disabled={isLoading}>
          {isLoading ? 'Saving...' : credential ? 'Save changes' : 'Add credential'}
        </Button>
      </div>
    </form>
  )
}

export default CredentialForm
