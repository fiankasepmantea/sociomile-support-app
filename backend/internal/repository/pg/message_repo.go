package pg

import (
	"context"
	"time"

	"backend/internal/appcontext"
	"backend/internal/model"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(ctx context.Context, message *model.Message) error {
	ac, err := appcontext.FromContext(ctx)
	if err != nil {
		return err
	}

	message.TenantID = ac.TenantID
	message.CreatedAt = time.Now()

	return r.db.Create(message).Error
}

func (r *MessageRepository) ListByConversation(ctx context.Context, convID uint, limit, offset int) ([]*model.Message, error) {
	db, err := r.getTenantDB(ctx)
	if err != nil {
		return nil, err
	}

	var messages []*model.Message
	err = db.Where("conversation_id = ?", convID).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error

	return messages, err
}

// Helper method (shared with other repos)
func (r *MessageRepository) getTenantDB(ctx context.Context) (*gorm.DB, error) {
	ac, err := appcontext.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	return r.db.Where("tenant_id = ? AND deleted_at IS NULL", ac.TenantID), nil
}

func (r *MessageRepository) ListByConversationID(
	ctx context.Context,
	convID uint,
) ([]model.Message, error) {

	var msgs []model.Message

	err := r.db.WithContext(ctx).
		Where("conversation_id = ?", convID).
		Order("id ASC").
		Find(&msgs).Error

	return msgs, err
}
