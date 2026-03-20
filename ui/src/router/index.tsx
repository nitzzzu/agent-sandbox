import { Navigate, Outlet, createHashRouter } from 'react-router-dom'

import AppShellLayout from '../layouts/AppShellLayout'
import { hasAuthToken } from '../lib/auth/token'
import FilesPage from '../pages/FilesPage'
import LoginPage from '../pages/LoginPage'
import LogsPage from '../pages/LogsPage'
import PoolDetailPage from '../pages/PoolDetailPage'
import PoolListPage from '../pages/PoolListPage'
import SandboxesPage from '../pages/SandboxesPage'
import TemplatesConfigPage from '../pages/TemplatesConfigPage'
import TerminalPage from '../pages/TerminalPage'

function RequireAuth() {
  if (!hasAuthToken()) {
    return <Navigate to="/login" replace />
  }
  return <Outlet />
}

function RedirectIfAuthed() {
  if (hasAuthToken()) {
    return <Navigate to="/sandboxes" replace />
  }
  return <LoginPage />
}

export const appRouter = createHashRouter([
  {
    path: '/login',
    element: <RedirectIfAuthed />,
  },
  {
    element: <RequireAuth />,
    children: [
      {
        path: '/',
        element: <AppShellLayout />,
        children: [
          {
            index: true,
            element: <Navigate to="sandboxes" replace />,
          },
          {
            path: 'sandboxes',
            element: <SandboxesPage />,
          },
          {
            path: 'pool',
            element: <PoolListPage />,
          },
          {
            path: 'logs',
            element: <LogsPage />,
          },
          {
            path: 'terminal',
            element: <TerminalPage />,
          },
          {
            path: 'files',
            element: <FilesPage />,
          },
          {
            path: 'pool/:poolName',
            element: <PoolDetailPage />,
          },
          {
            path: 'config/templates',
            element: <TemplatesConfigPage />,
          },
        ],
      },
    ],
  },
  {
    path: '*',
    element: <Navigate to="/sandboxes" replace />,
  },
])
