import { useQuery } from '@tanstack/react-query'

import { api } from '../api/client'

export const useCredentials = () => {
  return useQuery({
    queryKey: ['credentials'],
    queryFn: ({ signal }) => api.getCredentials(signal),
  })
}

export default useCredentials
