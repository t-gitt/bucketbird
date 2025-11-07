import { useMutation, useQueryClient } from '@tanstack/react-query'

import { api } from '../api/client'

type CreateCredentialInput = {
  name: string
  provider: string
  region: string
  endpoint: string
  accessKey: string
  secretKey: string
  useSSL: boolean
  logo?: string
}

type UpdateCredentialInput = {
  name: string
  provider: string
  region: string
  endpoint: string
  accessKey: string
  secretKey: string
  useSSL: boolean
  logo?: string
}

export function useCreateCredential() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: CreateCredentialInput) => api.createCredential(input),
    onSuccess: () => {
      // Invalidate credentials query to refetch the list
      queryClient.invalidateQueries({ queryKey: ['credentials'] })
    },
  })
}

export function useUpdateCredential() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: UpdateCredentialInput }) =>
      api.updateCredential(id, input),
    onSuccess: () => {
      // Invalidate credentials query to refetch the list
      queryClient.invalidateQueries({ queryKey: ['credentials'] })
    },
  })
}

export function useDeleteCredential() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => api.deleteCredential(id),
    onSuccess: () => {
      // Invalidate credentials query to refetch the list
      queryClient.invalidateQueries({ queryKey: ['credentials'] })
    },
  })
}

export function useTestCredential() {
  return useMutation({
    mutationFn: (id: string) => api.testCredential(id),
  })
}
