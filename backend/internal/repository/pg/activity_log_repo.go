package pg

import (
	"context"
	"time"

	"backend/internal/appcontext"
	"backend/internal/model"
	"gorm.io/gorm"
)

type ActivityLogRepository struct {
	db *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) *ActivityLogRepository {
	return &ActivityLogRepository{db: db}
}

func (r *ActivityLogRepository) Create(ctx context.Context, log *model.ActivityLog) error {
	ac, err := appcontext.FromContext(ctx)
	if err != nil {
		return err
	}

	// Auto-populate tenant and timestamp
	log.TenantID = ac.TenantID
	log.CreatedAt = time.Now()

	// If actor not set, use context user (for system events)
	if log.ActorUserID == nil && ac.UserID != 0 {
		log.ActorUserID = &ac.UserID
	}

	return r.db.Create(log).Error
}