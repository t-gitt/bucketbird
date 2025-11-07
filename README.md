# BucketBird

Full-stack S3 storage management application. The project includes a Go API and a React + Tailwind frontend with Docker Compose orchestration for easy deployment.

## Services

- **backend**: Go 1.23 HTTP API backed by PostgreSQL for metadata and MinIO (S3 compatible) for object storage. Exposes bucket CRUD, presigned URL generation, credentials, and profile endpoints.
- **frontend**: React 19 single-page application built with Vite, Tailwind CSS, and TanStack Query. The UI mirrors the supplied login, dashboard, bucket content, and settings screens and consumes the backend API.

## Prerequisites

- Go 1.23+ (for local backend development)
- Node.js 20+ (for local frontend development)
- Docker & Docker Compose (for full-stack runtime)

## Running with Docker Compose

```bash
docker compose up --build
```

The compose file now boots four containers:

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080 (health check at `/health`)
- PostgreSQL: exposed on `localhost:5432` (`bucketbird/bucketbird`)
- MinIO: S3 API on `localhost:9000`, console on http://localhost:9001 (`minioadmin/minioadmin`)

The frontend ships with `public/config.js`, which defaults the API URL to `http://localhost:8080`. Adjust that file (or rebuild with a different `VITE_API_URL` build argument) if you deploy to another host/port.

## Local Development

### Backend

```bash
cd backend
go build ./...
go run ./cmd/bucketbird-api
```

Configuration is provided through environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `BB_APP_NAME` | `bucketbird-api` | Service name used in logs |
| `BB_ENV` | `development` | Controls log verbosity |
| `BB_HTTP_PORT` | `8080` | Port the API listens on |
| `BB_ALLOWED_ORIGINS` | `*` | Comma-separated list of CORS origins |
| `BB_HTTP_READ_TIMEOUT` | `30m` | HTTP server read timeout |
| `BB_HTTP_WRITE_TIMEOUT` | `30m` | HTTP server write timeout |
| `BB_DB_DSN` | See below | PostgreSQL connection string (`postgres://user:pass@host:port/db?sslmode=disable`) |
| `BB_DB_HOST` | `postgres` | Database host |
| `BB_DB_PORT` | `5432` | Database port |
| `BB_DB_NAME` | `bucketbird` | Database name |
| `BB_DB_USER` | `bucketbird` | Database user |
| `BB_DB_PASSWORD` | `bucketbird` | Database password |
| `BB_JWT_SECRET` | _see config_ | JWT signing secret |
| `BB_ENCRYPTION_KEY` | _see config_ | Encryption key (must be 32 bytes) |
| `BB_ACCESS_TOKEN_TTL` | `15m` | Access token lifetime |
| `BB_REFRESH_TOKEN_TTL` | `7d` | Refresh token lifetime |
| `BB_ALLOW_REGISTRATION` | `true` | Set to `false` to disable self-service signup; accounts must be created by an administrator |

Key endpoints:

- `GET /health`
- `GET /api/v1/buckets`
- `POST /api/v1/buckets` – create a new S3 bucket + metadata record
- `DELETE /api/v1/buckets/:bucketID` – remove the bucket and metadata
- `GET /api/v1/buckets/:bucketID/objects`
- `POST /api/v1/buckets/:bucketID/objects/presign` – generate upload/download URLs
- `GET /api/v1/credentials`
- `GET /api/v1/profile`

### Frontend

```bash
cd frontend
cp .env.example .env  # adjust VITE_API_URL if needed
npm install
npm run dev
```

The Vite dev server runs on http://localhost:5173 by default. The app automatically points API calls to `VITE_API_URL` or to the runtime value from `public/config.js` when built.

### Testing & Builds

- Frontend lint/build: `npm run lint`, `npm run build`
- Backend build: `go build ./...`

## Admin Tasks

### Creating Users

If self-service registration is disabled (`BB_ALLOW_REGISTRATION=false`), administrators can provision accounts from the CLI:

```bash
go run ./backend/cmd/bucketbird-create-user \
  --email alex@example.com \
  --password 'ChangeMe!123' \
  --first-name Alex \
  --last-name Doe
```

The command connects to the configured database, creates the user, and sets up the matching profile record.

### Resetting User Passwords

BucketBird includes a CLI tool for administrators to reset user passwords. When users forget their password, they should contact an admin who can reset it using this tool.

#### Option 1: Using Docker (Recommended)

If running BucketBird with Docker Compose:

```bash
# Build the password reset tool inside the backend container
docker-compose exec backend go build -o /tmp/bucketbird-password-reset ./cmd/bucketbird-password-reset

# Reset a user's password
docker-compose exec backend /tmp/bucketbird-password-reset --email user@example.com --password "newSecurePass123"
```

#### Option 2: Local Build

If running the backend locally:

```bash
# Build the password reset tool
cd backend
go build -o bucketbird-password-reset ./cmd/bucketbird-password-reset

# Reset a user's password
./bucketbird-password-reset --email user@example.com --password "newSecurePass123"
```

#### Requirements

- Password must be at least 8 characters long
- The tool uses the same environment variables as the main API (database connection settings)
- The user must exist in the database (identified by email address)

#### Success Output

```
✓ Password successfully reset for user: user@example.com (John Doe)
```

For more detailed documentation, see: `backend/cmd/bucketbird-password-reset/README.md`

## Next Steps

- Extend the placeholder pages (Shared, Recent, Trash) with backend-backed data and sharing semantics
- Introduce automated testing (Go unit tests, React component tests, Playwright E2E) and CI pipelines
- Add email-based password reset flow for self-service recovery

## Repository Layout

```
backend/               # Go API (cmd, internal packages, Dockerfile)
  ├── cmd/
  │   ├── bucketbird-api/           # Main API server
  │   └── bucketbird-password-reset/ # Admin password reset CLI tool
  └── internal/        # Shared packages (config, data, http, security, storage)
frontend/              # React SPA (Vite project, Dockerfile)
docker-compose.yml     # Full-stack runtime definition
README.md              # This file
```
