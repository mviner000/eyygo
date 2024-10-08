package models

import (
	"time"
)

type Post struct {
	ID        uint `germ:"primaryKey"`
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
