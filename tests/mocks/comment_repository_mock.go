package mocks

import (
	"context"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/stretchr/testify/mock"
)

type CommentRepositoryMock struct {
	mock.Mock
}

func (m *CommentRepositoryMock) Insert(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	args := m.Called(ctx, comment)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Comment), args.Error(1)
}

func (m *CommentRepositoryMock) GetByEvent(ctx context.Context, eventID int) ([]*models.Comment, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Comment), args.Error(1)
}

func (m *CommentRepositoryMock) Get(ctx context.Context, id int) (*models.Comment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Comment), args.Error(1)
}

func (m *CommentRepositoryMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
