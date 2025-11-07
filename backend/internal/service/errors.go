package service

import "errors"

// Common errors used across services
var (
	// Auth errors
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrEmailAlreadyInUse   = errors.New("email already in use")

	// Credential errors
	ErrCredentialNotFound      = errors.New("credential not found")
	ErrCredentialAlreadyExists = errors.New("credential with this name already exists")
	ErrInvalidEncryptionKey    = errors.New("invalid encryption key")

	// Bucket errors
	ErrBucketNotFound      = errors.New("bucket not found")
	ErrBucketAlreadyExists = errors.New("bucket already exists")

	// Demo mode errors
	ErrDemoRestriction = errors.New("file preview and download are not available in demo mode")
)

// BucketProvisionError represents a failure when creating or ensuring the bucket
type BucketProvisionError struct {
	Reason string
	Err    error
}

func (e *BucketProvisionError) Error() string {
	if e == nil {
		return ""
	}
	if e.Reason != "" {
		return "bucket provisioning failed: " + e.Reason
	}
	return "bucket provisioning failed"
}

func (e *BucketProvisionError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func newBucketProvisionError(err error) error {
	if err == nil {
		return nil
	}
	return &BucketProvisionError{
		Reason: err.Error(),
		Err:    err,
	}
}

// CredentialDiscoveryError represents a failure when listing buckets for a credential
type CredentialDiscoveryError struct {
	Reason string
	Err    error
}

func (e *CredentialDiscoveryError) Error() string {
	if e == nil {
		return ""
	}
	if e.Reason != "" {
		return "bucket discovery failed: " + e.Reason
	}
	return "bucket discovery failed"
}

func (e *CredentialDiscoveryError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func newCredentialDiscoveryError(err error) error {
	if err == nil {
		return nil
	}
	return &CredentialDiscoveryError{
		Reason: err.Error(),
		Err:    err,
	}
}
