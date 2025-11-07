import { useQuery } from '@tanstack/react-query'

import { api } from '../api/client'

export const useBuckets = () => {
  return useQuery({
    queryKey: ['buckets'],
    queryFn: ({ signal }) => api.getBuckets(signal),
  })
}

export default useBuckets
