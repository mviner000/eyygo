package models

import (
	"time"
)

// Comment represents a comment made on a post
type Comment struct {
	ID        uint   `germ:"primaryKey"`
	PostID    uint   `germ:"not null;index;foreignKey:posts:ID"`
	AccountID uint   `germ:"not null;index;foreignKey:accounts:ID"`
	Content   string `germ:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Post    Post    `germ:"foreignKey:PostID"`
	Account Account `germ:"foreignKey:AccountID"`
}
