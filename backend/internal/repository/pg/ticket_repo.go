package pg

import (
	"context"

	"backend/internal/appcontext"
	"backend/internal/model"
	"backend/internal/repository"

	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) getTenantDB(ctx context.Context) (*gorm.DB, error) {
	ac, err := appcontext.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	return r.db.Where("tenant_id = ? AND deleted_at IS NULL", ac.TenantID), nil
}

// Add ticket methods as needed
func (r *TicketRepository) Create(ctx context.Context, ticket *model.Ticket) error {
	ac, err := appcontext.FromContext(ctx)
	if err != nil {
		return err
	}
	ticket.TenantID = ac.TenantID
	return r.db.Create(ticket).Error
}

func (r *TicketRepository) GetByID(ctx context.Context, id uint) (*model.Ticket, error) {
	db, err := r.getTenantDB(ctx)
	if err != nil {
		return nil, err
	}

	var t model.Ticket
	if err := db.First(&t, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *TicketRepository) List(ctx context.Context) ([]model.Ticket, error) {
	var out []model.Ticket
	err := r.db.WithContext(ctx).Find(&out).Error
	return out, err
}

func (r *TicketRepository) Get(ctx context.Context, id uint) (*model.Ticket, error) {
	var t model.Ticket
	err := r.db.WithContext(ctx).First(&t, id).Error
	return &t, err
}

func (r *TicketRepository) UpdateStatus(ctx context.Context, id uint, st model.TicketStatus) error {
	return r.db.WithContext(ctx).
		Model(&model.Ticket{}).
		Where("id = ?", id).
		Update("status", st).Error
}
