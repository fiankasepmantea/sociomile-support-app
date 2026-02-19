package handler

import (
	"net/http"

	"backend/internal/apperrors"
	"backend/internal/service"
	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	convService *service.ConversationService
}

func NewWebhookHandler(cs *service.ConversationService) *WebhookHandler {
	return &WebhookHandler{convService: cs}
}

// HandleChannel receives messages from external channels (WhatsApp, FB, etc)
func (h *WebhookHandler) HandleChannel(c *gin.Context) {
	var req service.WebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Process webhook (tenant context already injected by middleware)
	conv, err := h.convService.ProcessWebhook(c.Request.Context(), req)
	if err != nil {
		appErr, ok := err.(*apperrors.AppError)
		if ok {
			c.JSON(appErr.Status, gin.H{"error": appErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"conversation_id": conv.ID,
		"status":          conv.Status,
		"customer_id":     conv.CustomerID,
	})
}