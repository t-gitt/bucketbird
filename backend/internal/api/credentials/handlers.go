package credentials

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"bucketbird/backend/internal/middleware"
	"bucketbird/backend/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	credentialService *service.CredentialService
	logger            *slog.Logger
}

func NewHandler(credentialService *service.CredentialService, logger *slog.Logger) *Handler {
	return &Handler{
		credentialService: credentialService,
		logger:            logger,
	}
}

type CredentialDTO struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Provider  string  `json:"provider"`
	Region    string  `json:"region"`
	Endpoint  string  `json:"endpoint"`
	UseSSL    bool    `json:"useSSL"`
	Status    string  `json:"status"`
	Logo      *string `json:"logo"`
	CreatedAt string  `json:"createdAt"`
}

type DiscoveredBucketDTO struct {
	Name      string  `json:"name"`
	CreatedAt *string `json:"createdAt,omitempty"`
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	credentials, err := h.credentialService.List(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to list credentials", slog.Any("error", err))
		h.respondError(w, "Failed to list credentials", http.StatusInternalServerError)
		return
	}

	dtos := make([]CredentialDTO, len(credentials))
	for i, c := range credentials {
		dtos[i] = CredentialDTO{
			ID:        c.ID.String(),
			Name:      c.Name,
			Provider:  c.Provider,
			Region:    c.Region,
			Endpoint:  c.Endpoint,
			UseSSL:    c.UseSSL,
			Status:    c.Status,
			Logo:      c.Logo,
			CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	h.respondJSON(w, map[string]interface{}{"credentials": dtos}, http.StatusOK)
}

type CreateCredentialRequest struct {
	Name      string  `json:"name"`
	Provider  string  `json:"provider"`
	Region    string  `json:"region"`
	Endpoint  string  `json:"endpoint"`
	AccessKey string  `json:"accessKey"`
	SecretKey string  `json:"secretKey"`
	UseSSL    bool    `json:"useSSL"`
	Logo      *string `json:"logo"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	credential, err := h.credentialService.Create(r.Context(), service.CreateCredentialInput{
		UserID:    userID,
		Name:      req.Name,
		Provider:  req.Provider,
		Region:    req.Region,
		Endpoint:  req.Endpoint,
		AccessKey: req.AccessKey,
		SecretKey: req.SecretKey,
		UseSSL:    req.UseSSL,
		Logo:      req.Logo,
	})
	if err != nil {
		h.logger.Error("failed to create credential", slog.Any("error", err))
		h.respondError(w, "Failed to create credential", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, map[string]interface{}{"credential": CredentialDTO{
		ID:        credential.ID.String(),
		Name:      credential.Name,
		Provider:  credential.Provider,
		Region:    credential.Region,
		Endpoint:  credential.Endpoint,
		UseSSL:    credential.UseSSL,
		Status:    credential.Status,
		Logo:      credential.Logo,
		CreatedAt: credential.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}}, http.StatusCreated)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	credentialID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, "Invalid credential ID", http.StatusBadRequest)
		return
	}

	credential, err := h.credentialService.Get(r.Context(), credentialID, userID)
	if err != nil {
		if errors.Is(err, service.ErrCredentialNotFound) {
			h.respondError(w, "Credential not found", http.StatusNotFound)
			return
		}
		h.logger.Error("failed to get credential", slog.Any("error", err))
		h.respondError(w, "Failed to get credential", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, map[string]interface{}{"credential": CredentialDTO{
		ID:        credential.ID.String(),
		Name:      credential.Name,
		Provider:  credential.Provider,
		Region:    credential.Region,
		Endpoint:  credential.Endpoint,
		UseSSL:    credential.UseSSL,
		Status:    credential.Status,
		Logo:      credential.Logo,
		CreatedAt: credential.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}}, http.StatusOK)
}

type UpdateCredentialRequest struct {
	Name      string  `json:"name"`
	Provider  string  `json:"provider"`
	Region    string  `json:"region"`
	Endpoint  string  `json:"endpoint"`
	AccessKey string  `json:"accessKey"`
	SecretKey string  `json:"secretKey"`
	UseSSL    bool    `json:"useSSL"`
	Logo      *string `json:"logo"`
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	credentialID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, "Invalid credential ID", http.StatusBadRequest)
		return
	}

	var req UpdateCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.credentialService.Update(r.Context(), service.UpdateCredentialInput{
		ID:        credentialID,
		UserID:    userID,
		Name:      req.Name,
		Provider:  req.Provider,
		Region:    req.Region,
		Endpoint:  req.Endpoint,
		AccessKey: req.AccessKey,
		SecretKey: req.SecretKey,
		UseSSL:    req.UseSSL,
		Logo:      req.Logo,
	}); err != nil {
		if errors.Is(err, service.ErrCredentialNotFound) {
			h.respondError(w, "Credential not found", http.StatusNotFound)
			return
		}
		h.logger.Error("failed to update credential", slog.Any("error", err))
		h.respondError(w, "Failed to update credential", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	credentialID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, "Invalid credential ID", http.StatusBadRequest)
		return
	}

	if err := h.credentialService.Delete(r.Context(), credentialID, userID); err != nil {
		if errors.Is(err, service.ErrCredentialNotFound) {
			h.respondError(w, "Credential not found", http.StatusNotFound)
			return
		}
		h.logger.Error("failed to delete credential", slog.Any("error", err))
		h.respondError(w, "Failed to delete credential", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	credentialID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, "Invalid credential ID", http.StatusBadRequest)
		return
	}

	result, err := h.credentialService.Test(r.Context(), credentialID, userID)
	if err != nil {
		h.logger.Error("failed to test credential", slog.Any("error", err))
		h.respondError(w, "Failed to test credential", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, map[string]interface{}{"result": result}, http.StatusOK)
}

func (h *Handler) DiscoverBuckets(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	credentialID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, "Invalid credential ID", http.StatusBadRequest)
		return
	}

	buckets, err := h.credentialService.DiscoverBuckets(r.Context(), credentialID, userID)
	if err != nil {
		if errors.Is(err, service.ErrCredentialNotFound) {
			h.respondError(w, "Credential not found", http.StatusNotFound)
			return
		}
		var discoveryErr *service.CredentialDiscoveryError
		if errors.As(err, &discoveryErr) {
			h.respondError(w, discoveryErr.Error(), http.StatusBadRequest)
			return
		}
		h.logger.Error("failed to discover buckets for credential", slog.Any("error", err))
		h.respondError(w, "Failed to discover buckets", http.StatusInternalServerError)
		return
	}

	dtos := make([]DiscoveredBucketDTO, 0, len(buckets))
	for _, bucket := range buckets {
		dto := DiscoveredBucketDTO{
			Name: bucket.Name,
		}
		if bucket.CreatedAt != nil {
			formatted := bucket.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00")
			dto.CreatedAt = &formatted
		}
		dtos = append(dtos, dto)
	}

	h.respondJSON(w, map[string]interface{}{"buckets": dtos}, http.StatusOK)
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
