# tDrive Frontend

tDrive's frontend is a React 19 + TypeScript single-page application scaffolded with Vite and styled with Tailwind CSS. It implements the S3 drive experience defined in the design templates and is ready to integrate with the Go backend.

## Stack

- React 19 with Vite for fast local development
- TypeScript with strict settings and React Router for routing
- Tailwind CSS (v3) with design tokens matching the supplied templates
- @tanstack/react-query for data fetching and caching
- Lightweight mock API layer (`src/api`) to unblock UI work before the backend is ready
- Dark/light theme support with a custom ThemeProvider and Material Symbols icons

## Getting Started

```bash
cp .env.example .env # adjust VITE_API_URL if your backend runs elsewhere
npm install
npm run dev        # start Vite dev server at http://localhost:5173
npm run build      # type-check and create a production bundle
```

Linting uses ESLint (`npm run lint`).

## Project Layout

```
frontend/
├── src/
│   ├── app/               # Router and global providers
│   ├── api/               # Mock API client + types
│   ├── components/        # Layout, theme, and UI primitives
│   ├── hooks/             # React Query hooks around the API client
│   ├── pages/             # Route-level screens matching the design templates
│   ├── lib/               # Shared utilities (e.g., class name helper)
│   └── index.css          # Tailwind base styles
└── index.html             # Font/material icon setup and theme bootstrapping
```

## Mock Data

`src/api/mockData.ts` provides placeholder buckets, objects, credentials, and user profile data. Hooks in `src/hooks` expose these datasets via React Query so the UI behaves as if it were connected to a backend. Replace these implementations with real API calls as services come online.

## Theming

The UI respects a `tdrive-theme` localStorage key and the user's OS preference. `ThemeProvider` handles toggling dark/light mode and `ThemeToggle` inside the top bar exposes the switch to users.

## Next Steps for Backend Integration

1. **API Contracts** – Define REST/GraphQL responses for buckets, objects, credentials, and profile settings. Replace mock implementations in `src/api/client.ts` with real HTTP calls.
2. **Authentication Flow** – Wire the login/register screen to session endpoints (supporting TOTP/SSO later). Add route guards to redirect unauthenticated users.
3. **Mutations** – Implement create/update/delete mutations for buckets, objects, and credentials. Use React Query `useMutation` hooks and optimistic updates where applicable.
4. **Error & Loading States** – Expand loading indicators and surface server errors via toasts/snackbars once backend endpoints are available.
5. **Share Links & Activity** – Extend placeholder pages (Shared, Recent, Trash) with backend-driven data and actions.
6. **CI & Tests** – Add component tests (React Testing Library) and end-to-end coverage (Playwright) once data flows are stable.

With the frontend scaffold in place, the next milestone is to stand up the Go backend APIs, after which we can point the React Query hooks at real endpoints and iterate on the upload/download workflows.
