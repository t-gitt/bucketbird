import { createContext, useContext, useState, type ReactNode } from 'react'

type BucketModalContextType = {
  showCreateModal: boolean
  openCreateModal: () => void
  closeCreateModal: () => void
}

const BucketModalContext = createContext<BucketModalContextType | undefined>(undefined)

export const BucketModalProvider = ({ children }: { children: ReactNode }) => {
  const [showCreateModal, setShowCreateModal] = useState(false)

  const openCreateModal = () => setShowCreateModal(true)
  const closeCreateModal = () => setShowCreateModal(false)

  return (
    <BucketModalContext.Provider value={{ showCreateModal, openCreateModal, closeCreateModal }}>
      {children}
    </BucketModalContext.Provider>
  )
}

export const useBucketModal = () => {
  const context = useContext(BucketModalContext)
  if (!context) {
    throw new Error('useBucketModal must be used within BucketModalProvider')
  }
  return context
}
