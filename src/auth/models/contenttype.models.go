package models

// EyygoContentType represents the eyygo_content_type table
type EyygoContentType struct {
	ID   uint   `germ:"primaryKey;autoIncrement"`
	Name string `germ:"not null"`
}
