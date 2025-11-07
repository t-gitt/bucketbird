import { useQuery } from '@tanstack/react-query'

import { api } from '../api/client'

export const useCredentialBuckets = (credentialId?: string) => {
  return useQuery({
    queryKey: ['credentialBuckets', credentialId],
    queryFn: ({ signal }) => api.getCredentialBuckets(credentialId!, signal),
    enabled: Boolean(credentialId),
    staleTime: 1000 * 60 * 5,
  })
}

export default useCredentialBuckets
