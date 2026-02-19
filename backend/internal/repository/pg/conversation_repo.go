package pg

import (
	"context"
	"errors"
	"time"

	"backend/internal/appcontext"
	"backend/internal/model"
	"backend/internal/repository"

	"gorm.io/gorm"
)

type ConversationRepository struct {
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

func (r *ConversationRepository) getTenantDB(ctx context.Context) (*gorm.DB, error) {
	ac, err := appcontext.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	return r.db.Where("tenant_id = ? AND deleted_at IS NULL", ac.TenantID), nil
}

// FindOpenByCustomer uses PARTIAL INDEX for optimal performance
func (r *ConversationRepository) FindOpenByCustomer(
	ctx context.Context,
	customerID uint,
) (*model.Conversation, error) {

	var conv model.Conversation

	err := r.db.WithContext(ctx).
		Preload("Tenant").
		Preload("Customer").
		Where("customer_id = ? AND status = ?", customerID, model.StatusOpen).
		First(&conv).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &conv, nil
}
func (r *ConversationRepository) GetByID(
	ctx context.Context,
	id uint,
) (*model.Conversation, error) {

	var conv model.Conversation

	err := r.db.WithContext(ctx).
		Preload("Tenant").
		Preload("Customer").
		Preload("AssignedAgent").
		First(&conv, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &conv, nil
}
func (r *ConversationRepository) Create(ctx context.Context, conv *model.Conversation) error {
	ac, err := appcontext.FromContext(ctx)
	if err != nil {
		return err
	}
	conv.TenantID = ac.TenantID
	conv.Status = model.StatusOpen
	return r.db.Create(conv).Error
}

func (r *ConversationRepository) AssignAgent(ctx context.Context, convID uint, agentID uint) error {
	db, err := r.getTenantDB(ctx)
	if err != nil {
		return err
	}

	return db.Model(&model.Conversation{}).
		Where("id = ?", convID).
		Updates(map[string]interface{}{
			"assigned_agent_id": agentID,
			"status":            model.StatusAssigned,
		}).Error
}
func (r *ConversationRepository) UpdateLastMessageAt(ctx context.Context, convID uint, timestamp time.Time) error {
	db, err := r.getTenantDB(ctx)
	if err != nil {
		return err
	}

	return db.Model(&model.Conversation{}).
		Where("id = ?", convID).
		Update("last_message_at", timestamp).Error
}

func (r *ConversationRepository) List(
	ctx context.Context,
	filter repository.ConversationFilter,
	page, limit int,
) ([]*model.Conversation, int64, error) {

	var list []*model.Conversation
	var total int64

	q := r.db.WithContext(ctx).Model(&model.Conversation{})

	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
	}

	q.Count(&total)

	err := q.
		Preload("Tenant").
		Preload("Customer").
		Preload("AssignedAgent").
		Order("last_message_at DESC").
		Limit(limit).
		Offset((page-1)*limit).
		Find(&list).Error

	return list, total, err
}

func (r *ConversationRepository) GetDetail(ctx context.Context, id uint) (*model.Conversation, error) {
	var conv model.Conversation

	err := r.db.WithContext(ctx).
		Preload("Tenant").
		Preload("Customer").
		Preload("AssignedAgent").
		Where("id = ?", id).
		First(&conv).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &conv, nil
}

func (r *ConversationRepository) UpdateStatus(ctx context.Context, id uint, status model.ConversationStatus) error {
	return r.db.WithContext(ctx).
		Model(&model.Conversation{}).
		Where("id = ?", id).
		Update("status", status).Error
}

