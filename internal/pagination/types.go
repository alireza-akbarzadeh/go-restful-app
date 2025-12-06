package pagination

import "time"

// ===========================================
// PAGINATION TYPES & CONSTANTS
// ===========================================

// PaginationType defines the type of pagination
type PaginationType string

const (
	// OffsetPagination uses traditional page/offset based pagination
	OffsetPagination PaginationType = "offset"
	// CursorPagination uses cursor-based pagination (ideal for infinite scroll)
	CursorPagination PaginationType = "cursor"
)

// SortDirection defines the sort order
type SortDirection string

const (
	SortAsc  SortDirection = "asc"
	SortDesc SortDirection = "desc"
)

// FilterOperator defines the filter operation
type FilterOperator string

const (
	OpEqual        FilterOperator = "eq"      // Equal: field = value
	OpNotEqual     FilterOperator = "neq"     // Not Equal: field != value
	OpGreaterThan  FilterOperator = "gt"      // Greater Than: field > value
	OpGreaterEqual FilterOperator = "gte"     // Greater or Equal: field >= value
	OpLessThan     FilterOperator = "lt"      // Less Than: field < value
	OpLessEqual    FilterOperator = "lte"     // Less or Equal: field <= value
	OpLike         FilterOperator = "like"    // Like: field LIKE %value%
	OpILike        FilterOperator = "ilike"   // Case-insensitive Like
	OpIn           FilterOperator = "in"      // In: field IN (values)
	OpNotIn        FilterOperator = "nin"     // Not In: field NOT IN (values)
	OpIsNull       FilterOperator = "null"    // Is Null: field IS NULL
	OpIsNotNull    FilterOperator = "notnull" // Is Not Null: field IS NOT NULL
	OpBetween      FilterOperator = "between" // Between: field BETWEEN value1 AND value2
)

// ===========================================
// SORT & FILTER TYPES
// ===========================================

// SortField represents a single sort field with direction
type SortField struct {
	Field     string        `json:"field"`
	Direction SortDirection `json:"direction"`
}

// Filter represents a single filter condition
type Filter struct {
	Field    string         `json:"field"`
	Operator FilterOperator `json:"operator"`
	Value    interface{}    `json:"value"`
	Values   []interface{}  `json:"values,omitempty"` // For IN, NOT IN, BETWEEN operators
}

// CursorData holds the cursor information for cursor-based pagination
type CursorData struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	SortValue string    `json:"sort_value,omitempty"`
}
