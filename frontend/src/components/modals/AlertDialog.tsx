import { useEffect } from 'react'
import Button from '../ui/Button'

type AlertDialogProps = {
  isOpen: boolean
  onClose: () => void
  title: string
  message: string
  variant?: 'error' | 'success' | 'info' | 'warning'
}

export const AlertDialog = ({
  isOpen,
  onClose,
  title,
  message,
  variant = 'info',
}: AlertDialogProps) => {
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
    error: {
      icon: 'error',
      iconColor: 'text-red-500',
      bgColor: 'bg-red-50 dark:bg-red-900/20',
    },
    success: {
      icon: 'check_circle',
      iconColor: 'text-green-500',
      bgColor: 'bg-green-50 dark:bg-green-900/20',
    },
    info: {
      icon: 'info',
      iconColor: 'text-blue-500',
      bgColor: 'bg-blue-50 dark:bg-blue-900/20',
    },
    warning: {
      icon: 'warning',
      iconColor: 'text-yellow-500',
      bgColor: 'bg-yellow-50 dark:bg-yellow-900/20',
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
            <div className={`flex-shrink-0 rounded-full p-2 ${styles.bgColor}`}>
              <span className={`material-symbols-outlined text-2xl ${styles.iconColor}`}>
                {styles.icon}
              </span>
            </div>
            <div className="flex-1">
              <h3 className="text-lg font-semibold text-slate-900 dark:text-white">{title}</h3>
              <p className="mt-2 whitespace-pre-wrap text-sm text-slate-600 dark:text-slate-400">{message}</p>
            </div>
          </div>

          <div className="flex justify-end">
            <Button onClick={onClose}>Close</Button>
          </div>
        </div>
      </div>
    </div>
  )
}
