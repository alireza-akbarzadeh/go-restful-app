package repository

import (
	"context"
	"errors"

	appErrors "github.com/alireza-akbarzadeh/ginflow/internal/errors"
	"github.com/alireza-akbarzadeh/ginflow/internal/logging"
	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/alireza-akbarzadeh/ginflow/internal/query"
	"gorm.io/gorm"
)

// EventRepository handles event database operations
type EventRepository struct {
	DB *gorm.DB
}

// NewEventRepository creates a new EventRepository
func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{DB: db}
}

// Insert creates a new event in the database
func (r *EventRepository) Insert(ctx context.Context, event *models.Event) (*models.Event, error) {
	logging.Debug(ctx, "creating new event", "name", event.Name, "owner_id", event.OwnerID)

	result := r.DB.WithContext(ctx).Create(event)
	if result.Error != nil {
		logging.Error(ctx, "failed to create event", result.Error, "name", event.Name)
		return nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to create event")
	}

	logging.Info(ctx, "event created successfully", "event_id", event.ID, "name", event.Name)
	return event, nil
}

// Get retrieves an event by ID
func (r *EventRepository) Get(ctx context.Context, id int) (*models.Event, error) {
	logging.Debug(ctx, "retrieving event by ID", "event_id", id)

	var event models.Event
	result := r.DB.WithContext(ctx).First(&event, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logging.Debug(ctx, "event not found", "event_id", id)
			return nil, appErrors.Newf(appErrors.ErrNotFound, "event with ID %d not found", id)
		}
		logging.Error(ctx, "failed to retrieve event", result.Error, "event_id", id)
		return nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to retrieve event")
	}

	logging.Debug(ctx, "event retrieved successfully", "event_id", id, "name", event.Name)
	return &event, nil
}

// GetAll retrieves all events
func (r *EventRepository) GetAll(ctx context.Context) ([]*models.Event, error) {
	logging.Debug(ctx, "retrieving all events")

	var events []*models.Event
	result := r.DB.WithContext(ctx).Find(&events)
	if result.Error != nil {
		logging.Error(ctx, "failed to retrieve all events", result.Error)
		return nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to retrieve events")
	}

	logging.Info(ctx, "events retrieved successfully", "count", len(events))
	return events, nil
}

// Update updates an existing event
func (r *EventRepository) Update(ctx context.Context, event *models.Event) error {
	logging.Debug(ctx, "updating event", "event_id", event.ID, "name", event.Name)

	result := r.DB.WithContext(ctx).Save(event)
	if result.Error != nil {
		logging.Error(ctx, "failed to update event", result.Error, "event_id", event.ID)
		return appErrors.New(appErrors.ErrDatabaseOperation, "failed to update event")
	}

	logging.Info(ctx, "event updated successfully", "event_id", event.ID, "name", event.Name)
	return nil
}

// Delete removes an event by ID
func (r *EventRepository) Delete(ctx context.Context, id int) error {
	logging.Debug(ctx, "deleting event", "event_id", id)

	result := r.DB.WithContext(ctx).Delete(&models.Event{}, id)
	if result.Error != nil {
		logging.Error(ctx, "failed to delete event", result.Error, "event_id", id)
		return appErrors.New(appErrors.ErrDatabaseOperation, "failed to delete event")
	}

	if result.RowsAffected == 0 {
		logging.Debug(ctx, "no event found to delete", "event_id", id)
		return appErrors.Newf(appErrors.ErrNotFound, "event with ID %d not found", id)
	}

	logging.Info(ctx, "event deleted successfully", "event_id", id)
	return nil
}

// ListWithPagination retrieves events with pagination
func (r *EventRepository) ListWithPagination(ctx context.Context, req *query.PaginationRequest) ([]*models.Event, *query.PaginationResponse, error) {
	logging.Debug(ctx, "retrieving events with pagination", "page", req.Page, "page_size", req.PageSize)

	var events []*models.Event
	var total int64

	// Count total records
	if err := r.DB.WithContext(ctx).Model(&models.Event{}).Count(&total).Error; err != nil {
		logging.Error(ctx, "failed to count events", err)
		return nil, nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to count events")
	}

	// Get paginated records with owner preloaded
	if err := r.DB.WithContext(ctx).
		Preload("Owner").
		Offset(req.Offset()).
		Limit(req.PageSize).
		Find(&events).Error; err != nil {
		logging.Error(ctx, "failed to retrieve events", err)
		return nil, nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to retrieve events")
	}

	// Calculate pagination response
	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	paginationResp := &query.PaginationResponse{
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalItems: total,
		TotalPages: totalPages,
		HasNext:    req.Page < totalPages,
		HasPrev:    req.Page > 1,
	}

	logging.Info(ctx, "events retrieved successfully", "count", len(events), "total", total, "page", req.Page)
	return events, paginationResp, nil
}

// ListWithAdvancedPagination retrieves events with advanced pagination, filtering, sorting, and search
func (r *EventRepository) ListWithAdvancedPagination(ctx context.Context, req *query.QueryParams) ([]*models.Event, *query.PaginatedList, error) {
	logging.Debug(ctx, "retrieving events with advanced pagination",
		"page", req.Page,
		"page_size", req.PageSize,
		"type", req.Type,
		"search", req.Search,
	)

	var events []*models.Event
	var total int64

	// Build pagination query
	builder := query.NewQueryBuilder(r.DB.WithContext(ctx).Model(&models.Event{})).
		WithRequest(req).
		AllowFilters("name", "location", "owner_id", "start_date", "end_date", "created_at", "status").
		AllowSorts("name", "start_date", "end_date", "created_at", "updated_at").
		SearchColumns("name", "description", "location").
		DefaultSort("created_at", query.SortDesc)

	// Get count if needed
	if req.IncludeTotal {
		countQuery := r.DB.WithContext(ctx).Model(&models.Event{})
		// Apply filters and search for count
		for _, filter := range req.Filters {
			countQuery = query.FilterBy(filter)(countQuery)
		}
		if req.Search != "" {
			countQuery = query.Search(req.Search, "name", "description", "location")(countQuery)
		}
		countQuery.Count(&total)
	}

	// Execute main query
	dbQuery := builder.Build()
	if err := dbQuery.Preload("Owner").Find(&events).Error; err != nil {
		logging.Error(ctx, "failed to retrieve events with advanced pagination", err)
		return nil, nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to retrieve events")
	}

	// Get first and last IDs for cursor pagination
	var firstID, lastID int
	if len(events) > 0 {
		firstID = events[0].ID
		lastID = events[len(events)-1].ID
	}

	// Build response
	result := query.BuildResponse(events, req, total, len(events), firstID, lastID)

	logging.Info(ctx, "events retrieved with advanced pagination",
		"count", len(events),
		"total", total,
		"page", req.Page,
	)

	return events, result, nil
}

// GetByOwnerID retrieves events by owner ID
func (r *EventRepository) GetByOwnerID(ctx context.Context, ownerID int) ([]*models.Event, error) {
	logging.Debug(ctx, "retrieving events by owner ID", "owner_id", ownerID)

	var events []*models.Event
	result := r.DB.WithContext(ctx).Where("owner_id = ?", ownerID).Find(&events)
	if result.Error != nil {
		logging.Error(ctx, "failed to retrieve events by owner ID", result.Error, "owner_id", ownerID)
		return nil, appErrors.New(appErrors.ErrDatabaseOperation, "failed to retrieve events by owner")
	}

	logging.Info(ctx, "events retrieved by owner ID", "count", len(events), "owner_id", ownerID)
	return events, nil
}
