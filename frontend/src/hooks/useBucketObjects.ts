import { useQuery } from '@tanstack/react-query'

import { api } from '../api/client'

export const useBucketObjects = (bucketId: string | undefined, prefix = '') => {
  return useQuery({
    queryKey: ['bucketObjects', bucketId, prefix],
    queryFn: ({ signal }) => api.getBucketObjects(bucketId ?? '', prefix, signal),
    enabled: Boolean(bucketId),
  })
}

export default useBucketObjects
