package buckets

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"bucketbird/backend/internal/middleware"
	"bucketbird/backend/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	bucketService *service.BucketService
	encryptionKey []byte
	logger        *slog.Logger
}

func NewHandler(bucketService *service.BucketService, encryptionKey []byte, logger *slog.Logger) *Handler {
	return &Handler{
		bucketService: bucketService,
		encryptionKey: encryptionKey,
		logger:        logger,
	}
}

type BucketDTO struct {
	ID                 string  `json:"id"`
	Name               string  `json:"name"`
	Region             string  `json:"region"`
	Description        *string `json:"description"`
	Size               string  `json:"size"`
	SizeBytes          int64   `json:"sizeBytes"`
	CredentialID       string  `json:"credentialId"`
	CredentialName     string  `json:"credentialName"`
	CredentialProvider string  `json:"credentialProvider"`
	CreatedAt          string  `json:"createdAt"`
}

// formatByteSize formats bytes into human-readable format
func formatByteSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	if bytes == 0 {
		return "0 B"
	}

	if bytes < KB {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < MB {
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	} else if bytes < GB {
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	} else if bytes < TB {
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	}
	return fmt.Sprintf("%.2f TB", float64(bytes)/TB)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	buckets, err := h.bucketService.List(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to list buckets", slog.Any("error", err))
		h.respondError(w, "Failed to list buckets", http.StatusInternalServerError)
		return
	}

	dtos := make([]BucketDTO, len(buckets))
	for i, b := range buckets {
		dtos[i] = BucketDTO{
			ID:                 b.ID.String(),
			Name:               b.Name,
			Region:             b.Region,
			Description:        b.Description,
			Size:               formatByteSize(b.SizeBytes),
			SizeBytes:          b.SizeBytes,
			CredentialID:       b.CredentialID.String(),
			CredentialName:     b.CredentialName,
			CredentialProvider: b.CredentialProvider,
			CreatedAt:          b.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	h.respondJSON(w, map[string]interface{}{"buckets": dtos}, http.StatusOK)
}

type CreateBucketRequest struct {
	CredentialID string  `json:"credentialId"`
	Name         string  `json:"name"`
	Region       string  `json:"region"`
	Description  *string `json:"description"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateBucketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	credentialID, err := uuid.Parse(req.CredentialID)
	if err != nil {
		h.respondError(w, "Invalid credential ID", http.StatusBadRequest)
		return
	}

	bucket, err := h.bucketService.Create(r.Context(), service.CreateBucketInput{
		UserID:       userID,
		CredentialID: credentialID,
		Name:         req.Name,
		Region:       req.Region,
		Description:  req.Description,
	})
	if err != nil {
		if errors.Is(err, service.ErrCredentialNotFound) {
			h.respondError(w, "Credential not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrBucketAlreadyExists) {
			h.respondError(w, "Bucket already exists", http.StatusConflict)
			return
		}
		var provisionErr *service.BucketProvisionError
		if errors.As(err, &provisionErr) {
			h.respondError(w, provisionErr.Error(), http.StatusBadRequest)
			return
		}
		h.logger.Error("failed to create bucket", slog.Any("error", err))
		h.respondError(w, "Failed to create bucket", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, map[string]interface{}{"bucket": BucketDTO{
		ID:                 bucket.ID.String(),
		Name:               bucket.Name,
		Region:             bucket.Region,
		Description:        bucket.Description,
		Size:               formatByteSize(bucket.SizeBytes),
		SizeBytes:          bucket.SizeBytes,
		CredentialID:       bucket.CredentialID.String(),
		CredentialName:     bucket.CredentialName,
		CredentialProvider: bucket.CredentialProvider,
		CreatedAt:          bucket.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}}, http.StatusCreated)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	bucketID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, "Invalid bucket ID", http.StatusBadRequest)
		return
	}

	bucket, err := h.bucketService.Get(r.Context(), bucketID, userID)
	if err != nil {
		if errors.Is(err, service.ErrBucketNotFound) {
			h.respondError(w, "Bucket not found", http.StatusNotFound)
			return
		}
		h.logger.Error("failed to get bucket", slog.Any("error", err))
		h.respondError(w, "Failed to get bucket", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, map[string]interface{}{"bucket": BucketDTO{
		ID:                 bucket.ID.String(),
		Name:               bucket.Name,
		Region:             bucket.Region,
		Description:        bucket.Description,
		Size:               formatByteSize(bucket.SizeBytes),
		SizeBytes:          bucket.SizeBytes,
		CredentialID:       bucket.CredentialID.String(),
		CredentialName:     bucket.CredentialName,
		CredentialProvider: bucket.CredentialProvider,
		CreatedAt:          bucket.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}}, http.StatusOK)
}

type UpdateBucketRequest struct {
	Description *string `json:"description"`
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	bucketID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, "Invalid bucket ID", http.StatusBadRequest)
		return
	}

	var req UpdateBucketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.bucketService.Update(r.Context(), bucketID, userID, req.Description); err != nil {
		if errors.Is(err, service.ErrBucketNotFound) {
			h.respondError(w, "Bucket not found", http.StatusNotFound)
			return
		}
		h.logger.Error("failed to update bucket", slog.Any("error", err))
		h.respondError(w, "Failed to update bucket", http.StatusInternalServerError)
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

	bucketID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, "Invalid bucket ID", http.StatusBadRequest)
		return
	}

	deleteRemote := strings.EqualFold(r.URL.Query().Get("deleteRemote"), "true")

	if err := h.bucketService.Delete(r.Context(), bucketID, userID, deleteRemote); err != nil {
		if errors.Is(err, service.ErrBucketNotFound) {
			h.respondError(w, "Bucket not found", http.StatusNotFound)
			return
		}
		h.logger.Error("failed to delete bucket", slog.Any("error", err))
		h.respondError(w, "Failed to delete bucket", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RecalculateSize(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	bucketID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.respondError(w, "Invalid bucket ID", http.StatusBadRequest)
		return
	}

	if err := h.bucketService.RecalculateBucketSize(r.Context(), bucketID, userID, h.encryptionKey); err != nil {
		if errors.Is(err, service.ErrBucketNotFound) {
			h.respondError(w, "Bucket not found", http.StatusNotFound)
			return
		}
		h.logger.Error("failed to recalculate bucket size", slog.Any("error", err))
		h.respondError(w, "Failed to recalculate bucket size", http.StatusInternalServerError)
		return
	}

	// Get updated bucket info
	bucket, err := h.bucketService.Get(r.Context(), bucketID, userID)
	if err != nil {
		h.logger.Error("failed to get bucket after size calculation", slog.Any("error", err))
		h.respondError(w, "Failed to get updated bucket", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, map[string]interface{}{"bucket": BucketDTO{
		ID:                 bucket.ID.String(),
		Name:               bucket.Name,
		Region:             bucket.Region,
		Description:        bucket.Description,
		Size:               formatByteSize(bucket.SizeBytes),
		SizeBytes:          bucket.SizeBytes,
		CredentialID:       bucket.CredentialID.String(),
		CredentialName:     bucket.CredentialName,
		CredentialProvider: bucket.CredentialProvider,
		CreatedAt:          bucket.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
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
