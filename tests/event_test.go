package tests

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/alireza-akbarzadeh/ginflow/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestEventManagement tests the complete event management flow
func TestEventManagement(t *testing.T) {
	ts := SetupMockTestSuite(t)

	// Create a test user and get token
	userID := 1
	token, err := ts.GenerateToken(userID)
	assert.NoError(t, err)

	mockUserRepo := ts.Mocks.Users.(*mocks.UserRepositoryMock)
	mockUserRepo.On("Get", mock.Anything, userID).Return(&models.User{ID: userID, Email: "eventuser@example.com", Name: "Event User"}, nil)

	mockEventRepo := ts.Mocks.Events.(*mocks.EventRepositoryMock)

	t.Run("create event", func(t *testing.T) {
		event := models.Event{
			Name:        "Test Event",
			Description: "This is a test event description",
			Date:        "2025-12-31",
			Location:    "Test Location",
		}

		mockEventRepo.On("Insert", mock.Anything, mock.MatchedBy(func(e *models.Event) bool {
			return e.Name == event.Name && e.OwnerID == userID
		})).Return(&models.Event{ID: 1, Name: event.Name, Description: event.Description, Date: event.Date, Location: event.Location, OwnerID: userID}, nil).Once()

		w := ts.createAuthenticatedRequest("POST", "/api/v1/events", token, event)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createdEvent models.Event
		err := json.Unmarshal(w.Body.Bytes(), &createdEvent)
		assert.NoError(t, err)
		assert.Equal(t, event.Name, createdEvent.Name)
	})

	t.Run("get all events", func(t *testing.T) {
		events := []*models.Event{
			{Name: "Event 1", OwnerID: userID},
		}
		mockEventRepo.On("GetAll", mock.Anything).Return(events, nil).Once()

		w := ts.createRequest("GET", "/api/v1/events", nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var respEvents []models.Event
		err := json.Unmarshal(w.Body.Bytes(), &respEvents)
		assert.NoError(t, err)
		assert.Equal(t, len(events), len(respEvents))
	})

	t.Run("get single event", func(t *testing.T) {
		eventID := 1
		event := &models.Event{
			ID:          eventID,
			Name:        "Single Event",
			Description: "Description for single event",
			Date:        "2025-12-31",
			Location:    "Single Location",
			OwnerID:     userID,
		}

		mockEventRepo.On("Get", mock.Anything, eventID).Return(event, nil).Once()

		w := ts.createRequest("GET", "/api/v1/events/"+strconv.Itoa(eventID), nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var retrievedEvent models.Event
		err := json.Unmarshal(w.Body.Bytes(), &retrievedEvent)
		assert.NoError(t, err)
		assert.Equal(t, event.ID, retrievedEvent.ID)
	})

	t.Run("update event", func(t *testing.T) {
		eventID := 1
		event := &models.Event{
			ID:      eventID,
			Name:    "Original Event",
			OwnerID: userID,
		}

		updatedEvent := models.Event{
			Name:        "Updated Event",
			Description: "Updated description",
			Date:        "2025-12-31",
			Location:    "Updated Location",
		}

		mockEventRepo.On("Get", mock.Anything, eventID).Return(event, nil).Once()
		mockEventRepo.On("Update", mock.Anything, mock.MatchedBy(func(e *models.Event) bool {
			return e.ID == eventID && e.Name == updatedEvent.Name
		})).Return(nil).Once()

		w := ts.createAuthenticatedRequest("PUT", "/api/v1/events/"+strconv.Itoa(eventID), token, updatedEvent)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("delete event", func(t *testing.T) {
		eventID := 2
		event := &models.Event{
			ID:      eventID,
			Name:    "Event to Delete",
			OwnerID: userID,
		}

		mockAttendeeRepo := ts.Mocks.Attendees.(*mocks.AttendeeRepositoryMock)

		mockEventRepo.On("Get", mock.Anything, eventID).Return(event, nil).Once()
		mockAttendeeRepo.On("DeleteByEvent", mock.Anything, eventID).Return(nil).Once()
		mockEventRepo.On("Delete", mock.Anything, eventID).Return(nil).Once()

		w := ts.createAuthenticatedRequest("DELETE", "/api/v1/events/"+strconv.Itoa(eventID), token, nil)
		assert.Equal(t, http.StatusNoContent, w.Code)

		// Try to get the deleted event (should fail)
		mockEventRepo.On("Get", mock.Anything, eventID).Return(nil, nil).Once()
		w = ts.createRequest("GET", "/api/v1/events/"+strconv.Itoa(eventID), nil)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestEventAuthorization tests that users can only modify their own events
func TestEventAuthorization(t *testing.T) {
	ts := SetupMockTestSuite(t)

	user1ID := 1
	user2ID := 2
	token1, _ := ts.GenerateToken(user1ID)
	token2, _ := ts.GenerateToken(user2ID)

	mockUserRepo := ts.Mocks.Users.(*mocks.UserRepositoryMock)
	mockUserRepo.On("Get", mock.Anything, user1ID).Return(&models.User{ID: user1ID, Email: "user1@example.com"}, nil)
	mockUserRepo.On("Get", mock.Anything, user2ID).Return(&models.User{ID: user2ID, Email: "user2@example.com"}, nil)

	mockEventRepo := ts.Mocks.Events.(*mocks.EventRepositoryMock)

	eventID := 1
	event := &models.Event{
		ID:      eventID,
		Name:    "User1 Event",
		OwnerID: user1ID,
	}

	t.Run("user cannot update another user's event", func(t *testing.T) {
		updatedEvent := models.Event{
			Name:        "Hacked Event",
			Description: "This should not work",
			Date:        "2025-12-31",
			Location:    "Hacked Location",
		}

		mockEventRepo.On("Get", mock.Anything, eventID).Return(event, nil).Once()

		w := ts.createAuthenticatedRequest("PUT", "/api/v1/events/"+strconv.Itoa(eventID), token2, updatedEvent)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("user cannot delete another user's event", func(t *testing.T) {
		mockEventRepo.On("Get", mock.Anything, eventID).Return(event, nil).Once()

		w := ts.createAuthenticatedRequest("DELETE", "/api/v1/events/"+strconv.Itoa(eventID), token2, nil)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("owner can update their event", func(t *testing.T) {
		updatedEvent := models.Event{
			Name:        "Updated by Owner",
			Description: "Updated description by owner",
			Date:        "2025-12-31",
			Location:    "Updated Location",
		}

		mockEventRepo.On("Get", mock.Anything, eventID).Return(event, nil).Once()
		mockEventRepo.On("Update", mock.Anything, mock.MatchedBy(func(e *models.Event) bool {
			return e.ID == eventID && e.Name == updatedEvent.Name
		})).Return(nil).Once()

		w := ts.createAuthenticatedRequest("PUT", "/api/v1/events/"+strconv.Itoa(eventID), token1, updatedEvent)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestEventValidation tests input validation for events
func TestEventValidation(t *testing.T) {
	ts := SetupMockTestSuite(t)

	userID := 1
	token, _ := ts.GenerateToken(userID)

	mockUserRepo := ts.Mocks.Users.(*mocks.UserRepositoryMock)
	mockUserRepo.On("Get", mock.Anything, userID).Return(&models.User{ID: userID, Email: "validation@example.com"}, nil)

	t.Run("missing required fields", func(t *testing.T) {
		w := ts.createAuthenticatedRequest("POST", "/api/v1/events", token, map[string]string{})
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("name too short", func(t *testing.T) {
		event := models.Event{
			Name:        "A",
			Description: "Valid description",
			Date:        "2025-12-31",
			Location:    "Valid location",
		}
		w := ts.createAuthenticatedRequest("POST", "/api/v1/events", token, event)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("description too short", func(t *testing.T) {
		event := models.Event{
			Name:        "Valid Name",
			Description: "Short",
			Date:        "2025-12-31",
			Location:    "Valid location",
		}
		w := ts.createAuthenticatedRequest("POST", "/api/v1/events", token, event)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid date format", func(t *testing.T) {
		event := models.Event{
			Name:        "Valid Name",
			Description: "Valid description that is long enough",
			Date:        "invalid-date",
			Location:    "Valid location",
		}
		w := ts.createAuthenticatedRequest("POST", "/api/v1/events", token, event)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("location too short", func(t *testing.T) {
		event := models.Event{
			Name:        "Valid Name",
			Description: "Valid description that is long enough",
			Date:        "2025-12-31",
			Location:    "A",
		}
		w := ts.createAuthenticatedRequest("POST", "/api/v1/events", token, event)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
