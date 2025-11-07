# BucketBird Backend

BucketBird is a multi-cloud S3-compatible object storage management platform. This backend provides a RESTful API for managing S3 buckets, credentials, and objects across multiple S3-compatible storage providers.

## Architecture

The backend follows a clean architecture pattern with clear separation of concerns:

```
backend/
├── cmd/                      # Application entry points
│   ├── bucketbird-api/      # Main API server
│   ├── bucketbird-create-user/
│   └── bucketbird-password-reset/
├── internal/                 # Private application code
│   ├── api/                 # HTTP handlers (presentation layer)
│   │   ├── auth/           # Authentication endpoints
│   │   ├── buckets/        # Bucket management endpoints
│   │   ├── credentials/    # Credential management endpoints
│   │   └── profile/        # User profile endpoints
│   ├── service/            # Business logic layer
│   ├── repository/         # Data access layer
│   ├── middleware/         # HTTP middleware
│   ├── storage/            # S3 client abstraction
│   ├── config/             # Configuration management
│   └── logging/            # Logging setup
├── pkg/                     # Public reusable packages
│   ├── crypto/             # Encryption utilities
│   ├── jwt/                # JWT token management
│   └── validator/          # Input validation
├── migrations/              # Database migrations
└── queries/                 # SQL queries for sqlc

```

### Layer Responsibilities

- **API Layer** (`internal/api`): HTTP request/response handling, input validation, authentication
- **Service Layer** (`internal/service`): Business logic, orchestration, transaction management
- **Repository Layer** (`internal/repository`): Database operations using sqlc
- **Storage Layer** (`internal/storage`): S3-compatible storage operations using AWS SDK

## Tech Stack

- **Language**: Go 1.23+
- **Web Framework**: Chi router
- **Database**: PostgreSQL with pgx driver
- **Code Generation**: sqlc for type-safe SQL queries
- **Authentication**: JWT tokens with bcrypt password hashing
- **Storage**: AWS SDK v2 for S3-compatible storage
- **Encryption**: AES-256-GCM for credential encryption

## Features

### Authentication & Authorization
- JWT-based authentication with access and refresh tokens
- Secure password hashing with bcrypt
- Session management with refresh token rotation
- User registration and login

### Credential Management
- Encrypted storage of S3 credentials (access key, secret key)
- Support for multiple S3-compatible providers (AWS S3, MinIO, Wasabi, etc.)
- Connection testing before saving credentials
- AES-256-GCM encryption for sensitive data

### Bucket Management
- List, create, and delete S3 buckets
- Bucket size tracking and formatting
- Multi-credential support for different providers
- Metadata storage in PostgreSQL

### Object Operations
- List objects with folder navigation
- Upload files with progress tracking
- Download files and folders (as zip)
- Recursive search across all objects
- Folder creation and management
- Rename objects and folders (recursive)
- Delete objects and folders (recursive)
- Object metadata viewing
- Preview support for various file types

## Configuration

The application is configured via environment variables:

```bash
# Server
PORT=8080
ENV=development

# Database
DB_DSN=postgres://user:password@localhost:5432/bucketbird?sslmode=disable

# Security
JWT_SECRET=your-jwt-secret-key-min-32-chars
ENCRYPTION_KEY=your-encryption-key-must-be-32-bytes
ACCESS_TOKEN_TTL=15m
REFRESH_TOKEN_TTL=7d

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:5173
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization
CORS_MAX_AGE=300

# Features
REGISTRATION_ENABLED=true
```

## Database Setup

### Using Docker Compose (Recommended)

```bash
# Start PostgreSQL
docker compose up -d postgres

# The database will be automatically initialized
```

### Manual Setup

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path migrations -database "postgres://bucketbird:bucketbird@localhost:5432/bucketbird?sslmode=disable" up
```

### Using sqlc

The project uses sqlc for type-safe SQL queries. To regenerate code after modifying queries:

```bash
# Install sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Generate code
sqlc generate
```

Query files are in `queries/` and generated code goes to `internal/repository/sqlc/`.

## Building

```bash
# Build the API server
go build -o bucketbird-api ./cmd/bucketbird-api

# Build with optimizations for production
CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bucketbird-api ./cmd/bucketbird-api
```

## Running

### Development

```bash
# Set environment variables
export DB_DSN="postgres://bucketbird:bucketbird@localhost:5432/bucketbird?sslmode=disable"
export JWT_SECRET="your-secret-key-at-least-32-characters-long"
export ENCRYPTION_KEY="your-32-byte-encryption-key-here!!"
export CORS_ALLOWED_ORIGINS="http://localhost:5173"

# Run the server
go run ./cmd/bucketbird-api
```

### Production

```bash
# Build the binary
CGO_ENABLED=0 GOOS=linux go build -o bucketbird-api ./cmd/bucketbird-api

# Run with environment variables
./bucketbird-api
```

### Docker

```bash
# Build image
docker build -t bucketbird-api .

# Run container
docker run -p 8080:8080 \
  -e DB_DSN="postgres://..." \
  -e JWT_SECRET="..." \
  -e ENCRYPTION_KEY="..." \
  bucketbird-api
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login and get tokens
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - Logout and invalidate session
- `GET /api/v1/auth/me` - Get current user

### Credentials
- `GET /api/v1/credentials` - List all credentials
- `POST /api/v1/credentials` - Create new credential
- `GET /api/v1/credentials/:id` - Get credential details
- `PUT /api/v1/credentials/:id` - Update credential
- `DELETE /api/v1/credentials/:id` - Delete credential
- `POST /api/v1/credentials/:id/test` - Test credential connection

### Buckets
- `GET /api/v1/buckets` - List all buckets
- `POST /api/v1/buckets` - Create new bucket
- `GET /api/v1/buckets/:id` - Get bucket details
- `PUT /api/v1/buckets/:id` - Update bucket
- `DELETE /api/v1/buckets/:id` - Delete bucket

### Objects
- `GET /api/v1/buckets/:id/objects` - List objects (with prefix support)
- `GET /api/v1/buckets/:id/objects/search` - Search objects
- `POST /api/v1/buckets/:id/objects/upload` - Upload file
- `GET /api/v1/buckets/:id/objects/download` - Download file or folder
- `POST /api/v1/buckets/:id/objects/folders` - Create folder
- `POST /api/v1/buckets/:id/objects/delete` - Delete objects/folders
- `POST /api/v1/buckets/:id/objects/rename` - Rename object/folder
- `POST /api/v1/buckets/:id/objects/copy` - Copy object
- `GET /api/v1/buckets/:id/objects/metadata` - Get object metadata
- `POST /api/v1/buckets/:id/objects/presign` - Generate presigned URL

### Profile
- `GET /api/v1/profile` - Get user profile
- `PUT /api/v1/profile` - Update profile
- `PUT /api/v1/profile/password` - Change password

## Security

### Authentication
- JWT tokens with configurable expiration
- Refresh token rotation for enhanced security
- Secure password hashing with bcrypt (cost factor 10)
- Session management with database-backed refresh tokens

### Encryption
- AES-256-GCM encryption for S3 credentials
- Unique nonce per encryption operation
- Environment-based encryption key management

### HTTP Security Headers
- CORS configuration
- Security headers middleware
- Request ID tracking
- Panic recovery

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Code Quality

```bash
# Format code
go fmt ./...

# Lint code (requires golangci-lint)
golangci-lint run

# Vet code
go vet ./...
```

### Database Migrations

```bash
# Create a new migration
migrate create -ext sql -dir migrations -seq migration_name

# Check migration version
migrate -path migrations -database "$DB_DSN" version

# Force a specific version (use with caution)
migrate -path migrations -database "$DB_DSN" force VERSION
```

## Troubleshooting

### Database Connection Issues

```bash
# Test database connection
psql "$DB_DSN"

# Check if migrations are up to date
migrate -path migrations -database "$DB_DSN" version
```

### S3 Connection Issues

- Verify credentials are correct
- Check endpoint URL format (should include protocol: http:// or https://)
- Ensure SSL setting matches the endpoint
- Test credentials using the `/credentials/:id/test` endpoint

### Token Issues

- Ensure `JWT_SECRET` is at least 32 characters
- Check token expiration times are reasonable
- Verify system clock is synchronized

## Contributing

1. Follow the existing code structure and patterns
2. Use sqlc for all database queries
3. Write tests for new functionality
4. Follow Go best practices and idioms
5. Update documentation as needed

## License

See LICENSE file in the repository root.
