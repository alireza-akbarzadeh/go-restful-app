package mocks

import (
	"context"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/stretchr/testify/mock"
)

type AttendeeRepositoryMock struct {
	mock.Mock
}

func (m *AttendeeRepositoryMock) Insert(ctx context.Context, attendee *models.Attendee) (*models.Attendee, error) {
	args := m.Called(ctx, attendee)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Attendee), args.Error(1)
}

func (m *AttendeeRepositoryMock) GetByEventAndUser(ctx context.Context, eventID, userID int) (*models.Attendee, error) {
	args := m.Called(ctx, eventID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Attendee), args.Error(1)
}

func (m *AttendeeRepositoryMock) GetByEventAndAttendee(ctx context.Context, eventID, userID int) (*models.Attendee, error) {
	args := m.Called(ctx, eventID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Attendee), args.Error(1)
}

func (m *AttendeeRepositoryMock) GetAttendeesByEvent(ctx context.Context, eventID int) ([]*models.User, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *AttendeeRepositoryMock) GetEventsByAttendee(ctx context.Context, userID int) ([]*models.Event, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *AttendeeRepositoryMock) GetEventByAttendee(ctx context.Context, userID int) ([]*models.Event, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *AttendeeRepositoryMock) Delete(ctx context.Context, userID, eventID int) error {
	args := m.Called(ctx, userID, eventID)
	return args.Error(0)
}

func (m *AttendeeRepositoryMock) DeleteByEvent(ctx context.Context, eventID int) error {
	args := m.Called(ctx, eventID)
	return args.Error(0)
}

func (m *AttendeeRepositoryMock) DeleteByUser(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}