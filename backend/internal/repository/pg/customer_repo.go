package pg

import (
	"context"
	"errors"
	stderrors "errors"

	"backend/internal/appcontext"
	"backend/internal/model"
	"backend/internal/repository"

	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) getTenantDB(ctx context.Context) (*gorm.DB, error) {
	ac, err := appcontext.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	return r.db.Where("tenant_id = ? AND deleted_at IS NULL", ac.TenantID), nil
}
func (r *CustomerRepository) FindByExternalID(ctx context.Context, externalID string) (*model.Customer, error) {
	db, err := r.getTenantDB(ctx)
	if err != nil {
		return nil, err
	}

	var customer model.Customer
	if err := db.Where("external_id = ?", externalID).First(&customer).Error; err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) FindOrCreate(ctx context.Context, externalID string, name *string) (*model.Customer, error) {
	ac, err := appcontext.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Try to find existing customer
	customer, err := r.FindByExternalID(ctx, externalID)
	if err == nil {
		return customer, nil
	}
	if err != repository.ErrNotFound {
		return nil, err
	}

	// Create new customer
	customer = &model.Customer{
		TenantID:   ac.TenantID,
		ExternalID: externalID,
		Name:       name,
	}

	if err := r.db.Create(customer).Error; err != nil {
		return nil, err
	}
	return customer, nil
}

func (r *CustomerRepository) Create(ctx context.Context, customer *model.Customer) error {
	ac, err := appcontext.FromContext(ctx)
	if err != nil {
		return err
	}
	customer.TenantID = ac.TenantID
	return r.db.Create(customer).Error
}

func (r *CustomerRepository) List(ctx context.Context) ([]model.Customer, error) {
	var out []model.Customer
	err := r.db.WithContext(ctx).Find(&out).Error
	return out, err
}

func (r *CustomerRepository) GetByID(ctx context.Context, id uint) (*model.Customer, error) {
	var c model.Customer
	err := r.db.WithContext(ctx).First(&c, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}
