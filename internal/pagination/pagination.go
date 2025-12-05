package pagination

// PaginationRequest holds pagination parameters
type PaginationRequest struct {
	Page     int `json:"page" form:"page" binding:"min=1"`
	PageSize int `json:"page_size" form:"page_size" binding:"min=1,max=100"`
}

// PaginationResponse holds pagination metadata
type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// NewPaginationRequest creates a new pagination request with defaults
func NewPaginationRequest() *PaginationRequest {
	return &PaginationRequest{
		Page:     1,
		PageSize: 20,
	}
}

// Offset calculates the database offset
func (p *PaginationRequest) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// NewPaginationResponse creates pagination metadata
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

// PaginatedResult wraps paginated data with metadata
type PaginatedResult struct {
	Data       interface{}         `json:"data"`
	Pagination *PaginationResponse `json:"pagination"`
}
