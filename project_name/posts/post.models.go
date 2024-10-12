package models

import (
	"time"
)

// Role represents user roles such as admin, user, etc.
type Role struct {
	ID   uint   `germ:"primaryKey"`
	Name string `germ:"uniqueIndex;not null"`

	Accounts []Account `germ:"foreignKey:RoleID"`
}

// Post represents a post made by an account
type Post struct {
	ID        uint   `germ:"primaryKey"`
	AccountID uint   `germ:"not null;index;foreignKey:accounts:ID"`
	Content   string `germ:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Account Account `germ:"foreignKey:AccountID"`
}
