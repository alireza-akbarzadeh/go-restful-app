package mocks

import (
	"context"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/alireza-akbarzadeh/ginflow/internal/pagination"
	"github.com/stretchr/testify/mock"
)

type EventRepositoryMock struct {
	mock.Mock
}

func (m *EventRepositoryMock) Insert(ctx context.Context, event *models.Event) (*models.Event, error) {
	args := m.Called(ctx, event)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *EventRepositoryMock) Get(ctx context.Context, id int) (*models.Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *EventRepositoryMock) GetAll(ctx context.Context) ([]*models.Event, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *EventRepositoryMock) Update(ctx context.Context, event *models.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *EventRepositoryMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *EventRepositoryMock) ListWithPagination(ctx context.Context, req *pagination.PaginationRequest) ([]*models.Event, *pagination.PaginationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	if args.Get(1) == nil {
		return args.Get(0).([]*models.Event), nil, args.Error(2)
	}
	return args.Get(0).([]*models.Event), args.Get(1).(*pagination.PaginationResponse), args.Error(2)
}

func (m *EventRepositoryMock) ListWithAdvancedPagination(ctx context.Context, req *pagination.AdvancedPaginationRequest) ([]*models.Event, *pagination.AdvancedPaginatedResult, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	if args.Get(1) == nil {
		return args.Get(0).([]*models.Event), nil, args.Error(2)
	}
	return args.Get(0).([]*models.Event), args.Get(1).(*pagination.AdvancedPaginatedResult), args.Error(2)
}

func (m *EventRepositoryMock) GetByOwnerID(ctx context.Context, ownerID int) ([]*models.Event, error) {
	args := m.Called(ctx, ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Event), args.Error(1)
}
