package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/NirajDonga/dbpods/internal/repository"
	"github.com/gin-gonic/gin"
)

type PodHandler struct {
	podRepo *repository.PodRepository
}

func NewPodHandler(podRepo *repository.PodRepository) *PodHandler {
	return &PodHandler{podRepo: podRepo}
}

func (h *PodHandler) CreatePod(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}
	userID := userIDRaw.(int)

	tenantID := fmt.Sprintf("tenant-db-%d-%d", userID, time.Now().Unix())

	pod, err := h.podRepo.Create(c.Request.Context(), userID, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create pod in database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "pod created successfully",
		"pod":     pod,
	})
}

func (h *PodHandler) GetUserPods(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}
	userID := userIDRaw.(int)

	pods, err := h.podRepo.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch pods"})
		return
	}

	if pods == nil {
		c.JSON(http.StatusOK, gin.H{"pods": []interface{}{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pods": pods})
}
