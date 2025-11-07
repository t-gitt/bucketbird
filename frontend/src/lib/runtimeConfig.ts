declare global {
  interface Window {
    __BUCKETBIRD_CONFIG__?: {
      apiUrl?: string
      allowRegistration?: boolean
      enableDemoLogin?: boolean
    }
  }
}

type RuntimeConfig = {
  apiUrl?: string
  allowRegistration?: boolean
  enableDemoLogin?: boolean
}

const FALLBACK_API_URL = 'http://localhost:8080'

function parseBoolean(value: string | undefined, defaultValue: boolean): boolean {
  if (!value) {
    return defaultValue
  }

  const normalized = value.trim().toLowerCase()
  if (normalized === 'false' || normalized === '0' || normalized === 'no') {
    return false
  }
  if (normalized === 'true' || normalized === '1' || normalized === 'yes') {
    return true
  }
  return defaultValue
}

export function getRuntimeConfig(): RuntimeConfig {
  if (typeof window !== 'undefined' && window.__BUCKETBIRD_CONFIG__) {
    return window.__BUCKETBIRD_CONFIG__ ?? {}
  }
  return {}
}

export function resolveApiUrl(): string {
  const runtime = getRuntimeConfig()
  if (runtime.apiUrl && runtime.apiUrl.trim() !== '') {
    return runtime.apiUrl
  }
  const envValue = (import.meta.env.VITE_API_URL as string | undefined)
  return envValue && envValue.trim() !== '' ? envValue : FALLBACK_API_URL
}

export function resolveAllowRegistration(): boolean {
  const runtime = getRuntimeConfig()
  if (typeof runtime.allowRegistration === 'boolean') {
    return runtime.allowRegistration
  }
  const envValue = import.meta.env.VITE_ALLOW_REGISTRATION as string | undefined
  return parseBoolean(envValue, true)
}

export function resolveEnableDemoLogin(): boolean {
  const runtime = getRuntimeConfig()
  if (typeof runtime.enableDemoLogin === 'boolean') {
    return runtime.enableDemoLogin
  }
  const envValue = import.meta.env.VITE_ENABLE_DEMO_LOGIN as string | undefined
  return parseBoolean(envValue, false)
}

export type { RuntimeConfig }

