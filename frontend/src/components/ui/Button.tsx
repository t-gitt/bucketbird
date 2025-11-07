import type { ButtonHTMLAttributes } from 'react'
import { forwardRef } from 'react'

import { cn } from '../../lib/cn'

type ButtonVariant = 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger'

const variantClasses: Record<ButtonVariant, string> = {
  primary:
    'bg-primary text-white shadow-sm hover:bg-primary-strong focus-visible:ring-primary/40 disabled:bg-primary/50 disabled:text-white/70 dark:text-slate-900 dark:hover:text-white',
  secondary:
    'bg-white text-slate-800 border border-slate-200 hover:bg-slate-100 dark:bg-surface-dark dark:text-white dark:border-slate-700 dark:hover:bg-slate-800',
  outline:
    'border border-slate-300 text-slate-700 hover:bg-background-muted dark:border-slate-600 dark:text-white dark:hover:bg-surface-dark',
  ghost:
    'text-slate-600 hover:bg-background-muted dark:text-slate-300 dark:hover:bg-surface-dark/80',
  danger:
    'bg-red-500 text-white hover:bg-red-600 focus-visible:ring-red-500/50 disabled:bg-red-400',
}

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: ButtonVariant
  fullWidth?: boolean
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = 'primary', fullWidth = false, children, type = 'button', ...rest }, ref) => {
    return (
      <button
        ref={ref}
        type={type}
        className={cn(
          'inline-flex items-center justify-center gap-2 rounded-lg px-4 py-2 text-sm font-semibold transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed',
          fullWidth ? 'w-full' : 'w-auto',
          variantClasses[variant],
          className,
        )}
        {...rest}
      >
        {children}
      </button>
    )
  },
)

Button.displayName = 'Button'

export default Button
