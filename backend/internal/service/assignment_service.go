package service

import (
	"context"
	"fmt"

	"backend/internal/appcontext"
	"backend/internal/apperrors"
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/repository/pg"
)

type AssignmentService struct {
	repo *pg.Repository
}

func NewAssignmentService(repo *pg.Repository) *AssignmentService {
	return &AssignmentService{repo: repo}
}

// AssignConversation handles agent assignment with full validation and audit trail
func (s *AssignmentService) AssignConversation(ctx context.Context, convID uint, agentID uint) error {
	appCtx, err := appcontext.FromContext(ctx)
	if err != nil {
		return apperrors.ErrInternal("context error")
	}

	// 1. Get conversation (validates tenant ownership)
	conv, err := s.repo.Conversation.GetByID(ctx, convID)
	if err != nil {
		if err == repository.ErrNotFound {
			return apperrors.ErrNotFound("conversation not found")
		}
		return apperrors.ErrInternal("database error")
	}

	// 2. Get agent details (critical: verify same tenant!)
	agent, err := s.repo.User.GetByID(ctx, agentID)
	if err != nil {
		if err == repository.ErrNotFound {
			return apperrors.ErrNotFound("agent not found")
		}
		return apperrors.ErrInternal("database error")
	}

	// 🔒 CROSS-TENANT GUARD: Most critical validation
	if agent.TenantID != appCtx.TenantID {
		return apperrors.ErrForbidden("agent does not belong to this tenant")
	}

	// 🔒 ROLE VALIDATION
	if agent.Role != "agent" && agent.Role != "admin" {
		return apperrors.ErrInvalid("user is not an agent")
	}

	// 🔒 ACTIVE STATUS CHECK
	if !agent.IsActive {
		return apperrors.ErrInvalid("agent account is inactive")
	}

	// 3. Perform assignment (updates status to 'assigned')
	if err := s.repo.Conversation.AssignAgent(ctx, convID, agentID); err != nil {
		return apperrors.ErrInternal("failed to assign agent")
	}

	// 4. CREATE ACTIVITY LOG (audit trail)
	payload := model.JSONB{
		"previous_agent_id": conv.AssignedAgentID,
		"new_agent_id":      agentID,
		"assigned_by":       appCtx.UserID,
	}

	activityLog := &model.ActivityLog{
		Type:        "conversation.assigned",
		EntityType:  "conversation",
		EntityID:    int64(convID),
		ActorUserID: &appCtx.UserID,
		Payload:     payload,
	}

	if err := s.repo.ActivityLog.Create(ctx, activityLog); err != nil {
		return fmt.Errorf("assignment succeeded but audit log failed: %w", err)
	}

	return nil
}