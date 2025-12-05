package models

// Event represents an event in the system
type Event struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	OwnerID     int    `json:"ownerId" gorm:"not null"`
	Owner       User   `json:"owner,omitempty" gorm:"foreignKey:OwnerID;references:ID;constraint:OnDelete:CASCADE"`
	Name        string `json:"name" binding:"required,min=3" gorm:"not null"`
	Description string `json:"description" binding:"required,min=10" gorm:"not null"`
	Date        string `json:"date" binding:"required,datetime=2006-01-02" gorm:"not null"`
	Location    string `json:"location" binding:"required,min=3" gorm:"not null"`
}
