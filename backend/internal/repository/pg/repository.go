package pg

import "gorm.io/gorm"

type Repository struct {
	User         *UserRepository
	Conversation *ConversationRepository
	Customer     *CustomerRepository
	Message      *MessageRepository
	Ticket       *TicketRepository
	ActivityLog  *ActivityLogRepository
}

func NewPgRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:         NewUserRepository(db),
		Conversation: NewConversationRepository(db),
		Customer:     NewCustomerRepository(db),
		Message:      NewMessageRepository(db),
		Ticket:       NewTicketRepository(db),
		ActivityLog:  NewActivityLogRepository(db),
	}
}
