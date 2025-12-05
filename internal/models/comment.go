package models

import (
	"time"
)

// Comment represents a comment on an event
type Comment struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	UserID    int       `json:"userId" gorm:"not null"`
	EventID   int       `json:"eventId" gorm:"not null"`
	Content   string    `json:"content" binding:"required,min=1" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt"`

	// Associations (optional, for preloading if needed)
	User  User  `json:"user,omitempty" gorm:"foreignKey:UserID" binding:"-"`
	Event Event `json:"-" gorm:"foreignKey:EventID" binding:"-"`
}
