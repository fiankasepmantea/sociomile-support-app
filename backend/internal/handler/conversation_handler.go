package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"backend/internal/dto"
	"backend/internal/service"
)

type ConversationHandler struct {
	convService *service.ConversationService
}

func NewConversationHandler(
	conv *service.ConversationService,
) *ConversationHandler {
	return &ConversationHandler{
		convService: conv,
	}
}

func (h *ConversationHandler) ListConversations(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var statusPtr *string
	if s := c.Query("status"); s != "" {
		statusPtr = &s
	}

	convs, total, err := h.convService.ListConversations(
		c.Request.Context(),
		offset,
		limit,
		statusPtr,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// ✅ map ke DTO
	resp := make([]dto.ConversationResponse, 0, len(convs))
	for _, conv := range convs {
		resp = append(resp, dto.ToConversationResponse(conv))
	}

	c.JSON(http.StatusOK, dto.ListResponse{
		Data:   resp,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}
func (h *ConversationHandler) AssignAgent(c *gin.Context) {
	convID64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	var req struct {
		AgentID uint `json:"agent_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	convID64, _ = strconv.ParseUint(c.Param("id"), 10, 64)

	err = h.convService.AssignAgent(
		c.Request.Context(),
		uint(convID64),
		uint(req.AgentID),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "assigned"})
}

func (h *ConversationHandler) GetDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	conv, err := h.convService.GetConversationDetail(
		c.Request.Context(),
		uint(id),
	)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, dto.ToConversationResponse(conv))
}

func (h *ConversationHandler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.convService.UpdateStatus(
		c.Request.Context(),
		uint(id),
		req.Status,
	)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "updated"})
}
