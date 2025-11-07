import { useState } from 'react'

type CredentialCardProps = {
  credential: {
    id: string
    name: string
    provider: string
    region: string
    endpoint: string
    useSSL: boolean
    status: string
    logo?: string
  }
  onEdit: () => void
  onDelete: () => void
  onTest: () => void
  isTestLoading?: boolean
  testResult?: { success: boolean; message: string } | null | undefined
}

const CredentialCard = ({
  credential,
  onEdit,
  onDelete,
  onTest,
  isTestLoading,
  testResult,
}: CredentialCardProps) => {
  const [showMenu, setShowMenu] = useState(false)

  return (
    <article
      className={`flex flex-col gap-3 rounded-lg border p-4 ${
        credential.status === 'active'
          ? 'border-primary/40 bg-primary/10 dark:bg-primary/15'
          : 'border-slate-200 bg-white dark:border-slate-700 dark:bg-slate-900/60'
      }`}
    >
      <div className="flex items-start justify-between gap-4">
        <div className="flex items-center gap-3">
          <div className="flex h-9 w-9 items-center justify-center rounded bg-slate-200 dark:bg-slate-700">
            <span className="material-symbols-outlined text-lg text-slate-600 dark:text-slate-300">
              {credential.logo || 'storage'}
            </span>
          </div>
          <div>
            <h3 className="text-base font-semibold text-slate-900 dark:text-white">
              {credential.name}
            </h3>
            <p className="text-sm text-slate-500 dark:text-slate-400">{credential.region}</p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <span
            className={`rounded-full px-2.5 py-0.5 text-xs font-medium ${
              credential.status === 'active'
                ? 'bg-green-100 text-green-800 dark:bg-green-900/40 dark:text-green-300'
                : 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-300'
            }`}
          >
            {credential.status === 'active' ? 'Active' : 'Inactive'}
          </span>
          <div className="relative">
            <button
              type="button"
              onClick={() => setShowMenu(!showMenu)}
              className="flex h-8 w-8 items-center justify-center rounded-md text-slate-600 transition-colors hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800"
            >
              <span className="material-symbols-outlined text-lg">more_vert</span>
            </button>
            {showMenu && (
              <>
                <div
                  className="fixed inset-0 z-10"
                  onClick={() => setShowMenu(false)}
                  onKeyDown={(e) => e.key === 'Escape' && setShowMenu(false)}
                  role="button"
                  tabIndex={0}
                  aria-label="Close menu"
                />
                <div className="absolute right-0 top-full z-20 mt-1 w-40 rounded-lg border border-slate-200 bg-white py-1 shadow-lg dark:border-slate-700 dark:bg-slate-800">
                  <button
                    type="button"
                    onClick={() => {
                      onEdit()
                      setShowMenu(false)
                    }}
                    className="flex w-full items-center gap-2 px-4 py-2 text-left text-sm text-slate-700 transition-colors hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-slate-700"
                  >
                    <span className="material-symbols-outlined text-base">edit</span>
                    Edit
                  </button>
                  <button
                    type="button"
                    onClick={() => {
                      onTest()
                      setShowMenu(false)
                    }}
                    disabled={isTestLoading}
                    className="flex w-full items-center gap-2 px-4 py-2 text-left text-sm text-slate-700 transition-colors hover:bg-slate-100 disabled:opacity-50 dark:text-slate-200 dark:hover:bg-slate-700"
                  >
                    <span className="material-symbols-outlined text-base">
                      {isTestLoading ? 'progress_activity' : 'cable'}
                    </span>
                    {isTestLoading ? 'Testing...' : 'Test connection'}
                  </button>
                  <hr className="my-1 border-slate-200 dark:border-slate-700" />
                  <button
                    type="button"
                    onClick={() => {
                      onDelete()
                      setShowMenu(false)
                    }}
                    className="flex w-full items-center gap-2 px-4 py-2 text-left text-sm text-red-600 transition-colors hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-500/10"
                  >
                    <span className="material-symbols-outlined text-base">delete</span>
                    Delete
                  </button>
                </div>
              </>
            )}
          </div>
        </div>
      </div>

      {testResult && (
        <div
          className={`rounded-lg border p-3 text-sm ${
            testResult.success
              ? 'border-green-200 bg-green-50 text-green-800 dark:border-green-500/50 dark:bg-green-500/10 dark:text-green-200'
              : 'border-red-200 bg-red-50 text-red-600 dark:border-red-500/50 dark:bg-red-500/10 dark:text-red-200'
          }`}
        >
          <div className="flex items-center gap-2">
            <span className="material-symbols-outlined text-base">
              {testResult.success ? 'check_circle' : 'error'}
            </span>
            <span>{testResult.message}</span>
          </div>
        </div>
      )}

      <div className="text-xs text-slate-500 dark:text-slate-400">
        <div className="flex items-center justify-between gap-2">
          <div className="flex items-center gap-1">
            <span className="material-symbols-outlined text-sm">link</span>
            <span className="truncate">{credential.endpoint}</span>
          </div>
          <div
            className={`flex items-center gap-1 rounded-full px-2 py-0.5 font-medium ${
              credential.useSSL
                ? 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-300'
                : 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-300'
            }`}
          >
            <span className="material-symbols-outlined text-xs">
              {credential.useSSL ? 'lock' : 'vpn_lock'}
            </span>
            <span className="text-[10px] leading-none">{credential.useSSL ? 'HTTPS' : 'HTTP'}</span>
          </div>
        </div>
      </div>
    </article>
  )
}

export default CredentialCard
