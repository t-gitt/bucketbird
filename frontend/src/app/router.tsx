import { Navigate, createBrowserRouter } from 'react-router-dom'

import {
  BucketPage,
  DashboardPage,
  LoginPage,
  NotFoundPage,
  RecentPage,
  SettingsPage,
  SharedPage,
  TrashPage,
} from '../pages'
import { PrivateRoute } from '../components/auth/PrivateRoute'

export const appRouter = createBrowserRouter([
  {
    path: '/',
    element: <Navigate to="/dashboard" replace />,
  },
  {
    path: '/login',
    element: <LoginPage />,
  },
  {
    path: '/dashboard',
    element: (
      <PrivateRoute>
        <DashboardPage />
      </PrivateRoute>
    ),
  },
  {
    path: '/shared',
    element: (
      <PrivateRoute>
        <SharedPage />
      </PrivateRoute>
    ),
  },
  {
    path: '/recent',
    element: (
      <PrivateRoute>
        <RecentPage />
      </PrivateRoute>
    ),
  },
  {
    path: '/trash',
    element: (
      <PrivateRoute>
        <TrashPage />
      </PrivateRoute>
    ),
  },
  {
    path: '/buckets/:bucketId',
    element: (
      <PrivateRoute>
        <BucketPage />
      </PrivateRoute>
    ),
  },
  {
    path: '/settings/*',
    element: (
      <PrivateRoute>
        <SettingsPage />
      </PrivateRoute>
    ),
  },
  {
    path: '*',
    element: <NotFoundPage />,
  },
])
