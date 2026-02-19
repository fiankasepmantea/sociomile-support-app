package service

import (
	"context"
	"time"

	"backend/internal/appcontext"
	"backend/internal/apperrors"
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/repository/pg"
)

type ConversationService struct {
	repo *pg.Repository
}

func NewConversationService(repo *pg.Repository) *ConversationService {
	return &ConversationService{repo: repo}
}

type WebhookRequest struct {
	ExternalCustomerID string  `json:"external_customer_id" binding:"required"`
	CustomerName       *string `json:"customer_name"`
	MessageBody        string  `json:"message_body" binding:"required"`
}
func (s *ConversationService) ProcessWebhook(
	ctx context.Context,
	req WebhookRequest,
) (*model.Conversation, error) {

	ac, err := appcontext.FromContext(ctx)
	if err != nil {
		return nil, apperrors.ErrInternal("tenant context missing")
	}

	customer, err := s.repo.Customer.FindOrCreate(
		ctx,
		req.ExternalCustomerID,
		req.CustomerName,
	)
	if err != nil {
		return nil, apperrors.ErrInternal("failed to process customer")
	}

	conv, err := s.repo.Conversation.FindOpenByCustomer(
		ctx,
		customer.ID,
	)

	if err != nil && err != repository.ErrNotFound {
		return nil, apperrors.ErrInternal("failed to check conversations")
	}

	if err == repository.ErrNotFound {
		conv = &model.Conversation{
			TenantID:   ac.TenantID,
			CustomerID: customer.ID,
			Status:     model.StatusOpen,
		}

		if err := s.repo.Conversation.Create(ctx, conv); err != nil {
			return nil, apperrors.ErrInternal("failed to create conversation")
		}
	}

	msg := &model.Message{
		TenantID:       ac.TenantID,
		ConversationID: conv.ID,
		SenderType:     model.SenderCustomer,
		Body:           req.MessageBody,
	}

	if err := s.repo.Message.Create(ctx, msg); err != nil {
		return nil, apperrors.ErrInternal("failed to save message")
	}

	_ = s.repo.Conversation.UpdateLastMessageAt(
		ctx,
		conv.ID,
		time.Now(),
	)

	conv, err = s.repo.Conversation.GetByID(ctx, conv.ID)
	if err != nil {
		return nil, apperrors.ErrInternal("failed to refresh conversation")
	}

	return conv, nil
}
func (s *ConversationService) AssignAgent(
	ctx context.Context,
	convID uint,
	agentID uint,
) error {

	ac, err := appcontext.FromContext(ctx)
	if err != nil {
		return apperrors.ErrInternal("context error")
	}

	conv, err := s.repo.Conversation.GetByID(ctx, convID)
	if err != nil {
		return apperrors.ErrNotFound("conversation not found")
	}

	if conv.TenantID != ac.TenantID {
		return apperrors.ErrForbidden("cross tenant forbidden")
	}

	agent, err := s.repo.User.GetByID(ctx, agentID)
	if err != nil {
		return apperrors.ErrNotFound("agent not found")
	}

	if agent.TenantID != ac.TenantID {
		return apperrors.ErrForbidden("wrong tenant agent")
	}

	if agent.Role != "agent" && agent.Role != "admin" {
		return apperrors.ErrInvalid("not agent role")
	}

	if !agent.IsActive {
		return apperrors.ErrInvalid("agent inactive")
	}

	return s.repo.Conversation.AssignAgent(ctx, convID, agentID)
}

func (s *ConversationService) ListConversations(
	ctx context.Context,
	page, limit int,
	statusFilter *string,
) ([]*model.Conversation, int64, error) {

	var status *model.ConversationStatus

	if statusFilter != nil {
		switch *statusFilter {
		case "open":
			st := model.StatusOpen
			status = &st
		case "assigned":
			st := model.StatusAssigned
			status = &st
		case "closed":
			st := model.StatusClosed
			status = &st
		}
	}

	filter := repository.ConversationFilter{
		Status: status,
	}

	return s.repo.Conversation.List(ctx, filter, page, limit)
}

func (s *ConversationService) GetConversationDetail(
	ctx context.Context,
	id uint,
) (*model.Conversation, error) {

	return s.repo.Conversation.GetDetail(ctx, id)
}

func (s *ConversationService) UpdateStatus(
	ctx context.Context,
	id uint,
	status string,
) error {

	var st model.ConversationStatus

	switch status {
	case "open":
		st = model.StatusOpen
	case "assigned":
		st = model.StatusAssigned
	case "closed":
		st = model.StatusClosed
	default:
		return apperrors.ErrInvalid("invalid status")
	}

	return s.repo.Conversation.UpdateStatus(ctx, id, st)
}

func (s *ConversationService) CreateManual(
	ctx context.Context,
	customerID uint,
) (*model.Conversation, error) {

	ac, err := appcontext.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	conv := &model.Conversation{
		TenantID:   ac.TenantID,
		CustomerID: customerID,
		Status:     model.StatusOpen,
	}

	err = s.repo.Conversation.Create(ctx, conv)
	return conv, err
}
