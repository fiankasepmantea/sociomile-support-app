package dto

import (
	"time"

	"backend/internal/model"
)
type ListResponse struct {
	Data   any   `json:"data"`
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}
type LoginResponse struct {
	Token string `json:"token"`
}
type TenantResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToTenantResponse(t *model.Tenant) *TenantResponse {
	if t == nil {
		return nil
	}

	return &TenantResponse{
		ID:        uint(t.ID),
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
type CustomerResponse struct {
	ID         uint      `json:"id"`
	TenantID   uint      `json:"tenant_id"`
	ExternalID string    `json:"external_id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
}

func ToCustomerResponse(c *model.Customer) *CustomerResponse {
	if c == nil {
		return nil
	}

	name := ""
	if c.Name != nil {
		name = *c.Name
	}

	return &CustomerResponse{
		ID:         uint(c.ID),
		TenantID:   uint(c.TenantID),
		ExternalID: c.ExternalID,
		Name:       name,
		CreatedAt:  c.CreatedAt,
	}
}
type UserResponse struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

func ToUserResponse(u *model.User) *UserResponse {
	if u == nil {
		return nil
	}

	return &UserResponse{
		ID:       u.ID,
		Email:    u.Email,
		Role:     u.Role,
		IsActive: u.IsActive,
	}
}
type MessageResponse struct {
	ID             uint      `json:"id"`
	ConversationID uint      `json:"conversation_id"`
	SenderType     string    `json:"sender_type"`
	Body           string    `json:"body"`
	CreatedAt      time.Time `json:"created_at"`
}

func ToMessageResponse(m *model.Message) MessageResponse {
	return MessageResponse{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		SenderType:     string(m.SenderType),
		Body:           m.Body,
		CreatedAt:      m.CreatedAt,
	}
}

func ToMessageResponseList(list []model.Message) []MessageResponse {
	out := make([]MessageResponse, 0, len(list))

	for i := range list {
		out = append(out, MessageResponse{
			ID:             list[i].ID,
			ConversationID: list[i].ConversationID,
			SenderType:     string(list[i].SenderType),
			Body:           list[i].Body,
			CreatedAt:      list[i].CreatedAt,
		})
	}

	return out
}
type ConversationResponse struct {
	ID              uint               `json:"id"`
	Tenant          *TenantResponse    `json:"tenant,omitempty"`
	Customer        *CustomerResponse  `json:"customer"`
	AssignedAgent   *UserResponse      `json:"assigned_agent,omitempty"`
	Status          string             `json:"status"`
	LastMessageAt   *time.Time         `json:"last_message_at,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	Messages        []MessageResponse  `json:"messages,omitempty"`
}

func ToConversationResponse(c *model.Conversation) ConversationResponse {
	resp := ConversationResponse{
		ID:            c.ID,
		Tenant:        ToTenantResponse(&c.Tenant),
		Customer:      ToCustomerResponse(&c.Customer),
		AssignedAgent: ToUserResponse(c.AssignedAgent),
		Status:        string(c.Status),
		LastMessageAt: c.LastMessageAt,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}

	if c.Messages != nil {
		resp.Messages = ToMessageResponseList(c.Messages)
	}

	return resp
}
func ToConversationResponseList(list []*model.Conversation) []ConversationResponse {
	out := make([]ConversationResponse, 0, len(list))

	for _, c := range list {
		if c != nil {
			out = append(out, ToConversationResponse(c))
		}
	}

	return out
}

type TicketResponse struct {
    ID             uint   `json:"id"`
    TenantID       uint   `json:"tenant_id"`
    ConversationID uint   `json:"conversation_id"`
    Title          string `json:"title"`
    Description    string `json:"description"`
    Status         string `json:"status"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}

func ToTicketResponse(t *model.Ticket) TicketResponse {
    return TicketResponse{
        ID: t.ID,
        TenantID: t.TenantID,
        ConversationID: t.ConversationID,
        Title: t.Title,
		Description: t.Description,
        Status: string(t.Status),
        CreatedAt: t.CreatedAt,
        UpdatedAt: t.UpdatedAt,
    }
}


