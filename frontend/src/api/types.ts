export type BucketSummary = {
  id: string
  name: string
  region: string
  description?: string | null
  credentialId: string
  credentialName: string
  credentialProvider: string
  createdAt: string
  size: string
}

export type BucketObject = {
  key: string
  name: string
  kind: 'folder' | 'file'
  icon: string
  iconColor: string
  lastModified: string
  size: string
}

export type CredentialSet = {
  id: string
  name: string
  provider: string
  region: string
  endpoint: string
  useSSL: boolean
  status: 'Active' | 'Paused'
  logo?: string
}

export type UserProfile = {
  id: string
  firstName: string
  lastName: string
  email: string
}
