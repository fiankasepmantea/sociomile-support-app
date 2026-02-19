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

type MessageService struct {
	repo *pg.Repository
}

func NewMessageService(repo *pg.Repository) *MessageService {
	return &MessageService{repo: repo}
}

type SendMessageRequest struct {
	ConversationID uint   `json:"conversation_id" binding:"required"`
	Body           string `json:"body" binding:"required,min=1,max=10000"`
}

// SendAgentMessage handles agent replies with strict validation
func (s *MessageService) SendAgentMessage(ctx context.Context, req SendMessageRequest) error {
	appCtx, err := appcontext.FromContext(ctx)
	if err != nil {
		return apperrors.ErrInternal("context error")
	}

	// 1. Get conversation (validates tenant ownership)
	conv, err := s.repo.Conversation.GetByID(ctx, req.ConversationID)
	if err != nil {
		if err == repository.ErrNotFound {
			return apperrors.ErrNotFound("conversation not found")
		}
		return apperrors.ErrInternal("database error")
	}

	// 2. 🔒 CRITICAL VALIDATION: Check if user can reply to this conversation
	if err := s.validateAgentPermission(ctx, appCtx, conv); err != nil {
		return err
	}

	// 3. 🔒 Check conversation status (cannot reply to closed conversations)
	if conv.Status == model.StatusClosed {
		return apperrors.ErrInvalid("cannot reply to closed conversation")
	}

	// 4. Create agent message
	message := &model.Message{
		ConversationID: req.ConversationID,
		SenderType:     model.SenderAgent,
		SenderUserID:   &appCtx.UserID,
		Body:           req.Body,
	}
	if err := s.repo.Message.Create(ctx, message); err != nil {
		return apperrors.ErrInternal("failed to save message")
	}

	// 5. Update conversation's last_message_at
	now := time.Now()
	if err := s.repo.Conversation.UpdateLastMessageAt(ctx, req.ConversationID, now); err != nil {
		// Non-critical: log but don't fail
	}

	// 6. Update conversation status to 'assigned' if still 'open'
	if conv.Status == model.StatusOpen {
		if err := s.repo.Conversation.UpdateStatus(ctx, req.ConversationID, model.StatusAssigned); err != nil {
			// Non-critical
		}
	}

	// 7. Create activity log
	payload := model.JSONB{
		"message_id":   message.ID,
		"body_preview": req.Body,
		"sent_by":      appCtx.UserID,
	}
	activityLog := &model.ActivityLog{
		Type:        "message.sent",
		EntityType:  "conversation",
		EntityID:    int64(req.ConversationID),
		ActorUserID: &appCtx.UserID,
		Payload:     payload,
	}
	s.repo.ActivityLog.Create(ctx, activityLog) // Fire-and-forget for audit

	return nil
}

// validateAgentPermission checks if current user can reply to conversation
func (s *MessageService) validateAgentPermission(ctx context.Context, appCtx *appcontext.AppContext, conv *model.Conversation) error {
	// Admins can reply to any conversation in tenant
	_ = ctx

	if appCtx.Role == string(model.RoleAdmin) {
		return nil
	}

	// Agents must be assigned to this conversation
	if appCtx.Role == string(model.RoleAgent) {
		if conv.AssignedAgentID == nil {
			return apperrors.ErrForbidden("conversation is not assigned to any agent")
		}
		if *conv.AssignedAgentID != appCtx.UserID {
			return apperrors.ErrForbidden("you are not assigned to this conversation")
		}
		return nil
	}

	// Other roles (e.g., customer) cannot send agent messages
	return apperrors.ErrForbidden("only agents or admins can send replies")
}

func (s *MessageService) ListByConversation(
	ctx context.Context,
	convID uint,
) ([]model.Message, error) {

	return s.repo.Message.ListByConversationID(ctx, convID)
}
