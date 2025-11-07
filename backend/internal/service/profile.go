package service

import (
	"context"

	"bucketbird/backend/internal/repository"
	"bucketbird/backend/pkg/crypto"

	"github.com/google/uuid"
)

type ProfileService struct {
	users repository.UserRepository
}

func NewProfileService(users repository.UserRepository) *ProfileService {
	return &ProfileService{
		users: users,
	}
}

type ProfileData struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Email     string
}

type UpdateProfileInput struct {
	FirstName string
	LastName  string
	Email     string
}

func (s *ProfileService) Get(ctx context.Context, userID uuid.UUID) (*ProfileData, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &ProfileData{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}, nil
}

func (s *ProfileService) Update(ctx context.Context, userID uuid.UUID, input UpdateProfileInput) (*ProfileData, error) {
	// Check if email is already in use by another user
	existingUser, err := s.users.GetByEmail(ctx, input.Email)
	if err == nil && existingUser.ID != userID {
		return nil, ErrEmailAlreadyInUse
	}

	if err := s.users.Update(ctx, userID, input.Email, input.FirstName, input.LastName); err != nil {
		return nil, err
	}

	return s.Get(ctx, userID)
}

func (s *ProfileService) UpdatePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify current password
	valid, err := crypto.VerifyPassword(user.PasswordHash, currentPassword)
	if err != nil {
		return err
	}
	if !valid {
		return ErrInvalidCredentials
	}

	// Hash new password
	newHash, err := crypto.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	return s.users.UpdatePassword(ctx, userID, newHash)
}
