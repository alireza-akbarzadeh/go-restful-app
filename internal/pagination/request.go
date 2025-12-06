package pagination

// ===========================================
// BASIC PAGINATION REQUEST (Backwards Compatible)
// ===========================================

// PaginationRequest holds basic pagination parameters
type PaginationRequest struct {
	Page     int `json:"page" form:"page" binding:"min=1"`
	PageSize int `json:"page_size" form:"page_size" binding:"min=1,max=100"`
}

// NewPaginationRequest creates a new pagination request with defaults
func NewPaginationRequest() *PaginationRequest {
	return &PaginationRequest{
		Page:     1,
		PageSize: 20,
	}
}

// Offset calculates the database offset for the current page
func (p *PaginationRequest) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Validate ensures pagination values are within acceptable ranges
func (p *PaginationRequest) Validate() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

// ===========================================
// ADVANCED PAGINATION REQUEST
// ===========================================

// AdvancedPaginationRequest holds all pagination, filtering, and sorting parameters
type AdvancedPaginationRequest struct {
	// Pagination type (offset or cursor)
	Type PaginationType `json:"type" form:"type"`

	// Offset-based pagination
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`

	// Cursor-based pagination
	Cursor string `json:"cursor" form:"cursor"`
	After  string `json:"after" form:"after"`
	Before string `json:"before" form:"before"`
	First  int    `json:"first" form:"first"`
	Last   int    `json:"last" form:"last"`

	// Sorting - parsed from query string
	Sort    []SortField `json:"sort"`
	SortRaw string      `json:"-" form:"sort"` // Raw: "field:direction,field2:direction" or "-field"

	// Filtering - parsed from query string
	Filters   []Filter `json:"filters"`
	FilterRaw string   `json:"-" form:"filter"` // Raw filter string

	// Search
	Search       string   `json:"search" form:"search"`
	SearchFields []string `json:"search_fields" form:"search_fields"`

	// Performance: Include total count (can be disabled for large datasets)
	IncludeTotal bool `json:"include_total" form:"include_total"`

	// Base URL for HATEOAS links
	BaseURL string `json:"-"`
}

// NewAdvancedPaginationRequest creates a new advanced pagination request with defaults
func NewAdvancedPaginationRequest() *AdvancedPaginationRequest {
	return &AdvancedPaginationRequest{
		Type:         OffsetPagination,
		Page:         1,
		PageSize:     20,
		IncludeTotal: true,
	}
}

// Validate ensures pagination values are within acceptable ranges
func (r *AdvancedPaginationRequest) Validate() {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.PageSize <= 0 {
		r.PageSize = 20
	}
	if r.PageSize > 100 {
		r.PageSize = 100
	}
}

// Offset calculates the database offset for the current page
func (r *AdvancedPaginationRequest) Offset() int {
	return (r.Page - 1) * r.PageSize
}

// IsCursorBased returns true if using cursor-based pagination
func (r *AdvancedPaginationRequest) IsCursorBased() bool {
	return r.Type == CursorPagination
}

// HasFilters returns true if any filters are applied
func (r *AdvancedPaginationRequest) HasFilters() bool {
	return len(r.Filters) > 0
}

// HasSort returns true if any sorting is applied
func (r *AdvancedPaginationRequest) HasSort() bool {
	return len(r.Sort) > 0
}

// HasSearch returns true if a search term is provided
func (r *AdvancedPaginationRequest) HasSearch() bool {
	return r.Search != ""
}
