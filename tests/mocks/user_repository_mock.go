package mocks

import (
	"context"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/alireza-akbarzadeh/ginflow/internal/pagination"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) Insert(ctx context.Context, user *models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepositoryMock) Get(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepositoryMock) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepositoryMock) GetById(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepositoryMock) UpdatePassword(ctx context.Context, userID int, hashedPassword string) error {
	args := m.Called(ctx, userID, hashedPassword)
	return args.Error(0)
}

func (m *UserRepositoryMock) GetAll(ctx context.Context) ([]*models.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *UserRepositoryMock) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepositoryMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *UserRepositoryMock) UpdateLastLogin(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *UserRepositoryMock) ListWithPagination(ctx context.Context, req *pagination.PaginationRequest) ([]*models.User, *pagination.PaginationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	if args.Get(1) == nil {
		return args.Get(0).([]*models.User), nil, args.Error(2)
	}
	return args.Get(0).([]*models.User), args.Get(1).(*pagination.PaginationResponse), args.Error(2)
}