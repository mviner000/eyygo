package germ

import "time"

// Model a basic GoLang struct which includes the following fields: ID, CreatedAt, UpdatedAt, DeletedAt
// It may be embedded into your model or you may build your own model without it
//
//	type User struct {
//	  germ.Model
//	}
type Model struct {
	ID        uint `germ:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt DeletedAt `germ:"index"`
}
