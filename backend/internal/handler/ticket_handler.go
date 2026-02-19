package handler

import (
	"net/http"
	"strconv"

	"backend/internal/dto"
	"backend/internal/service"
	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	svc *service.TicketService
}

func NewTicketHandler(s *service.TicketService) *TicketHandler {
	return &TicketHandler{svc: s}
}

func (h *TicketHandler) List(c *gin.Context) {
	out, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, out)
}

func (h *TicketHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	t, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, t)
}

func (h *TicketHandler) Create(c *gin.Context) {
    var req dto.CreateTicketRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    t, err := h.svc.Create(c.Request.Context(), req.ConversationID, req.Title, req.Description)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, dto.ToTicketResponse(t))
}
func (h *TicketHandler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateStatus(c.Request.Context(), uint(id), req.Status); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "updated"})
}
