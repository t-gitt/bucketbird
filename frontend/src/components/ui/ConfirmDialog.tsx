import type { ReactNode } from 'react'
import { useEffect, useRef } from 'react'

import Button from './Button'

type ConfirmDialogProps = {
  open: boolean
  onClose: () => void
  onConfirm: () => void
  title: string
  description: string | ReactNode
  confirmText?: string
  cancelText?: string
  variant?: 'danger' | 'primary'
  isLoading?: boolean
}

export const ConfirmDialog = ({
  open,
  onClose,
  onConfirm,
  title,
  description,
  confirmText = 'Confirm',
  cancelText = 'Cancel',
  variant = 'primary',
  isLoading = false,
}: ConfirmDialogProps) => {
  const dialogRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (open) {
      // Focus the dialog when it opens
      dialogRef.current?.focus()

      // Prevent body scroll when dialog is open
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = ''
    }

    return () => {
      document.body.style.overflow = ''
    }
  }, [open])

  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && open && !isLoading) {
        onClose()
      }
    }

    document.addEventListener('keydown', handleEscape)
    return () => document.removeEventListener('keydown', handleEscape)
  }, [open, onClose, isLoading])

  if (!open) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/50 backdrop-blur-sm"
        onClick={isLoading ? undefined : onClose}
        aria-hidden="true"
      />

      {/* Dialog */}
      <div
        ref={dialogRef}
        role="dialog"
        aria-modal="true"
        aria-labelledby="dialog-title"
        aria-describedby="dialog-description"
        tabIndex={-1}
        className="relative w-full max-w-md rounded-xl border border-slate-200 bg-white p-6 shadow-xl dark:border-slate-700 dark:bg-slate-900"
      >
        {/* Icon */}
        <div
          className={`mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full ${
            variant === 'danger'
              ? 'bg-red-100 dark:bg-red-900/30'
              : 'bg-primary/10 dark:bg-primary/20'
          }`}
        >
          <span
            className={`material-symbols-outlined text-2xl ${
              variant === 'danger' ? 'text-red-600 dark:text-red-400' : 'text-primary'
            }`}
          >
            {variant === 'danger' ? 'warning' : 'info'}
          </span>
        </div>

        {/* Title */}
        <h2
          id="dialog-title"
          className="mb-2 text-center text-lg font-semibold text-slate-900 dark:text-white"
        >
          {title}
        </h2>

        {/* Description */}
        <div
          id="dialog-description"
          className="mb-6 text-center text-sm text-slate-600 dark:text-slate-300"
        >
          {description}
        </div>

        {/* Actions */}
        <div className="flex gap-3">
          <Button
            type="button"
            variant="outline"
            onClick={onClose}
            disabled={isLoading}
            fullWidth
            className="border-slate-300 text-slate-700 dark:border-slate-700 dark:text-white"
          >
            {cancelText}
          </Button>
          <Button
            type="button"
            variant={variant}
            onClick={onConfirm}
            disabled={isLoading}
            fullWidth
          >
            {isLoading ? 'Please wait...' : confirmText}
          </Button>
        </div>
      </div>
    </div>
  )
}

export default ConfirmDialog
