package models

import (
	"time"

	"gorm.io/gorm"
)

type Note struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Title       string `gorm:"size:200;not null"`
	Content     string `gorm:"type:text"`
	Author      *User  `gorm:"foreignkey:AuthorID"`
	AuthorID    uint
	IsPublished bool   `gorm:"default:false"`
	Tags        string `gorm:"size:500"`
}
