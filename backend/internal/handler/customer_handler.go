package handler

import (
	"strconv"

	"backend/internal/service"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	service *service.CustomerService
}

func NewCustomerHandler(s *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: s}
}

func (h *CustomerHandler) List(c *gin.Context) {
	out, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, out)
}

func (h *CustomerHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	cust, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, cust)
}

func (h *CustomerHandler) Create(c *gin.Context) {
	var req struct {
		ExternalID string `json:"external_id"`
		Name       string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	out, err := h.service.Create(c.Request.Context(), req.ExternalID, req.Name)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, out)
}
