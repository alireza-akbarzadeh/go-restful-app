package pagination

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ===========================================
// REQUEST PARSER
// ===========================================

// ParseFromContext parses pagination request from gin context
func ParseFromContext(c *gin.Context) *AdvancedPaginationRequest {
	req := NewAdvancedPaginationRequest()

	// Parse pagination type
	parsePaginationType(c, req)

	// Parse pagination params
	parseOffsetParams(c, req)
	parseCursorParams(c, req)

	// Parse sorting
	parseSorting(c, req)

	// Parse filtering
	parseFiltering(c, req)

	// Parse search
	parseSearch(c, req)

	// Parse options
	parseOptions(c, req)

	// Set base URL for HATEOAS links
	req.BaseURL = c.Request.URL.Path

	return req
}

func parsePaginationType(c *gin.Context, req *AdvancedPaginationRequest) {
	if pType := c.Query("type"); pType == "cursor" {
		req.Type = CursorPagination
	}
}

func parseOffsetParams(c *gin.Context, req *AdvancedPaginationRequest) {
	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			req.Page = p
		}
	}

	if pageSize := c.Query("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 && ps <= 100 {
			req.PageSize = ps
		}
	}
}

func parseCursorParams(c *gin.Context, req *AdvancedPaginationRequest) {
	req.Cursor = c.Query("cursor")
	req.After = c.Query("after")
	req.Before = c.Query("before")

	if first := c.Query("first"); first != "" {
		if f, err := strconv.Atoi(first); err == nil && f > 0 {
			req.First = f
			req.PageSize = f
		}
	}

	if last := c.Query("last"); last != "" {
		if l, err := strconv.Atoi(last); err == nil && l > 0 {
			req.Last = l
			req.PageSize = l
		}
	}
}

func parseSorting(c *gin.Context, req *AdvancedPaginationRequest) {
	if sortRaw := c.Query("sort"); sortRaw != "" {
		req.SortRaw = sortRaw
		req.Sort = ParseSortString(sortRaw)
	}
}

func parseFiltering(c *gin.Context, req *AdvancedPaginationRequest) {
	// Parse filter query param (JSON format)
	if filterRaw := c.Query("filter"); filterRaw != "" {
		req.FilterRaw = filterRaw
		req.Filters = ParseFilterString(filterRaw)
	}

	// Parse individual filter params (field[operator]=value format)
	for key, values := range c.Request.URL.Query() {
		if strings.Contains(key, "[") && strings.Contains(key, "]") {
			filter := parseFilterParam(key, values)
			if filter != nil {
				req.Filters = append(req.Filters, *filter)
			}
		}
	}
}

func parseSearch(c *gin.Context, req *AdvancedPaginationRequest) {
	req.Search = c.Query("search")
	req.SearchFields = c.QueryArray("search_fields")
}

func parseOptions(c *gin.Context, req *AdvancedPaginationRequest) {
	if includeTotal := c.Query("include_total"); includeTotal != "" {
		req.IncludeTotal = includeTotal == "true" || includeTotal == "1"
	}
}

// ===========================================
// SORT STRING PARSER
// ===========================================

// ParseSortString parses a sort string in format:
// - "field:direction,field2:direction"
// - "-field,field2" (- prefix means descending)
func ParseSortString(sortStr string) []SortField {
	var sorts []SortField

	parts := strings.Split(sortStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		field, direction := parseSortPart(part)
		if field != "" {
			sorts = append(sorts, SortField{Field: field, Direction: direction})
		}
	}

	return sorts
}

func parseSortPart(part string) (string, SortDirection) {
	// Check for - prefix (descending)
	if strings.HasPrefix(part, "-") {
		return part[1:], SortDesc
	}

	// Check for :direction suffix
	if strings.Contains(part, ":") {
		fieldParts := strings.SplitN(part, ":", 2)
		field := fieldParts[0]
		direction := SortAsc

		if len(fieldParts) > 1 && strings.ToLower(fieldParts[1]) == "desc" {
			direction = SortDesc
		}

		return field, direction
	}

	// Default to ascending
	return part, SortAsc
}

// ===========================================
// FILTER STRING PARSER
// ===========================================

// ParseFilterString parses a filter string:
// - JSON format: [{"field":"name","operator":"eq","value":"test"}]
// - Simple format: "key:value,key2:value2"
func ParseFilterString(filterStr string) []Filter {
	var filters []Filter

	// Try to parse as JSON array first
	if err := json.Unmarshal([]byte(filterStr), &filters); err == nil {
		return filters
	}

	// Fall back to simple key:value format
	return parseSimpleFilterString(filterStr)
}

func parseSimpleFilterString(filterStr string) []Filter {
	var filters []Filter

	parts := strings.Split(filterStr, ",")
	for _, part := range parts {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) == 2 {
			filters = append(filters, Filter{
				Field:    strings.TrimSpace(kv[0]),
				Operator: OpEqual,
				Value:    strings.TrimSpace(kv[1]),
			})
		}
	}

	return filters
}

// parseFilterParam parses a filter parameter in format field[operator]=value
// Examples: name[eq]=test, price[gte]=100, status[in]=active&status[in]=pending
func parseFilterParam(key string, values []string) *Filter {
	start := strings.Index(key, "[")
	end := strings.Index(key, "]")

	if start == -1 || end == -1 || start >= end {
		return nil
	}

	field := key[:start]
	operator := FilterOperator(key[start+1 : end])

	if len(values) == 0 {
		return nil
	}

	filter := &Filter{
		Field:    field,
		Operator: operator,
	}

	// Handle multi-value operators
	switch operator {
	case OpIn, OpNotIn:
		filter.Values = stringsToInterfaces(values)
	case OpBetween:
		if len(values) >= 2 {
			filter.Values = []interface{}{values[0], values[1]}
		}
	default:
		filter.Value = values[0]
	}

	return filter
}

func stringsToInterfaces(strs []string) []interface{} {
	result := make([]interface{}, len(strs))
	for i, s := range strs {
		result[i] = s
	}
	return result
}
