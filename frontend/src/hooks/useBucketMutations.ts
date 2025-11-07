import { useMutation, useQueryClient } from '@tanstack/react-query'

import { api, type CreateBucketInput } from '../api/client'

export function useCreateBucket() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: CreateBucketInput) => api.createBucket(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['buckets'] })
    },
  })
}

type DeleteBucketInput = {
  bucketId: string
  deleteRemote?: boolean
}

export function useDeleteBucket() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ bucketId, deleteRemote }: DeleteBucketInput) =>
      api.deleteBucket(bucketId, { deleteRemote }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['buckets'] })
    },
  })
}

export function useRecalculateBucketSize() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (bucketId: string) => api.recalculateBucketSize(bucketId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['buckets'] })
    },
  })
}
