package handler

import (
	"net/http"
	"strconv"

	"backend/internal/dto"
	"backend/internal/service"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

func (h *MessageHandler) SendAgentMessage(c *gin.Context) {
	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	var req struct {
		Body string `json:"body" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	err = h.messageService.SendAgentMessage(
		c.Request.Context(),
		service.SendMessageRequest{
			ConversationID: uint(conversationID),
			Body:           req.Body,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "sent",
	})
}

func (h *MessageHandler) ListMessages(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	msgs, err := h.messageService.ListByConversation(
		c.Request.Context(),
		uint(id),
	)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, dto.ToMessageResponseList(msgs))
}

