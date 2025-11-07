package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"bucketbird/backend/internal/middleware"
	"bucketbird/backend/internal/service"
)

const refreshTokenCookieName = "bb_refresh_token"

type Handler struct {
	authService     *service.AuthService
	logger          *slog.Logger
	cookieSecure    bool
	enableDemoLogin bool
}

func NewHandler(authService *service.AuthService, logger *slog.Logger, cookieSecure bool, enableDemoLogin bool) *Handler {
	return &Handler{
		authService:     authService,
		logger:          logger,
		cookieSecure:    cookieSecure,
		enableDemoLogin: enableDemoLogin,
	}
}

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type AuthResponse struct {
	User UserDTO       `json:"user"`
	Auth AuthTokensDTO `json:"auth"`
}

type AuthTokensDTO struct {
	AccessToken   string `json:"accessToken"`
	AccessExpiry  int64  `json:"accessExpiry"`
	RefreshExpiry int64  `json:"refreshExpiry"`
}

type UserDTO struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	IsReadonly bool   `json:"isReadonly"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.authService.Register(r.Context(), service.RegisterInput{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyInUse) {
			h.respondError(w, "Email already in use", http.StatusConflict)
			return
		}
		h.logger.Error("register failed", slog.Any("error", err))
		h.respondError(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	h.setRefreshTokenCookie(w, result.RefreshToken, result.RefreshExpiry)
	h.respondJSON(w, AuthResponse{
		User: UserDTO{
			ID:         result.User.ID.String(),
			Email:      result.User.Email,
			FirstName:  result.User.FirstName,
			LastName:   result.User.LastName,
			IsReadonly: result.User.IsDemo,
		},
		Auth: AuthTokensDTO{
			AccessToken:   result.AccessToken,
			AccessExpiry:  result.AccessExpiry.Unix(),
			RefreshExpiry: result.RefreshExpiry.Unix(),
		},
	}, http.StatusCreated)
}

type LoginRequest struct{
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			h.respondError(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		h.logger.Error("login failed", slog.Any("error", err))
		h.respondError(w, "Login failed", http.StatusInternalServerError)
		return
	}

	h.setRefreshTokenCookie(w, result.RefreshToken, result.RefreshExpiry)
	h.respondJSON(w, AuthResponse{
		User: UserDTO{
			ID:         result.User.ID.String(),
			Email:      result.User.Email,
			FirstName:  result.User.FirstName,
			LastName:   result.User.LastName,
			IsReadonly: result.User.IsDemo,
		},
		Auth: AuthTokensDTO{
			AccessToken:   result.AccessToken,
			AccessExpiry:  result.AccessExpiry.Unix(),
			RefreshExpiry: result.RefreshExpiry.Unix(),
		},
	}, http.StatusOK)
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken,omitempty"`
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	refreshToken := h.refreshTokenFromRequest(r, req.RefreshToken)
	if refreshToken == "" {
		h.respondError(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	result, err := h.authService.Refresh(r.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRefreshToken) {
			h.respondError(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}
		h.logger.Error("refresh failed", slog.Any("error", err))
		h.respondError(w, "Refresh failed", http.StatusInternalServerError)
		return
	}

	h.setRefreshTokenCookie(w, result.RefreshToken, result.RefreshExpiry)
	h.respondJSON(w, AuthResponse{
		User: UserDTO{
			ID:         result.User.ID.String(),
			Email:      result.User.Email,
			FirstName:  result.User.FirstName,
			LastName:   result.User.LastName,
			IsReadonly: result.User.IsDemo,
		},
		Auth: AuthTokensDTO{
			AccessToken:   result.AccessToken,
			AccessExpiry:  result.AccessExpiry.Unix(),
			RefreshExpiry: result.RefreshExpiry.Unix(),
		},
	}, http.StatusOK)
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var req LogoutRequest
	// Logout can work without a body (just invalidates current session)
	_ = json.NewDecoder(r.Body).Decode(&req)

	refreshToken := h.refreshTokenFromRequest(r, req.RefreshToken)
	if refreshToken != "" {
		if err := h.authService.Logout(r.Context(), refreshToken); err != nil {
			h.logger.Warn("logout failed", slog.Any("error", err))
		}
	}
	h.clearRefreshTokenCookie(w)

	h.respondJSON(w, map[string]interface{}{
		"message": "logged out successfully",
	}, http.StatusOK)
}

func (h *Handler) DemoLogin(w http.ResponseWriter, r *http.Request) {
	if !h.enableDemoLogin {
		h.respondError(w, "Demo login is disabled", http.StatusNotFound)
		return
	}

	result, err := h.authService.DemoLogin(r.Context())
	if err != nil {
		h.logger.Error("demo login failed", slog.Any("error", err))
		h.respondError(w, "Demo login failed", http.StatusInternalServerError)
		return
	}

	h.setRefreshTokenCookie(w, result.RefreshToken, result.RefreshExpiry)
	h.respondJSON(w, AuthResponse{
		User: UserDTO{
			ID:         result.User.ID.String(),
			Email:      result.User.Email,
			FirstName:  result.User.FirstName,
			LastName:   result.User.LastName,
			IsReadonly: result.User.IsDemo,
		},
		Auth: AuthTokensDTO{
			AccessToken:   result.AccessToken,
			AccessExpiry:  result.AccessExpiry.Unix(),
			RefreshExpiry: result.RefreshExpiry.Unix(),
		},
	}, http.StatusOK)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		h.respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	h.respondJSON(w, map[string]interface{}{"user": UserDTO{
		ID:         user.ID.String(),
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		IsReadonly: user.IsDemo,
	}}, http.StatusOK)
}

func (h *Handler) respondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to encode response", slog.Any("error", err))
	}
}

func (h *Handler) respondError(w http.ResponseWriter, message string, status int) {
	h.respondJSON(w, map[string]string{"error": message}, status)
}

func (h *Handler) setRefreshTokenCookie(w http.ResponseWriter, token string, expires time.Time) {
	if token == "" {
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expires.UTC(),
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *Handler) clearRefreshTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *Handler) refreshTokenFromRequest(r *http.Request, provided string) string {
	if token := strings.TrimSpace(provided); token != "" {
		return token
	}
	cookie, err := r.Cookie(refreshTokenCookieName)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(cookie.Value)
}
