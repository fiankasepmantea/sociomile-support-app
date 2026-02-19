package service

import (
	"context"

	"backend/internal/appcontext"
	"backend/internal/model"
	"backend/internal/repository/pg"
)

type TicketService struct {
	repo *pg.Repository
}

func NewTicketService(r *pg.Repository) *TicketService {
	return &TicketService{repo: r}
}

func (s *TicketService) Create(
	ctx context.Context,
	conversationID uint,
	title string,
	description string,
) (*model.Ticket, error) {

	ac, err := appcontext.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	t := &model.Ticket{
		TenantID:       ac.TenantID,
		ConversationID: conversationID,
		Title:          title,
		Description:    description,
		Status:         model.TicketStatusOpen,
	}

	err = s.repo.Ticket.Create(ctx, t)
	return t, err
}

func (s *TicketService) List(ctx context.Context) ([]model.Ticket, error) {
	return s.repo.Ticket.List(ctx)
}

func (s *TicketService) Get(ctx context.Context, id uint) (*model.Ticket, error) {
	return s.repo.Ticket.Get(ctx, id)
}

func (s *TicketService) UpdateStatus(ctx context.Context, id uint, status string) error {
	return s.repo.Ticket.UpdateStatus(ctx, id, model.TicketStatus(status))
}
