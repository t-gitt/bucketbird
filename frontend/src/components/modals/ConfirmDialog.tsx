import { useEffect, type ReactNode } from 'react'
import Button from '../ui/Button'

type ConfirmDialogProps = {
  isOpen: boolean
  onClose: () => void
  onConfirm: () => void
  title: string
  message: ReactNode
  confirmText?: string
  cancelText?: string
  variant?: 'danger' | 'warning' | 'info'
  isLoading?: boolean
}

export const ConfirmDialog = ({
  isOpen,
  onClose,
  onConfirm,
  title,
  message,
  confirmText = 'Confirm',
  cancelText = 'Cancel',
  variant = 'danger',
  isLoading = false,
}: ConfirmDialogProps) => {
  useEffect(() => {
    if (!isOpen) return

    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        onClose()
      }
    }

    window.addEventListener('keydown', handleEscape)
    return () => window.removeEventListener('keydown', handleEscape)
  }, [isOpen, onClose])

  if (!isOpen) return null

  const variantStyles = {
    danger: {
      icon: 'warning',
      iconColor: 'text-red-500',
      confirmButton: 'bg-red-600 hover:bg-red-700 text-white',
    },
    warning: {
      icon: 'warning',
      iconColor: 'text-yellow-500',
      confirmButton: 'bg-yellow-600 hover:bg-yellow-700 text-white',
    },
    info: {
      icon: 'info',
      iconColor: 'text-blue-500',
      confirmButton: 'bg-primary hover:bg-primary-strong text-white',
    },
  }

  const styles = variantStyles[variant]

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
      <div
        className="relative w-full max-w-md rounded-lg bg-white shadow-xl dark:bg-slate-900"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex flex-col gap-4 p-6">
          <div className="flex items-start gap-4">
            <div className={`flex-shrink-0 ${styles.iconColor}`}>
              <span className="material-symbols-outlined text-3xl">{styles.icon}</span>
            </div>
            <div className="flex-1">
              <h3 className="text-lg font-semibold text-slate-900 dark:text-white">{title}</h3>
              <div className="mt-2 text-sm text-slate-600 dark:text-slate-400">{message}</div>
            </div>
          </div>

          <div className="flex justify-end gap-3">
            <Button variant="outline" onClick={onClose}>
              {cancelText}
            </Button>
            <button
              onClick={() => {
                if (isLoading) {
                  return
                }
                onConfirm()
                onClose()
              }}
              disabled={isLoading}
              className={`rounded-lg px-4 py-2 text-sm font-medium transition-colors ${styles.confirmButton} ${
                isLoading ? 'cursor-not-allowed opacity-70' : ''
              }`}
            >
              {isLoading ? 'Please wait…' : confirmText}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
