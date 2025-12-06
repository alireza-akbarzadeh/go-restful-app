package pagination

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// ===========================================
// PAGINATION BUILDER
// ===========================================

// PaginationBuilder helps build pagination queries with a fluent API
type PaginationBuilder struct {
	request        *AdvancedPaginationRequest
	db             *gorm.DB
	allowedFilters map[string]bool
	allowedSorts   map[string]bool
	defaultSort    []SortField
	searchColumns  []string
}

// NewPaginationBuilder creates a new pagination builder
func NewPaginationBuilder(db *gorm.DB) *PaginationBuilder {
	return &PaginationBuilder{
		db:             db,
		allowedFilters: make(map[string]bool),
		allowedSorts:   make(map[string]bool),
		defaultSort:    []SortField{{Field: "created_at", Direction: SortDesc}},
	}
}

// ===========================================
// BUILDER CONFIGURATION METHODS
// ===========================================

// WithRequest sets the pagination request
func (pb *PaginationBuilder) WithRequest(req *AdvancedPaginationRequest) *PaginationBuilder {
	pb.request = req
	return pb
}

// AllowFilters sets allowed filter fields (for security)
func (pb *PaginationBuilder) AllowFilters(fields ...string) *PaginationBuilder {
	for _, f := range fields {
		pb.allowedFilters[f] = true
	}
	return pb
}

// AllowSorts sets allowed sort fields (for security)
func (pb *PaginationBuilder) AllowSorts(fields ...string) *PaginationBuilder {
	for _, f := range fields {
		pb.allowedSorts[f] = true
	}
	return pb
}

// DefaultSort sets the default sort order when no sort is specified
func (pb *PaginationBuilder) DefaultSort(field string, direction SortDirection) *PaginationBuilder {
	pb.defaultSort = []SortField{{Field: field, Direction: direction}}
	return pb
}

// SearchColumns sets the columns to search in
func (pb *PaginationBuilder) SearchColumns(columns ...string) *PaginationBuilder {
	pb.searchColumns = columns
	return pb
}

// ===========================================
// BUILD METHODS
// ===========================================

// Build applies pagination, filtering, and sorting to the query
func (pb *PaginationBuilder) Build() *gorm.DB {
	pb.ensureRequest()

	query := pb.db
	query = pb.applyFilters(query)
	query = pb.applySearch(query)
	query = pb.applySorting(query)
	query = pb.applyPagination(query)

	return query
}

// BuildWithCount applies pagination and returns both query and total count
func (pb *PaginationBuilder) BuildWithCount(model interface{}) (*gorm.DB, int64) {
	pb.ensureRequest()

	var total int64

	// Count query (without pagination)
	countQuery := pb.db.Model(model)
	countQuery = pb.applyFilters(countQuery)
	countQuery = pb.applySearch(countQuery)
	countQuery.Count(&total)

	// Main query with pagination
	query := pb.Build()

	return query, total
}

func (pb *PaginationBuilder) ensureRequest() {
	if pb.request == nil {
		pb.request = NewAdvancedPaginationRequest()
	}
}

// ===========================================
// FILTER APPLICATION
// ===========================================

func (pb *PaginationBuilder) applyFilters(query *gorm.DB) *gorm.DB {
	for _, filter := range pb.request.Filters {
		if pb.isFilterAllowed(filter.Field) {
			query = ApplyFilter(query, filter)
		}
	}
	return query
}

func (pb *PaginationBuilder) isFilterAllowed(field string) bool {
	// If no allowed filters specified, allow all
	if len(pb.allowedFilters) == 0 {
		return true
	}
	return pb.allowedFilters[field]
}

// ===========================================
// SEARCH APPLICATION
// ===========================================

func (pb *PaginationBuilder) applySearch(query *gorm.DB) *gorm.DB {
	if pb.request.Search == "" {
		return query
	}

	searchFields := pb.getSearchFields()
	if len(searchFields) == 0 {
		return query
	}

	return ApplySearch(query, pb.request.Search, searchFields...)
}

func (pb *PaginationBuilder) getSearchFields() []string {
	// Use request search fields if provided, otherwise use builder's
	if len(pb.request.SearchFields) > 0 {
		return pb.request.SearchFields
	}
	return pb.searchColumns
}

// ===========================================
// SORT APPLICATION
// ===========================================

func (pb *PaginationBuilder) applySorting(query *gorm.DB) *gorm.DB {
	sorts := pb.request.Sort
	if len(sorts) == 0 {
		sorts = pb.defaultSort
	}

	for _, sort := range sorts {
		if pb.isSortAllowed(sort.Field) {
			query = ApplySortField(query, sort)
		}
	}

	return query
}

func (pb *PaginationBuilder) isSortAllowed(field string) bool {
	// If no allowed sorts specified, allow all
	if len(pb.allowedSorts) == 0 {
		return true
	}
	return pb.allowedSorts[field]
}

// ===========================================
// PAGINATION APPLICATION
// ===========================================

func (pb *PaginationBuilder) applyPagination(query *gorm.DB) *gorm.DB {
	limit := pb.getLimit()

	if pb.request.Type == CursorPagination {
		return pb.applyCursorPagination(query, limit)
	}

	return pb.applyOffsetPagination(query, limit)
}

func (pb *PaginationBuilder) getLimit() int {
	limit := pb.request.PageSize
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return limit
}

func (pb *PaginationBuilder) applyCursorPagination(query *gorm.DB, limit int) *gorm.DB {
	if pb.request.Cursor != "" {
		cursor, err := DecodeCursor(pb.request.Cursor)
		if err == nil {
			query = query.Where("id > ?", cursor.ID)
		}
	}
	return query.Limit(limit)
}

func (pb *PaginationBuilder) applyOffsetPagination(query *gorm.DB, limit int) *gorm.DB {
	page := pb.request.Page
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	return query.Offset(offset).Limit(limit)
}

// ===========================================
// STANDALONE FILTER/SEARCH/SORT FUNCTIONS
// ===========================================

// ApplyFilter applies a single filter to a GORM query
func ApplyFilter(query *gorm.DB, filter Filter) *gorm.DB {
	switch filter.Operator {
	case OpEqual:
		return query.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
	case OpNotEqual:
		return query.Where(fmt.Sprintf("%s != ?", filter.Field), filter.Value)
	case OpGreaterThan:
		return query.Where(fmt.Sprintf("%s > ?", filter.Field), filter.Value)
	case OpGreaterEqual:
		return query.Where(fmt.Sprintf("%s >= ?", filter.Field), filter.Value)
	case OpLessThan:
		return query.Where(fmt.Sprintf("%s < ?", filter.Field), filter.Value)
	case OpLessEqual:
		return query.Where(fmt.Sprintf("%s <= ?", filter.Field), filter.Value)
	case OpLike:
		return query.Where(fmt.Sprintf("%s LIKE ?", filter.Field), "%"+fmt.Sprint(filter.Value)+"%")
	case OpILike:
		return query.Where(fmt.Sprintf("LOWER(%s) LIKE LOWER(?)", filter.Field), "%"+fmt.Sprint(filter.Value)+"%")
	case OpIn:
		return query.Where(fmt.Sprintf("%s IN ?", filter.Field), filter.Values)
	case OpNotIn:
		return query.Where(fmt.Sprintf("%s NOT IN ?", filter.Field), filter.Values)
	case OpIsNull:
		return query.Where(fmt.Sprintf("%s IS NULL", filter.Field))
	case OpIsNotNull:
		return query.Where(fmt.Sprintf("%s IS NOT NULL", filter.Field))
	case OpBetween:
		if len(filter.Values) >= 2 {
			return query.Where(fmt.Sprintf("%s BETWEEN ? AND ?", filter.Field), filter.Values[0], filter.Values[1])
		}
	}
	return query
}

// ApplySearch applies search to a GORM query across multiple columns
func ApplySearch(query *gorm.DB, term string, columns ...string) *gorm.DB {
	if term == "" || len(columns) == 0 {
		return query
	}

	searchTerm := "%" + strings.ToLower(term) + "%"
	var conditions []string
	var args []interface{}

	for _, col := range columns {
		conditions = append(conditions, fmt.Sprintf("LOWER(%s) LIKE ?", col))
		args = append(args, searchTerm)
	}

	return query.Where(strings.Join(conditions, " OR "), args...)
}

// ApplySortField applies a single sort field to a GORM query
func ApplySortField(query *gorm.DB, sort SortField) *gorm.DB {
	direction := "ASC"
	if sort.Direction == SortDesc {
		direction = "DESC"
	}
	return query.Order(fmt.Sprintf("%s %s", sort.Field, direction))
}
