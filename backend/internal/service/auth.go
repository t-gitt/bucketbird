package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"bucketbird/backend/internal/repository"
	"bucketbird/backend/pkg/crypto"
	"bucketbird/backend/pkg/jwt"

	"github.com/google/uuid"
)

type AuthService struct {
	users           repository.UserRepository
	sessions        repository.SessionRepository
	tokenManager    *jwt.TokenManager
	refreshTokenTTL time.Duration
	logger          *slog.Logger
}

func NewAuthService(
	users repository.UserRepository,
	sessions repository.SessionRepository,
	tokenManager *jwt.TokenManager,
	refreshTokenTTL time.Duration,
	logger *slog.Logger,
) *AuthService {
	return &AuthService{
		users:           users,
		sessions:        sessions,
		tokenManager:    tokenManager,
		refreshTokenTTL: refreshTokenTTL,
		logger:          logger,
	}
}

type RegisterInput struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

type AuthResult struct {
	User          *repository.User
	AccessToken   string
	AccessExpiry  time.Time
	RefreshToken  string
	RefreshExpiry time.Time
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*AuthResult, error) {
	// Normalize and validate email
	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" {
		return nil, errors.New("email is required")
	}

	// Validate password
	if input.Password == "" {
		return nil, errors.New("password is required")
	}

	// Check if email already exists
	if _, err := s.users.GetByEmail(ctx, email); err == nil {
		return nil, ErrEmailAlreadyInUse
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	// Hash password
	hash, err := crypto.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	firstName := strings.TrimSpace(input.FirstName)
	lastName := strings.TrimSpace(input.LastName)
	user, err := s.users.Create(ctx, email, hash, firstName, lastName)
	if err != nil {
		return nil, err
	}

	// Issue tokens
	tokens, err := s.issueTokens(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:          user,
		AccessToken:   tokens.accessToken,
		AccessExpiry:  tokens.accessExpiry,
		RefreshToken:  tokens.refreshToken,
		RefreshExpiry: tokens.refreshExpiry,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	// Normalize email
	normalized := strings.TrimSpace(strings.ToLower(email))

	// Get user by email
	user, err := s.users.GetByEmail(ctx, normalized)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Verify password
	ok, err := crypto.VerifyPassword(user.PasswordHash, password)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrInvalidCredentials
	}

	// Issue tokens
	tokens, err := s.issueTokens(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:          user,
		AccessToken:   tokens.accessToken,
		AccessExpiry:  tokens.accessExpiry,
		RefreshToken:  tokens.refreshToken,
		RefreshExpiry: tokens.refreshExpiry,
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*AuthResult, error) {
	if refreshToken == "" {
		return nil, ErrInvalidRefreshToken
	}

	// Hash and lookup session
	hash := crypto.HashRefreshToken(refreshToken)
	session, err := s.sessions.GetByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, err
	}

	// Check expiry
	if time.Now().After(session.ExpiresAt) {
		_ = s.sessions.DeleteByHash(ctx, hash)
		return nil, ErrInvalidRefreshToken
	}

	// Get user
	user, err := s.users.GetByID(ctx, session.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, err
	}

	// Rotate session
	tokens, err := s.rotateSession(ctx, session.ID, user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:          user,
		AccessToken:   tokens.accessToken,
		AccessExpiry:  tokens.accessExpiry,
		RefreshToken:  tokens.refreshToken,
		RefreshExpiry: tokens.refreshExpiry,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	hash := crypto.HashRefreshToken(refreshToken)
	return s.sessions.DeleteByHash(ctx, hash)
}

func (s *AuthService) ValidateAccessToken(ctx context.Context, token string) (*repository.User, error) {
	// Validate token
	userID, err := s.tokenManager.Validate(token)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Get user
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	return user, nil
}

// DemoLogin authenticates as the demo user without password
func (s *AuthService) DemoLogin(ctx context.Context) (*AuthResult, error) {
	const demoEmail = "demo@bucketbird.app"

	// Get demo user by email
	user, err := s.users.GetByEmail(ctx, demoEmail)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, errors.New("demo user not found")
		}
		return nil, err
	}

	// Verify it's actually a demo user
	if !user.IsDemo {
		return nil, errors.New("user is not a demo account")
	}

	// Issue tokens
	tokens, err := s.issueTokens(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:          user,
		AccessToken:   tokens.accessToken,
		AccessExpiry:  tokens.accessExpiry,
		RefreshToken:  tokens.refreshToken,
		RefreshExpiry: tokens.refreshExpiry,
	}, nil
}

type tokens struct {
	accessToken   string
	accessExpiry  time.Time
	refreshToken  string
	refreshExpiry time.Time
}

func (s *AuthService) issueTokens(ctx context.Context, userID uuid.UUID) (*tokens, error) {
	// Generate access token
	accessToken, accessExpiry, err := s.tokenManager.Generate(userID)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := crypto.GenerateRandomToken(32)
	if err != nil {
		return nil, err
	}
	refreshExpiry := time.Now().Add(s.refreshTokenTTL)
	hash := crypto.HashRefreshToken(refreshToken)

	// Create session
	if _, err := s.sessions.Create(ctx, userID, hash, refreshExpiry); err != nil {
		return nil, err
	}

	return &tokens{
		accessToken:   accessToken,
		accessExpiry:  accessExpiry,
		refreshToken:  refreshToken,
		refreshExpiry: refreshExpiry,
	}, nil
}

func (s *AuthService) rotateSession(ctx context.Context, sessionID, userID uuid.UUID) (*tokens, error) {
	// Generate new access token
	accessToken, accessExpiry, err := s.tokenManager.Generate(userID)
	if err != nil {
		return nil, err
	}

	// Generate new refresh token
	refreshToken, err := crypto.GenerateRandomToken(32)
	if err != nil {
		return nil, err
	}
	refreshExpiry := time.Now().Add(s.refreshTokenTTL)
	hash := crypto.HashRefreshToken(refreshToken)

	// Update session
	if err := s.sessions.UpdateToken(ctx, sessionID, hash, refreshExpiry); err != nil {
		return nil, err
	}

	return &tokens{
		accessToken:   accessToken,
		accessExpiry:  accessExpiry,
		refreshToken:  refreshToken,
		refreshExpiry: refreshExpiry,
	}, nil
}
