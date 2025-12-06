package pagination

import (
	"math"
)

// ===========================================
// BASIC PAGINATION RESPONSE
// ===========================================

// PaginationResponse holds basic pagination metadata
type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// NewPaginationResponse creates pagination metadata from page info and total count
func NewPaginationResponse(page, pageSize int, totalItems int64) *PaginationResponse {
	totalPages := int((totalItems + int64(pageSize) - 1) / int64(pageSize))

	return &PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// PaginatedResult wraps paginated data with basic metadata
type PaginatedResult struct {
	Data       interface{}         `json:"data"`
	Pagination *PaginationResponse `json:"pagination"`
}

// ===========================================
// ADVANCED PAGINATION RESPONSE
// ===========================================

// AdvancedPaginationResponse holds enhanced pagination metadata
type AdvancedPaginationResponse struct {
	// Offset-based pagination info
	Page       int   `json:"page,omitempty"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`

	// Cursor-based pagination info
	StartCursor string `json:"start_cursor,omitempty"`
	EndCursor   string `json:"end_cursor,omitempty"`
	HasNextPage bool   `json:"has_next_page"`
	HasPrevPage bool   `json:"has_prev_page"`

	// Current result count
	Count int `json:"count"`

	// HATEOAS navigation links
	Links PaginationLinks `json:"links"`

	// Applied filters & sorts for transparency
	AppliedFilters []Filter    `json:"applied_filters,omitempty"`
	AppliedSort    []SortField `json:"applied_sort,omitempty"`
}

// PaginationLinks holds HATEOAS navigation links
type PaginationLinks struct {
	Self  string `json:"self"`
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
}

// AdvancedPaginatedResult wraps paginated data with enhanced metadata
type AdvancedPaginatedResult struct {
	Success    bool                        `json:"success"`
	Data       interface{}                 `json:"data"`
	Pagination *AdvancedPaginationResponse `json:"pagination"`
	Meta       map[string]interface{}      `json:"meta,omitempty"`
}

// ===========================================
// RESPONSE BUILDER
// ===========================================

// ResponseBuilder helps construct paginated responses
type ResponseBuilder struct {
	data    interface{}
	request *AdvancedPaginationRequest
	total   int64
	count   int
	firstID int
	lastID  int
	meta    map[string]interface{}
}

// NewResponseBuilder creates a new response builder
func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{
		meta: make(map[string]interface{}),
	}
}

// WithData sets the data payload
func (rb *ResponseBuilder) WithData(data interface{}) *ResponseBuilder {
	rb.data = data
	return rb
}

// WithRequest sets the original pagination request
func (rb *ResponseBuilder) WithRequest(req *AdvancedPaginationRequest) *ResponseBuilder {
	rb.request = req
	return rb
}

// WithTotal sets the total count
func (rb *ResponseBuilder) WithTotal(total int64) *ResponseBuilder {
	rb.total = total
	return rb
}

// WithCount sets the current page count
func (rb *ResponseBuilder) WithCount(count int) *ResponseBuilder {
	rb.count = count
	return rb
}

// WithCursorIDs sets the first and last IDs for cursor pagination
func (rb *ResponseBuilder) WithCursorIDs(firstID, lastID int) *ResponseBuilder {
	rb.firstID = firstID
	rb.lastID = lastID
	return rb
}

// WithMeta adds metadata to the response
func (rb *ResponseBuilder) WithMeta(key string, value interface{}) *ResponseBuilder {
	rb.meta[key] = value
	return rb
}

// Build constructs the final paginated response
func (rb *ResponseBuilder) Build() *AdvancedPaginatedResult {
	if rb.request == nil {
		rb.request = NewAdvancedPaginationRequest()
	}

	resp := &AdvancedPaginationResponse{
		PageSize: rb.request.PageSize,
		Count:    rb.count,
	}

	if rb.request.Type == CursorPagination {
		rb.buildCursorResponse(resp)
	} else {
		rb.buildOffsetResponse(resp)
	}

	// Build HATEOAS links
	resp.Links = buildLinks(rb.request, resp)

	// Include applied filters/sorts for transparency
	resp.AppliedFilters = rb.request.Filters
	resp.AppliedSort = rb.request.Sort

	result := &AdvancedPaginatedResult{
		Success:    true,
		Data:       rb.data,
		Pagination: resp,
	}

	if len(rb.meta) > 0 {
		result.Meta = rb.meta
	}

	return result
}

func (rb *ResponseBuilder) buildCursorResponse(resp *AdvancedPaginationResponse) {
	if rb.firstID > 0 {
		resp.StartCursor = EncodeCursor(&CursorData{ID: rb.firstID})
	}
	if rb.lastID > 0 {
		resp.EndCursor = EncodeCursor(&CursorData{ID: rb.lastID})
	}
	resp.HasNextPage = rb.count >= rb.request.PageSize
	resp.HasPrevPage = rb.request.Cursor != "" || rb.request.After != ""
}

func (rb *ResponseBuilder) buildOffsetResponse(resp *AdvancedPaginationResponse) {
	resp.Page = rb.request.Page
	resp.TotalItems = rb.total
	resp.TotalPages = int(math.Ceil(float64(rb.total) / float64(rb.request.PageSize)))
	resp.HasNextPage = rb.request.Page < resp.TotalPages
	resp.HasPrevPage = rb.request.Page > 1
}

// BuildResponse is a convenience function to build a response quickly
func BuildResponse(
	data interface{},
	req *AdvancedPaginationRequest,
	total int64,
	count int,
	firstID, lastID int,
) *AdvancedPaginatedResult {
	return NewResponseBuilder().
		WithData(data).
		WithRequest(req).
		WithTotal(total).
		WithCount(count).
		WithCursorIDs(firstID, lastID).
		Build()
}
