package models

import (
	"time"
)

// Follower represents an account who is following another account
type Follower struct {
	ID         uint `germ:"primaryKey"`
	AccountID  uint `germ:"not null;index;foreignKey:accounts:ID"`
	FollowerID uint `germ:"not null;index;foreignKey:accounts:ID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Account  Account `germ:"foreignKey:AccountID"`
	Follower Account `germ:"foreignKey:FollowerID"`
}
