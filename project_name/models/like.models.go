package models

import (
	"time"
)

// Like represents a like made on a post
type Like struct {
	ID        uint `germ:"primaryKey"`
	PostID    uint `germ:"not null;index;foreignKey:posts:ID"`
	AccountID uint `germ:"not null;index;foreignKey:accounts:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Post    Post    `germ:"foreignKey:PostID"`
	Account Account `germ:"foreignKey:AccountID"`
}
