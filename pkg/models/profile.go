package models

import "time"

type Profile struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	UserID    int    `json:"userId" gorm:"uniqueIndex;not null"`
	User      User   `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Bio       string `json:"bio" gorm:"type:text"`
	AvatarURL string `json:"avatarUrl" gorm:"size:500"`
	Phone     string `json:"phone" gorm:"size:25"`

	DateOfBirth *time.Time `json:"dateOfBirth"`
	Country     string     `json:"country" gorm:"size:100"`
	City        string     `json:"city" gorm:"size:100"`
	Timezone    string     `json:"timezone" gorm:"size:100;default:'UTC'"`

	Website  string `json:"website" gorm:"size:255"`
	Twitter  string `json:"twitter" gorm:"size:255"`
	LinkedIn string `json:"linkedin" gorm:"size:255"`
	GitHub   string `json:"github" gorm:"size:255"`

	Theme              string    `json:"theme" gorm:"size:20;default:light"`
	Language           string    `json:"language" gorm:"size:10;default:en"`
	IsPublic           bool      `json:"isPublic" gorm:"default:true"`
	EmailNotifications bool      `json:"emailNotifications" gorm:"default:true"`
	PushNotifications  bool      `json:"pushNotifications" gorm:"default:true"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}
