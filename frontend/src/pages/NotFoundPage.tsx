import { Link } from 'react-router-dom'

const NotFoundPage = () => {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-background-light text-slate-800">
      <h1 className="text-4xl font-bold">404</h1>
      <p className="mt-2 text-slate-500">The page you requested does not exist.</p>
      <Link
        to="/"
        className="mt-6 inline-flex items-center justify-center rounded-lg bg-primary px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary/90"
      >
        Go home
      </Link>
    </div>
  )
}

export default NotFoundPage
