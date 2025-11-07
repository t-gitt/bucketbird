import { RouterProvider } from 'react-router-dom'

import AppProviders from './app/providers'
import { appRouter } from './app/router'

const App = () => {
  return (
    <AppProviders>
      <RouterProvider router={appRouter} />
    </AppProviders>
  )
}

export default App
