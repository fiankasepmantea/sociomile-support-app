package repository

import (
	"context"
	"errors"
	"backend/internal/model"
	"time"
)

// Common repository errors
var (
	ErrNotFound = errors.New("record not found")
	ErrConflict = errors.New("conflict")
	ErrInvalid  = errors.New("invalid input")
)

type ConversationFilter struct {
	Status  *model.ConversationStatus
	AgentID *uint
}

// UserRepository interface
type UserRepository interface {
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
}

// ConversationRepository interface
type ConversationRepository interface {
	FindOpenByCustomer(ctx context.Context, customerID uint) (*model.Conversation, error)
	GetByID(ctx context.Context, id uint) (*model.Conversation, error)
	Create(ctx context.Context, conv *model.Conversation) error
	AssignAgent(ctx context.Context, convID uint, agentID uint) error
	UpdateStatus(ctx context.Context, convID uint, status model.ConversationStatus) error
	UpdateLastMessageAt(ctx context.Context, convID uint, timestamp time.Time) error
	List(ctx context.Context, filter ConversationFilter, page, limit int) ([]*model.Conversation, int64, error)
}

// CustomerRepository interface
type CustomerRepository interface {
	GetByID(ctx context.Context, id uint) (*model.Customer, error)
	FindByExternalID(ctx context.Context, externalID string) (*model.Customer, error)
	FindOrCreate(ctx context.Context, externalID string, name *string) (*model.Customer, error)
	Create(ctx context.Context, customer *model.Customer) error
}

// MessageRepository interface
type MessageRepository interface {
	Create(ctx context.Context, message *model.Message) error
	ListByConversation(ctx context.Context, convID uint, limit, offset int) ([]*model.Message, error)
}

// TicketRepository interface
type TicketRepository interface {
	Create(ctx context.Context, ticket *model.Ticket) error
}

// ActivityLogRepository interface
type ActivityLogRepository interface {
	Create(ctx context.Context, log *model.ActivityLog) error
}