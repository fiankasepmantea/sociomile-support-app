package pg

import (
	"context"
	stderrors "errors"

	"backend/internal/appcontext"
	"backend/internal/model"
	"backend/internal/repository"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) getTenantDB(ctx context.Context) (*gorm.DB, error) {
	ac, err := appcontext.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	return r.db.Where("tenant_id = ? AND deleted_at IS NULL", ac.TenantID), nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	db, err := r.getTenantDB(ctx)
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	db, err := r.getTenantDB(ctx)
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	ac, err := appcontext.FromContext(ctx)
	if err != nil {
		return err
	}
	user.TenantID = ac.TenantID
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	db, err := r.getTenantDB(ctx)
	if err != nil {
		return err
	}

	return db.Save(user).Error
}