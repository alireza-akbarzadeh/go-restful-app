package mocks

import (
	"context"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/stretchr/testify/mock"
)

type ProfileRepositoryMock struct {
	mock.Mock
}

func (m *ProfileRepositoryMock) Insert(ctx context.Context, profile *models.Profile) (*models.Profile, error) {
	args := m.Called(ctx, profile)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Profile), args.Error(1)
}

func (m *ProfileRepositoryMock) GetByUserID(ctx context.Context, userID int) (*models.Profile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Profile), args.Error(1)
}

func (m *ProfileRepositoryMock) GetByUserIDWithUser(ctx context.Context, userID int) (*models.Profile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Profile), args.Error(1)
}

func (m *ProfileRepositoryMock) Update(ctx context.Context, profile *models.Profile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *ProfileRepositoryMock) UpdateByUserID(ctx context.Context, userID int, updates map[string]interface{}) error {
	args := m.Called(ctx, userID, updates)
	return args.Error(0)
}

func (m *ProfileRepositoryMock) DeleteByUserID(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
