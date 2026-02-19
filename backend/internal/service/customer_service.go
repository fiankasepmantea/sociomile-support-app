package service

import (
	"context"

	"backend/internal/model"
	"backend/internal/repository/pg"
)

type CustomerService struct {
	repo *pg.Repository
}

func NewCustomerService(r *pg.Repository) *CustomerService {
	return &CustomerService{repo: r}
}

func (s *CustomerService) List(ctx context.Context) ([]model.Customer, error) {
	return s.repo.Customer.List(ctx)
}

func (s *CustomerService) Get(ctx context.Context, id uint) (*model.Customer, error) {
	return s.repo.Customer.GetByID(ctx, id)
}

func (s *CustomerService) Create(
	ctx context.Context,
	extID string,
	name string,
) (*model.Customer, error) {

	c := &model.Customer{
		ExternalID: extID,
		Name:       &name,
	}

	err := s.repo.Customer.Create(ctx, c)
	return c, err
}
