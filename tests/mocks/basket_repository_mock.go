package mocks

import (
	"context"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/stretchr/testify/mock"
)

type BasketRepositoryMock struct {
	mock.Mock
}

func (m *BasketRepositoryMock) GetActiveBasket(ctx context.Context, userID int) (*models.Basket, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Basket), args.Error(1)
}

func (m *BasketRepositoryMock) CreateBasket(ctx context.Context, basket *models.Basket) error {
	args := m.Called(ctx, basket)
	return args.Error(0)
}

func (m *BasketRepositoryMock) AddItem(ctx context.Context, basketID int, item *models.BasketItem) error {
	args := m.Called(ctx, basketID, item)
	return args.Error(0)
}

func (m *BasketRepositoryMock) UpdateItemQuantity(ctx context.Context, itemID int, quantity int) error {
	args := m.Called(ctx, itemID, quantity)
	return args.Error(0)
}

func (m *BasketRepositoryMock) RemoveItem(ctx context.Context, itemID int) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}

func (m *BasketRepositoryMock) ClearBasket(ctx context.Context, basketID int) error {
	args := m.Called(ctx, basketID)
	return args.Error(0)
}
