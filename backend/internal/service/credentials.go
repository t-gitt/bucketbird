package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"bucketbird/backend/internal/repository"
	"bucketbird/backend/internal/storage"
	"bucketbird/backend/pkg/crypto"

	"github.com/google/uuid"
)

type CredentialService struct {
	credentials   repository.CredentialRepository
	encryptionKey []byte
	logger        *slog.Logger
}

func NewCredentialService(
	credentials repository.CredentialRepository,
	encryptionKey []byte,
	logger *slog.Logger,
) *CredentialService {
	return &CredentialService{
		credentials:   credentials,
		encryptionKey: encryptionKey,
		logger:        logger,
	}
}

type CreateCredentialInput struct {
	UserID    uuid.UUID
	Name      string
	Provider  string
	Region    string
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Logo      *string
}

func (s *CredentialService) Create(ctx context.Context, input CreateCredentialInput) (*repository.Credential, error) {
	// Encrypt credentials
	encryptedAccessKey, err := crypto.EncryptAES(input.AccessKey, s.encryptionKey)
	if err != nil {
		return nil, err
	}

	encryptedSecretKey, err := crypto.EncryptAES(input.SecretKey, s.encryptionKey)
	if err != nil {
		return nil, err
	}

	// Test connection before saving
	if err := s.testConnection(ctx, input.Endpoint, input.Region, input.AccessKey, input.SecretKey, input.UseSSL); err != nil {
		s.logger.Warn("failed to connect to S3", slog.Any("error", err))
		// Don't fail here, just log - user might be adding credentials for later use
	}

	// Create credential
	cred := &repository.Credential{
		UserID:             input.UserID,
		Name:               input.Name,
		Provider:           input.Provider,
		Region:             input.Region,
		Endpoint:           input.Endpoint,
		EncryptedAccessKey: encryptedAccessKey,
		EncryptedSecretKey: encryptedSecretKey,
		UseSSL:             input.UseSSL,
		Status:             "active",
		Logo:               input.Logo,
	}

	return s.credentials.Create(ctx, cred)
}

func (s *CredentialService) List(ctx context.Context, userID uuid.UUID) ([]*repository.Credential, error) {
	return s.credentials.List(ctx, userID)
}

func (s *CredentialService) Get(ctx context.Context, id, userID uuid.UUID) (*repository.Credential, error) {
	cred, err := s.credentials.Get(ctx, id, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrCredentialNotFound
		}
		return nil, err
	}
	return cred, nil
}

type UpdateCredentialInput struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Provider  string
	Region    string
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Logo      *string
}

func (s *CredentialService) Update(ctx context.Context, input UpdateCredentialInput) error {
	// Verify credential exists
	existing, err := s.credentials.Get(ctx, input.ID, input.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrCredentialNotFound
		}
		return err
	}

	// Encrypt new credentials
	encryptedAccessKey, err := crypto.EncryptAES(input.AccessKey, s.encryptionKey)
	if err != nil {
		return err
	}

	encryptedSecretKey, err := crypto.EncryptAES(input.SecretKey, s.encryptionKey)
	if err != nil {
		return err
	}

	// Test connection
	if err := s.testConnection(ctx, input.Endpoint, input.Region, input.AccessKey, input.SecretKey, input.UseSSL); err != nil {
		s.logger.Warn("failed to connect to S3", slog.Any("error", err))
	}

	// Update credential
	existing.Name = input.Name
	existing.Provider = input.Provider
	existing.Region = input.Region
	existing.Endpoint = input.Endpoint
	existing.EncryptedAccessKey = encryptedAccessKey
	existing.EncryptedSecretKey = encryptedSecretKey
	existing.UseSSL = input.UseSSL
	existing.Logo = input.Logo

	return s.credentials.Update(ctx, existing)
}

func (s *CredentialService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	// Verify credential exists
	if _, err := s.credentials.Get(ctx, id, userID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrCredentialNotFound
		}
		return err
	}

	// Delete will cascade to buckets automatically via database constraint
	return s.credentials.Delete(ctx, id, userID)
}

// GetDecryptedCredentials returns decrypted access and secret keys
func (s *CredentialService) GetDecryptedCredentials(ctx context.Context, id, userID uuid.UUID) (accessKey, secretKey string, err error) {
	cred, err := s.credentials.Get(ctx, id, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", "", ErrCredentialNotFound
		}
		return "", "", err
	}

	accessKey, err = crypto.DecryptAES(cred.EncryptedAccessKey, s.encryptionKey)
	if err != nil {
		return "", "", err
	}

	secretKey, err = crypto.DecryptAES(cred.EncryptedSecretKey, s.encryptionKey)
	if err != nil {
		return "", "", err
	}

	return accessKey, secretKey, nil
}

type TestCredentialResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type DiscoveredBucket struct {
	Name      string
	CreatedAt *time.Time
}

func (s *CredentialService) Test(ctx context.Context, id, userID uuid.UUID) (*TestCredentialResult, error) {
	// Get credential
	cred, err := s.credentials.Get(ctx, id, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return &TestCredentialResult{
				Success: false,
				Message: "Credential not found",
			}, nil
		}
		return &TestCredentialResult{
			Success: false,
			Message: "Failed to retrieve credential",
		}, nil
	}

	// Decrypt credentials
	accessKey, err := crypto.DecryptAES(cred.EncryptedAccessKey, s.encryptionKey)
	if err != nil {
		return &TestCredentialResult{
			Success: false,
			Message: "Failed to decrypt access key",
		}, nil
	}

	secretKey, err := crypto.DecryptAES(cred.EncryptedSecretKey, s.encryptionKey)
	if err != nil {
		return &TestCredentialResult{
			Success: false,
			Message: "Failed to decrypt secret key",
		}, nil
	}

	// Test connection
	if err := s.testConnection(ctx, cred.Endpoint, cred.Region, accessKey, secretKey, cred.UseSSL); err != nil {
		return &TestCredentialResult{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &TestCredentialResult{
		Success: true,
		Message: "Connection successful",
	}, nil
}

func (s *CredentialService) testConnection(ctx context.Context, endpoint, region, accessKey, secretKey string, useSSL bool) error {
	store, err := storage.NewObjectStoreWithCredentials(ctx, endpoint, region, accessKey, secretKey, useSSL)
	if err != nil {
		return err
	}
	return store.TestConnection(ctx)
}

func (s *CredentialService) DiscoverBuckets(ctx context.Context, id, userID uuid.UUID) ([]DiscoveredBucket, error) {
	cred, err := s.credentials.Get(ctx, id, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrCredentialNotFound
		}
		return nil, err
	}

	accessKey, err := crypto.DecryptAES(cred.EncryptedAccessKey, s.encryptionKey)
	if err != nil {
		return nil, err
	}

	secretKey, err := crypto.DecryptAES(cred.EncryptedSecretKey, s.encryptionKey)
	if err != nil {
		return nil, err
	}

	store, err := storage.NewObjectStoreWithCredentials(
		ctx,
		cred.Endpoint,
		cred.Region,
		accessKey,
		secretKey,
		cred.UseSSL,
	)
	if err != nil {
		return nil, newCredentialDiscoveryError(err)
	}

	buckets, err := store.ListBuckets(ctx)
	if err != nil {
		return nil, newCredentialDiscoveryError(err)
	}

	discovered := make([]DiscoveredBucket, 0, len(buckets))
	for _, bucket := range buckets {
		if bucket.Name == nil || *bucket.Name == "" {
			continue
		}

		var createdAt *time.Time
		if bucket.CreationDate != nil {
			t := bucket.CreationDate.UTC()
			createdAt = &t
		}

		discovered = append(discovered, DiscoveredBucket{
			Name:      *bucket.Name,
			CreatedAt: createdAt,
		})
	}

	return discovered, nil
}

// Helper function used by BucketService
func decryptCredential(encrypted string, key []byte) (string, error) {
	return crypto.DecryptAES(encrypted, key)
}
