package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("record not found")

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(ctx context.Context, dsn string) (*Store, error) {
	if dsn == "" {
		return nil, errors.New("database dsn is required")
	}
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	return &Store{pool: pool}, nil
}

func (s *Store) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}

func (s *Store) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

func (s *Store) InitSchema(ctx context.Context) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			refresh_token_hash TEXT NOT NULL UNIQUE,
			expires_at TIMESTAMPTZ NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS sessions_user_id_idx ON sessions(user_id);`,
		`CREATE TABLE IF NOT EXISTS credentials (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			provider TEXT NOT NULL,
			region TEXT NOT NULL,
			endpoint TEXT NOT NULL,
			access_key TEXT NOT NULL,
			secret_key TEXT NOT NULL,
			encrypted_access_key TEXT NOT NULL DEFAULT '',
			encrypted_secret_key TEXT NOT NULL DEFAULT '',
			use_ssl BOOLEAN NOT NULL DEFAULT true,
			status TEXT NOT NULL,
			logo TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
		`ALTER TABLE credentials ADD COLUMN IF NOT EXISTS encrypted_access_key TEXT NOT NULL DEFAULT '';`,
		`ALTER TABLE credentials ADD COLUMN IF NOT EXISTS encrypted_secret_key TEXT NOT NULL DEFAULT '';`,
		`ALTER TABLE credentials ADD COLUMN IF NOT EXISTS use_ssl BOOLEAN NOT NULL DEFAULT true;`,
		`CREATE INDEX IF NOT EXISTS credentials_user_id_idx ON credentials(user_id);`,
		`CREATE TABLE IF NOT EXISTS buckets (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			credential_id UUID NOT NULL REFERENCES credentials(id) ON DELETE RESTRICT,
			name TEXT NOT NULL,
			region TEXT NOT NULL,
			description TEXT,
			size_bytes BIGINT NOT NULL DEFAULT 0,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE (user_id, name)
		);`,
		`ALTER TABLE buckets ADD COLUMN IF NOT EXISTS user_id UUID REFERENCES users(id) ON DELETE CASCADE;`,
		`ALTER TABLE buckets ADD COLUMN IF NOT EXISTS credential_id UUID REFERENCES credentials(id) ON DELETE RESTRICT;`,
		`ALTER TABLE buckets ADD COLUMN IF NOT EXISTS size_bytes BIGINT NOT NULL DEFAULT 0;`,
		`ALTER TABLE buckets ADD COLUMN IF NOT EXISTS region TEXT NOT NULL DEFAULT 'us-east-1';`,
		`ALTER TABLE buckets DROP CONSTRAINT IF EXISTS buckets_name_key;`,
		`CREATE UNIQUE INDEX IF NOT EXISTS buckets_user_id_name_idx ON buckets(user_id, name);`,
		`CREATE INDEX IF NOT EXISTS buckets_user_id_idx ON buckets(user_id);`,
		`CREATE INDEX IF NOT EXISTS buckets_credential_id_idx ON buckets(credential_id);`,
		`CREATE TABLE IF NOT EXISTS profiles (
			id UUID PRIMARY KEY,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT NOT NULL,
			language TEXT NOT NULL,
			timezone TEXT NOT NULL,
			avatar_url TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
	}

	for _, stmt := range stmts {
		if _, err := s.pool.Exec(ctx, stmt); err != nil {
			return err
		}
	}

	return nil
}

type BucketRecord struct {
	ID                 uuid.UUID
	UserID             uuid.UUID
	CredentialID       uuid.UUID
	Name               string
	Region             string
	Description        *string
	SizeBytes          int64
	CreatedAt          time.Time
	UpdatedAt          time.Time
	CredentialName     string
	CredentialProvider string
}

func (s *Store) ListBuckets(ctx context.Context, userID uuid.UUID) ([]BucketRecord, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT
			b.id,
			b.user_id,
			b.credential_id,
			b.name,
			b.region,
			b.description,
			b.size_bytes,
			b.created_at,
			b.updated_at,
			COALESCE(c.name, ''),
			COALESCE(c.provider, '')
		FROM buckets b
		JOIN credentials c ON c.id = b.credential_id
		WHERE b.user_id = $1
		ORDER BY b.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []BucketRecord
	for rows.Next() {
		var rec BucketRecord
		if err := rows.Scan(
			&rec.ID,
			&rec.UserID,
			&rec.CredentialID,
			&rec.Name,
			&rec.Region,
			&rec.Description,
			&rec.SizeBytes,
			&rec.CreatedAt,
			&rec.UpdatedAt,
			&rec.CredentialName,
			&rec.CredentialProvider,
		); err != nil {
			return nil, err
		}
		result = append(result, rec)
	}
	return result, rows.Err()
}

func (s *Store) GetBucket(ctx context.Context, id, userID uuid.UUID) (*BucketRecord, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT
			b.id,
			b.user_id,
			b.credential_id,
			b.name,
			b.region,
			b.description,
			b.size_bytes,
			b.created_at,
			b.updated_at,
			COALESCE(c.name, ''),
			COALESCE(c.provider, '')
		FROM buckets b
		JOIN credentials c ON c.id = b.credential_id
		WHERE b.id = $1 AND b.user_id = $2
	`, id, userID)
	var rec BucketRecord
	if err := row.Scan(
		&rec.ID,
		&rec.UserID,
		&rec.CredentialID,
		&rec.Name,
		&rec.Region,
		&rec.Description,
		&rec.SizeBytes,
		&rec.CreatedAt,
		&rec.UpdatedAt,
		&rec.CredentialName,
		&rec.CredentialProvider,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &rec, nil
}

func (s *Store) GetBucketByName(ctx context.Context, userID uuid.UUID, name string) (*BucketRecord, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT
			b.id,
			b.user_id,
			b.credential_id,
			b.name,
			b.region,
			b.description,
			b.size_bytes,
			b.created_at,
			b.updated_at,
			COALESCE(c.name, ''),
			COALESCE(c.provider, '')
		FROM buckets b
		JOIN credentials c ON c.id = b.credential_id
		WHERE b.user_id = $1 AND b.name = $2
	`, userID, name)
	var rec BucketRecord
	if err := row.Scan(
		&rec.ID,
		&rec.UserID,
		&rec.CredentialID,
		&rec.Name,
		&rec.Region,
		&rec.Description,
		&rec.SizeBytes,
		&rec.CreatedAt,
		&rec.UpdatedAt,
		&rec.CredentialName,
		&rec.CredentialProvider,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &rec, nil
}

func (s *Store) InsertBucket(ctx context.Context, userID, credentialID uuid.UUID, name, region string, description *string) (BucketRecord, error) {
	id := uuid.New()
	row := s.pool.QueryRow(ctx, `
		INSERT INTO buckets (id, user_id, credential_id, name, region, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, credential_id, name, region, description, size_bytes, created_at, updated_at
	`, id, userID, credentialID, name, region, description)
	var rec BucketRecord
	if err := row.Scan(
		&rec.ID,
		&rec.UserID,
		&rec.CredentialID,
		&rec.Name,
		&rec.Region,
		&rec.Description,
		&rec.SizeBytes,
		&rec.CreatedAt,
		&rec.UpdatedAt,
	); err != nil {
		return BucketRecord{}, err
	}
	return rec, nil
}

func (s *Store) UpdateBucketSize(ctx context.Context, id uuid.UUID, sizeBytes int64) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE buckets SET size_bytes = $2, updated_at = NOW()
		WHERE id = $1
	`, id, sizeBytes)
	return err
}

func (s *Store) UpdateBucket(ctx context.Context, id, userID uuid.UUID, description *string) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE buckets SET description = $3, updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`, id, userID, description)
	return err
}

func (s *Store) DeleteBucket(ctx context.Context, id, userID uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM buckets WHERE id = $1 AND user_id = $2`, id, userID)
	return err
}

type CredentialRecord struct {
	ID                 uuid.UUID
	UserID             uuid.UUID
	Name               string
	Provider           string
	Region             string
	Endpoint           string
	AccessKey          string
	SecretKey          string
	EncryptedAccessKey string
	EncryptedSecretKey string
	UseSSL             bool
	Status             string
	Logo               *string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (s *Store) ListCredentials(ctx context.Context, userID uuid.UUID) ([]CredentialRecord, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, user_id, name, provider, region, endpoint, access_key, secret_key, encrypted_access_key, encrypted_secret_key, use_ssl, status, logo, created_at, updated_at
		FROM credentials
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []CredentialRecord
	for rows.Next() {
		var rec CredentialRecord
		if err := rows.Scan(
			&rec.ID,
			&rec.UserID,
			&rec.Name,
			&rec.Provider,
			&rec.Region,
			&rec.Endpoint,
			&rec.AccessKey,
			&rec.SecretKey,
			&rec.EncryptedAccessKey,
			&rec.EncryptedSecretKey,
			&rec.UseSSL,
			&rec.Status,
			&rec.Logo,
			&rec.CreatedAt,
			&rec.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, rec)
	}
	return result, rows.Err()
}

type ProfileRecord struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Language  string
	Timezone  string
	AvatarURL *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRecord struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type SessionRecord struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	RefreshTokenHash string
	ExpiresAt        time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (s *Store) GetProfileByID(ctx context.Context, id uuid.UUID) (*ProfileRecord, error) {
	row := s.pool.QueryRow(ctx, `SELECT id, first_name, last_name, email, language, timezone, avatar_url, created_at, updated_at FROM profiles WHERE id = $1`, id)
	var rec ProfileRecord
	if err := row.Scan(&rec.ID, &rec.FirstName, &rec.LastName, &rec.Email, &rec.Language, &rec.Timezone, &rec.AvatarURL, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &rec, nil
}

func (s *Store) UpsertProfile(ctx context.Context, rec ProfileRecord) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO profiles (id, first_name, last_name, email, language, timezone, avatar_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			email = EXCLUDED.email,
			language = EXCLUDED.language,
			timezone = EXCLUDED.timezone,
			avatar_url = EXCLUDED.avatar_url,
			updated_at = NOW()
	`, rec.ID, rec.FirstName, rec.LastName, rec.Email, rec.Language, rec.Timezone, rec.AvatarURL)
	return err
}

func (s *Store) InsertUser(ctx context.Context, email, passwordHash, firstName, lastName string) (UserRecord, error) {
	id := uuid.New()
	row := s.pool.QueryRow(ctx, `
		INSERT INTO users (id, email, password_hash, first_name, last_name)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, email, password_hash, first_name, last_name, created_at, updated_at
	`, id, email, passwordHash, firstName, lastName)
	var rec UserRecord
	if err := row.Scan(&rec.ID, &rec.Email, &rec.PasswordHash, &rec.FirstName, &rec.LastName, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
		return UserRecord{}, err
	}
	return rec, nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*UserRecord, error) {
	row := s.pool.QueryRow(ctx, `SELECT id, email, password_hash, first_name, last_name, created_at, updated_at FROM users WHERE email = $1`, email)
	var rec UserRecord
	if err := row.Scan(&rec.ID, &rec.Email, &rec.PasswordHash, &rec.FirstName, &rec.LastName, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &rec, nil
}

func (s *Store) GetUserByID(ctx context.Context, id uuid.UUID) (*UserRecord, error) {
	row := s.pool.QueryRow(ctx, `SELECT id, email, password_hash, first_name, last_name, created_at, updated_at FROM users WHERE id = $1`, id)
	var rec UserRecord
	if err := row.Scan(&rec.ID, &rec.Email, &rec.PasswordHash, &rec.FirstName, &rec.LastName, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &rec, nil
}

func (s *Store) UpdateUser(ctx context.Context, id uuid.UUID, email, firstName, lastName string) error {
	tag, err := s.pool.Exec(ctx, `
		UPDATE users
		SET email = $2,
			first_name = $3,
			last_name = $4,
			updated_at = NOW()
		WHERE id = $1
	`, id, email, firstName, lastName)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) UpdateUserPassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	tag, err := s.pool.Exec(ctx, `
		UPDATE users
		SET password_hash = $2,
			updated_at = NOW()
		WHERE id = $1
	`, id, passwordHash)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) CreateSession(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) (SessionRecord, error) {
	id := uuid.New()
	row := s.pool.QueryRow(ctx, `
		INSERT INTO sessions (id, user_id, refresh_token_hash, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, refresh_token_hash, expires_at, created_at, updated_at
	`, id, userID, tokenHash, expiresAt)
	var rec SessionRecord
	if err := row.Scan(&rec.ID, &rec.UserID, &rec.RefreshTokenHash, &rec.ExpiresAt, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
		return SessionRecord{}, err
	}
	return rec, nil
}

func (s *Store) GetSessionByHash(ctx context.Context, hash string) (*SessionRecord, error) {
	row := s.pool.QueryRow(ctx, `SELECT id, user_id, refresh_token_hash, expires_at, created_at, updated_at FROM sessions WHERE refresh_token_hash = $1`, hash)
	var rec SessionRecord
	if err := row.Scan(&rec.ID, &rec.UserID, &rec.RefreshTokenHash, &rec.ExpiresAt, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &rec, nil
}

func (s *Store) UpdateSessionToken(ctx context.Context, sessionID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE sessions
		SET refresh_token_hash = $2,
			expires_at = $3,
			updated_at = NOW()
		WHERE id = $1
	`, sessionID, tokenHash, expiresAt)
	return err
}

func (s *Store) DeleteSessionByHash(ctx context.Context, hash string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM sessions WHERE refresh_token_hash = $1`, hash)
	return err
}

func (s *Store) DeleteSessionsForUser(ctx context.Context, userID uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM sessions WHERE user_id = $1`, userID)
	return err
}

func (s *Store) CreateCredential(ctx context.Context, userID uuid.UUID, name, provider, region, endpoint, encryptedAccessKey, encryptedSecretKey string, useSSL bool, status string, logo *string) (CredentialRecord, error) {
	id := uuid.New()
	row := s.pool.QueryRow(ctx, `
		INSERT INTO credentials (
			id, user_id, name, provider, region, endpoint,
			access_key, secret_key, encrypted_access_key, encrypted_secret_key,
			use_ssl, status, logo
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, user_id, name, provider, region, endpoint, access_key, secret_key, encrypted_access_key, encrypted_secret_key, use_ssl, status, logo, created_at, updated_at
	`, id, userID, name, provider, region, endpoint, encryptedAccessKey, encryptedSecretKey, encryptedAccessKey, encryptedSecretKey, useSSL, status, logo)
	var rec CredentialRecord
	if err := row.Scan(
		&rec.ID,
		&rec.UserID,
		&rec.Name,
		&rec.Provider,
		&rec.Region,
		&rec.Endpoint,
		&rec.AccessKey,
		&rec.SecretKey,
		&rec.EncryptedAccessKey,
		&rec.EncryptedSecretKey,
		&rec.UseSSL,
		&rec.Status,
		&rec.Logo,
		&rec.CreatedAt,
		&rec.UpdatedAt,
	); err != nil {
		return CredentialRecord{}, err
	}
	return rec, nil
}

func (s *Store) GetCredential(ctx context.Context, id, userID uuid.UUID) (*CredentialRecord, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT id, user_id, name, provider, region, endpoint, access_key, secret_key, encrypted_access_key, encrypted_secret_key, use_ssl, status, logo, created_at, updated_at
		FROM credentials
		WHERE id = $1 AND user_id = $2
	`, id, userID)
	var rec CredentialRecord
	if err := row.Scan(
		&rec.ID,
		&rec.UserID,
		&rec.Name,
		&rec.Provider,
		&rec.Region,
		&rec.Endpoint,
		&rec.AccessKey,
		&rec.SecretKey,
		&rec.EncryptedAccessKey,
		&rec.EncryptedSecretKey,
		&rec.UseSSL,
		&rec.Status,
		&rec.Logo,
		&rec.CreatedAt,
		&rec.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &rec, nil
}

func (s *Store) UpdateCredential(ctx context.Context, id, userID uuid.UUID, name, provider, region, endpoint, encryptedAccessKey, encryptedSecretKey string, useSSL bool, status string, logo *string) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE credentials
		SET name = $3,
			provider = $4,
			region = $5,
			endpoint = $6,
			access_key = $7,
			secret_key = $8,
			encrypted_access_key = $9,
			encrypted_secret_key = $10,
			use_ssl = $11,
			status = $12,
			logo = $13,
			updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`, id, userID, name, provider, region, endpoint, encryptedAccessKey, encryptedSecretKey, encryptedAccessKey, encryptedSecretKey, useSSL, status, logo)
	return err
}

func (s *Store) DeleteCredential(ctx context.Context, id, userID uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM credentials WHERE id = $1 AND user_id = $2`, id, userID)
	return err
}
