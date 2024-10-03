package admin

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string `gorm:"unique;not null"`
	Email       string `gorm:"unique;not null"`
	Password    string `gorm:"not null"`
	DateJoined  time.Time
	IsActive    bool    `gorm:"default:true"`
	IsStaff     bool    `gorm:"default:false"`
	IsSuperuser bool    `gorm:"default:false"`
	Groups      []Group `gorm:"many2many:user_groups;"`
}

type Group struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description string
	Users       []User `gorm:"many2many:user_groups;"`
}

// UserGroup represents the many-to-many relationship between User and Group
type UserGroup struct {
	UserID  uint `gorm:"primaryKey"`
	GroupID uint `gorm:"primaryKey"`
}
