package mocks

import (
	"context"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/stretchr/testify/mock"
)

type ProductRepositoryMock struct {
	mock.Mock
}

func (m *ProductRepositoryMock) Insert(ctx context.Context, product *models.Product) (*models.Product, error) {
	args := m.Called(ctx, product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *ProductRepositoryMock) GetAll(ctx context.Context, page, limit int, search string, categoryID int) ([]models.Product, int64, error) {
	args := m.Called(ctx, page, limit, search, categoryID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

func (m *ProductRepositoryMock) Get(ctx context.Context, id int) (*models.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *ProductRepositoryMock) Update(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *ProductRepositoryMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ProductRepositoryMock) GetBySlug(ctx context.Context, slug string) (*models.Product, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *ProductRepositoryMock) GetByUser(ctx context.Context, userID int) ([]models.Product, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *ProductRepositoryMock) GetByCategory(ctx context.Context, categoryID int) ([]models.Product, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *ProductRepositoryMock) GetByCategoryID(ctx context.Context, categoryID int) ([]*models.Product, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Product), args.Error(1)
}
