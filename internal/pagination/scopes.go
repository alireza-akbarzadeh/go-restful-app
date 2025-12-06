package pagination

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// ===========================================
// GORM SCOPES
// These are reusable GORM scope functions for common operations
// ===========================================

// Paginate is a GORM scope for simple offset pagination
// Usage: db.Scopes(pagination.Paginate(page, pageSize)).Find(&users)
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Validate and set defaults
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 20
		}
		if pageSize > 100 {
			pageSize = 100
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// Search is a GORM scope for case-insensitive search across multiple columns
// Usage: db.Scopes(pagination.Search("john", "name", "email")).Find(&users)
func Search(term string, columns ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if term == "" || len(columns) == 0 {
			return db
		}

		searchTerm := "%" + strings.ToLower(term) + "%"
		var conditions []string
		var args []interface{}

		for _, col := range columns {
			conditions = append(conditions, fmt.Sprintf("LOWER(%s) LIKE ?", col))
			args = append(args, searchTerm)
		}

		return db.Where(strings.Join(conditions, " OR "), args...)
	}
}

// SortBy is a GORM scope for sorting by a single field
// Usage: db.Scopes(pagination.SortBy("created_at", pagination.SortDesc)).Find(&users)
func SortBy(field string, direction SortDirection) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if field == "" {
			return db
		}

		dir := "ASC"
		if direction == SortDesc {
			dir = "DESC"
		}

		return db.Order(fmt.Sprintf("%s %s", field, dir))
	}
}

// FilterBy is a GORM scope for applying multiple filters
// Usage: db.Scopes(pagination.FilterBy(filters...)).Find(&users)
func FilterBy(filters ...Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, filter := range filters {
			db = ApplyFilter(db, filter)
		}
		return db
	}
}

// WhereEqual is a convenience scope for simple equality filter
// Usage: db.Scopes(pagination.WhereEqual("status", "active")).Find(&users)
func WhereEqual(field string, value interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s = ?", field), value)
	}
}

// WhereIn is a convenience scope for IN filter
// Usage: db.Scopes(pagination.WhereIn("status", []string{"active", "pending"})).Find(&users)
func WhereIn(field string, values interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s IN ?", field), values)
	}
}

// WhereBetween is a convenience scope for BETWEEN filter
// Usage: db.Scopes(pagination.WhereBetween("price", 100, 500)).Find(&products)
func WhereBetween(field string, min, max interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", field), min, max)
	}
}

// WhereNotNull is a convenience scope for NOT NULL filter
// Usage: db.Scopes(pagination.WhereNotNull("deleted_at")).Find(&users)
func WhereNotNull(field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s IS NOT NULL", field))
	}
}

// WhereNull is a convenience scope for NULL filter
// Usage: db.Scopes(pagination.WhereNull("deleted_at")).Find(&users)
func WhereNull(field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s IS NULL", field))
	}
}

// OrderByCreatedAt is a convenience scope for ordering by created_at
// Usage: db.Scopes(pagination.OrderByCreatedAt(pagination.SortDesc)).Find(&users)
func OrderByCreatedAt(direction SortDirection) func(db *gorm.DB) *gorm.DB {
	return SortBy("created_at", direction)
}

// OrderByUpdatedAt is a convenience scope for ordering by updated_at
// Usage: db.Scopes(pagination.OrderByUpdatedAt(pagination.SortDesc)).Find(&users)
func OrderByUpdatedAt(direction SortDirection) func(db *gorm.DB) *gorm.DB {
	return SortBy("updated_at", direction)
}

// Limit is a simple GORM scope for limiting results
// Usage: db.Scopes(pagination.Limit(10)).Find(&users)
func Limit(limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if limit <= 0 {
			limit = 20
		}
		if limit > 100 {
			limit = 100
		}
		return db.Limit(limit)
	}
}
