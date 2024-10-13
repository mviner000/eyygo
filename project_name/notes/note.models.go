package notes

import (
	"time"
)

type Note struct {
	ID        uint   `germ:"primaryKey"`
	Content   string `germ:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
