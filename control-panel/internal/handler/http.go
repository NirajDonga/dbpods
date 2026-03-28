package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/NirajDonga/dbpods/internal/core"
)

var validTenantID = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

type APIHandler struct {
	provisioner core.DBProvisioner
}

func NewAPIHandler(p core.DBProvisioner) *APIHandler {
	return &APIHandler{provisioner: p}
}

type CreateDBRequest struct {
	TenantID string `json:"tenantId"`
}

func (h *APIHandler) HandleCreateDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateDBRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	req.TenantID = strings.TrimSpace(req.TenantID)
	if req.TenantID == "" {
		http.Error(w, "tenantId is required", http.StatusBadRequest)
		return
	}
	if len(req.TenantID) > 63 {
		http.Error(w, "tenantId must be 63 characters or fewer", http.StatusBadRequest)
		return
	}
	if !validTenantID.MatchString(req.TenantID) {
		http.Error(w, "tenantId may only contain letters, digits, hyphens, and underscores, and must start with a letter or digit", http.StatusBadRequest)
		return
	}

	dbPassword, err := generateSecurePassword(16)
	if err != nil {
		http.Error(w, "Failed to generate password", http.StatusInternalServerError)
		return
	}

	err = h.provisioner.CreateTenantDatabase(r.Context(), req.TenantID, dbPassword)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to provision DB: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Database provisioning successfully started",
		"tenant":   req.TenantID,
		"password": dbPassword,
	})
}

func generateSecurePassword(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}
