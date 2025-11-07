import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../contexts/useAuth'
import { ThemeToggle } from '../components/theme/ThemeToggle'
import { resolveAllowRegistration, resolveEnableDemoLogin } from '../lib/runtimeConfig'

const backgroundImage =
  'https://lh3.googleusercontent.com/aida-public/AB6AXuAhhTC9oOiVT3yal95oRfmzGA0VwsSWtkM6TEcP8yskGSpCxrkn3QxMsWu3myeFEFdALDe9jIrcj9NZoOAaXYCgQtVtGSVplNQnf5_VKV_UlDvznJEzMB9sHMek-i-wuTQ6ACkHRwDB57Js2biDS9a6PwESziK4_rM0tIGneJbSCfwY0sUSpWHzZKBeViolpOd0vqh22JP30j20AVBkJhHy00Ds8LabErZcZJOaZtsonzsXan9uwsz-FOUqL_6OMWuCM1AE2jEW_w'

const LoginPage = () => {
  const allowRegistration = resolveAllowRegistration()
  const enableDemoLogin = resolveEnableDemoLogin()
  const [mode, setMode] = useState<'login' | 'register'>('login')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [firstName, setFirstName] = useState('')
  const [lastName, setLastName] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  const { login, demoLogin, register, isAuthenticated } = useAuth()
  const navigate = useNavigate()
  const isRegisterMode = allowRegistration && mode === 'register'

  // Redirect if already authenticated
  useEffect(() => {
    if (isAuthenticated) {
      navigate('/dashboard', { replace: true })
    }
  }, [isAuthenticated, navigate])

  useEffect(() => {
    if (!allowRegistration && mode === 'register') {
      setMode('login')
    }
  }, [allowRegistration, mode])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setIsLoading(true)

    if (!allowRegistration && mode === 'register') {
      setError('Registration is currently disabled.')
      setMode('login')
      setIsLoading(false)
      return
    }

    try {
      if (mode === 'login' || !allowRegistration) {
        await login(email, password)
      } else {
        await register(email, password, firstName, lastName)
      }
      // Navigation handled by useEffect when isAuthenticated changes
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setIsLoading(false)
    }
  }

  const handleModeChange = (newMode: 'login' | 'register') => {
    if (newMode === 'register' && !allowRegistration) {
      return
    }
    setMode(newMode)
    setError('')
  }

  const handleDemoLogin = async () => {
    setError('')
    setIsLoading(true)
    try {
      await demoLogin()
      // Navigation handled by useEffect when isAuthenticated changes
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Demo login failed')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="relative flex min-h-screen items-center justify-center bg-gradient-to-br from-background-light via-white to-background-muted p-4 font-display dark:bg-gradient-to-br dark:from-background-dark dark:via-surface-dark dark:to-background-dark sm:p-6">
      {/* Theme Toggle Button */}
      <div className="absolute right-4 top-4 z-20 sm:right-6 sm:top-6">
        <ThemeToggle />
      </div>

      <div
        className="absolute inset-0 -z-10 opacity-10"
        style={{
          backgroundImage: `url("${backgroundImage}")`,
          backgroundSize: 'cover',
          backgroundPosition: 'center',
          filter: 'blur(40px) saturate(1.5)',
        }}
      />

      <div className="relative z-10 flex w-full max-w-lg flex-col items-center gap-8">
        <div className="flex flex-col items-center gap-3">
          <img
            src="/bucketbird.png"
            alt="BucketBird"
            className="h-32 w-32 rounded-lg border border-accent/40 bg-white p-2 shadow-md dark:border-accent-soft/50 dark:bg-white/10 dark:[filter:drop-shadow(0_0_1px_rgba(255,255,255,0.1))_drop-shadow(0_0_2px_rgba(255,255,255,0.1))]"
          />
          <div className="flex flex-col">
            <span className="text-3xl font-bold text-slate-800 dark:text-white">BucketBird</span>
          </div>
        </div>

        <div className="w-full rounded-xl border border-slate-200/60 bg-white shadow-xl backdrop-blur-lg dark:border-slate-700/60 dark:bg-surface-dark/95 dark:shadow-black/30">
          <div className="flex flex-col">
            <div className="flex p-2">
              <div className="flex h-10 flex-1 items-center justify-center gap-2 rounded-lg bg-slate-100 p-1 dark:bg-surface-dark/70">
                <button
                  type="button"
                  className={`flex h-full flex-1 items-center justify-center rounded-md px-2 text-sm font-medium transition-all ${
                    mode === 'login'
                      ? 'bg-white text-slate-900 shadow-sm dark:bg-slate-900/80 dark:text-white'
                      : 'text-slate-500 hover:text-slate-700 dark:text-slate-400 dark:hover:text-slate-200'
                  }`}
                  onClick={() => handleModeChange('login')}
                >
                  Login
                </button>
                {allowRegistration && (
                  <button
                    type="button"
                    className={`flex h-full flex-1 items-center justify-center rounded-md px-2 text-sm font-medium transition-all ${
                      isRegisterMode
                        ? 'bg-white text-slate-900 shadow-sm dark:bg-slate-900/80 dark:text-white'
                        : 'text-slate-500 hover:text-slate-700 dark:text-slate-400 dark:hover:text-slate-200'
                    }`}
                    onClick={() => handleModeChange('register')}
                  >
                    Register
                  </button>
                )}
              </div>
            </div>

            {!allowRegistration && (
              <p className="px-6 text-center text-xs text-slate-500 dark:text-slate-400">
                Self-service registration is disabled. Contact an administrator if you need an account.
              </p>
            )}

            <div className="flex flex-col gap-6 px-6 pb-8 pt-4">
              <header className="flex flex-col gap-2">
                <p className="text-2xl font-bold leading-tight text-slate-900 dark:text-white">
                  {isRegisterMode ? 'Create your account' : 'Welcome Back'}
                </p>
                <p className="text-sm text-slate-500 dark:text-slate-400">
                  {isRegisterMode
                    ? 'Create an account to start exploring your storage.'
                    : 'Log in to manage your buckets with ease.'}
                </p>
              </header>

              {error && (
                <div className="rounded-lg bg-red-50 p-4 dark:bg-red-900/20">
                  <p className="text-sm text-red-800 dark:text-red-200">{error}</p>
                </div>
              )}

              <form className="flex flex-col gap-4" onSubmit={handleSubmit}>
                {isRegisterMode && (
                  <div className="flex flex-col gap-4">
                    <label className="flex flex-1 flex-col gap-2">
                      <span className="text-sm font-medium text-slate-800 dark:text-slate-200">First Name</span>
                      <input
                        type="text"
                        placeholder="First name"
                        value={firstName}
                        onChange={(e) => setFirstName(e.target.value)}
                        className="form-input rounded-lg border border-slate-300 bg-white px-4 py-2.5 text-sm text-slate-900 placeholder:text-slate-400 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20 dark:border-slate-700 dark:bg-background-dark dark:text-white dark:placeholder:text-slate-500"
                      />
                    </label>
                    <label className="flex flex-1 flex-col gap-2">
                      <span className="text-sm font-medium text-slate-800 dark:text-slate-200">Last Name</span>
                      <input
                        type="text"
                        placeholder="Last name"
                        value={lastName}
                        onChange={(e) => setLastName(e.target.value)}
                        className="form-input rounded-lg border border-slate-300 bg-white px-4 py-2.5 text-sm text-slate-900 placeholder:text-slate-400 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20 dark:border-slate-700 dark:bg-background-dark dark:text-white dark:placeholder:text-slate-500"
                      />
                    </label>
                  </div>
                )}

                <label className="flex flex-col gap-2">
                  <span className="text-sm font-medium text-slate-800 dark:text-slate-200">Email Address</span>
                  <input
                    type="email"
                    placeholder="Enter your email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                    className="form-input rounded-lg border border-slate-300 bg-white px-4 py-2.5 text-sm text-slate-900 placeholder:text-slate-400 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20 dark:border-slate-600 dark:bg-surface-dark/80 dark:text-white dark:placeholder:text-slate-400"
                  />
                </label>

                <label className="flex flex-col gap-2">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium text-slate-800 dark:text-slate-200">Password</span>
                    {mode === 'login'}
                  </div>
                  <div className="flex items-stretch">
                    <input
                      type={showPassword ? 'text' : 'password'}
                      placeholder="Enter your password"
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      required
                      className="form-input flex-1 rounded-l-lg border border-r-0 border-slate-300 bg-white px-4 py-2.5 text-sm text-slate-900 placeholder:text-slate-400 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20 dark:border-slate-600 dark:bg-surface-dark/80 dark:text-white dark:placeholder:text-slate-400"
                    />
                    <button
                      type="button"
                      onClick={() => setShowPassword(!showPassword)}
                      className="flex items-center justify-center rounded-r-lg border border-l-0 border-slate-300 bg-white px-3 text-slate-500 dark:border-slate-600 dark:bg-surface-dark/80 dark:text-slate-300"
                    >
                      <span className="material-symbols-outlined text-xl">
                        {showPassword ? 'visibility_off' : 'visibility'}
                      </span>
                    </button>
                  </div>
                </label>

                <button
                  type="submit"
                  disabled={isLoading}
                  className="mt-2 flex h-11 w-full items-center justify-center rounded-lg bg-primary text-base font-semibold text-white shadow-sm transition-colors hover:bg-primary-strong focus:outline-none focus:ring-2 focus:ring-primary/40 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-primary/80 dark:hover:bg-primary-strong"
                >
                  {isLoading ? (
                    <div className="flex items-center gap-2">
                      <div className="h-5 w-5 animate-spin rounded-full border-2 border-solid border-white border-r-transparent"></div>
                      <span>{isRegisterMode ? 'Creating account...' : 'Logging in...'}</span>
                    </div>
                  ) : (
                    <span>{isRegisterMode ? 'Create account' : 'Login'}</span>
                  )}
                </button>
              </form>

              {mode === 'login' && enableDemoLogin && (
                <div className="flex items-center gap-3">
                  <div className="h-px flex-1 bg-slate-200 dark:bg-slate-700" />
                  <span className="text-xs text-slate-500 dark:text-slate-400">OR</span>
                  <div className="h-px flex-1 bg-slate-200 dark:bg-slate-700" />
                </div>
              )}

              {mode === 'login' && enableDemoLogin && (
                <button
                  type="button"
                  onClick={handleDemoLogin}
                  disabled={isLoading}
                  className="flex h-11 w-full items-center justify-center gap-2 rounded-lg border-2 border-primary bg-white text-base font-semibold text-primary shadow-sm transition-colors hover:bg-primary/5 focus:outline-none focus:ring-2 focus:ring-primary/40 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-surface-dark/80 dark:hover:bg-primary/10"
                >
                  <span className="material-symbols-outlined text-xl">rocket_launch</span>
                  <span>Try Demo</span>
                </button>
              )}

              <p className="text-center text-sm text-slate-500 dark:text-slate-400">
                {allowRegistration ? (
                  mode === 'login' ? (
                    <>
                      Don't have an account?{' '}
                      <button
                        className="font-semibold text-primary hover:underline"
                        type="button"
                        onClick={() => handleModeChange('register')}
                      >
                        Register
                      </button>
                    </>
                  ) : (
                    <>
                      Already joined?{' '}
                      <button
                        className="font-semibold text-primary hover:underline"
                        type="button"
                        onClick={() => handleModeChange('login')}
                      >
                        Log in
                      </button>
                    </>
                  )
                ) : (
                  <>Need an account? Please contact your administrator.</>
                )}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default LoginPage
