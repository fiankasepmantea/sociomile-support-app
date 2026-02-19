package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AssignAgentRequest struct {
	AgentID uint `json:"agent_id" binding:"required"`
}

type SendMessageRequest struct {
	Body string `json:"body" binding:"required"`
}

type WebhookChannelRequest struct {
	ExternalCustomerID string  `json:"external_customer_id" binding:"required"`
	CustomerName       *string `json:"customer_name"`
	MessageBody        string  `json:"message_body" binding:"required"`
}

type CreateTicketRequest struct {
    ConversationID uint   `json:"conversation_id" binding:"required"`
    Title          string `json:"title" binding:"required"`
    Description    string `json:"description,omitempty"`
}


