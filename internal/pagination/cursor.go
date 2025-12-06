package pagination

import (
	"encoding/base64"
	"encoding/json"
)

// ===========================================
// CURSOR ENCODING/DECODING
// ===========================================

// EncodeCursor encodes cursor data to a base64 string
func EncodeCursor(data *CursorData) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(jsonData)
}

// DecodeCursor decodes a base64 cursor string back to CursorData
func DecodeCursor(cursor string) (*CursorData, error) {
	data, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	var cursorData CursorData
	if err := json.Unmarshal(data, &cursorData); err != nil {
		return nil, err
	}

	return &cursorData, nil
}

// ===========================================
// CURSOR HELPERS
// ===========================================

// CreateCursor creates a cursor from an ID
func CreateCursor(id int) string {
	return EncodeCursor(&CursorData{ID: id})
}

// ExtractCursorID extracts the ID from a cursor string
// Returns 0 if cursor is invalid
func ExtractCursorID(cursor string) int {
	if cursor == "" {
		return 0
	}

	data, err := DecodeCursor(cursor)
	if err != nil {
		return 0
	}

	return data.ID
}

// ValidateCursor checks if a cursor is valid
func ValidateCursor(cursor string) bool {
	if cursor == "" {
		return false
	}

	_, err := DecodeCursor(cursor)
	return err == nil
}
