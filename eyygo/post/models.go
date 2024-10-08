package models

import (
	"time"

	"github.com/mviner000/eyymi/eyygo/germ"
	"github.com/mviner000/eyymi/eyygo/germ/driver/sqlite"
)

type Post struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func InitDB() (*germ.DB, error) {
	db, err := germ.Open(sqlite.Open("test.db"), &germ.Config{})
	if err != nil {
		return nil, err
	}

	// Automatically migrate the schema for the Post struct
	db.AutoMigrate(&Post{})
	return db, nil
}
