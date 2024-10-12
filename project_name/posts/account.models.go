package models

import (
	"time"
)

// Account represents a user account in the social media platform
type Account struct {
	ID        uint   `germ:"primaryKey"`
	Username  string `germ:"uniqueIndex;not null"`
	Email     string `germ:"uniqueIndex;not null"`
	Password  string `germ:"not null"`
	RoleID    uint   `germ:"not null;index;foreignKey:roles:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Role Role `germ:"foreignKey:RoleID"`
}
