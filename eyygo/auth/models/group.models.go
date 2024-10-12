package models

// AuthGroup represents the auth_group table
type AuthGroup struct {
	ID   uint   `germ:"primaryKey;autoIncrement"`
	Name string `germ:"unique;not null"`
}
