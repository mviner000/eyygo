package models

import (
	"time"
)

// Borrower represents a person who can borrow books from the library
type Borrower struct {
	ID        uint   `germ:"primaryKey"`
	Name      string `germ:"not null"`
	Email     string `germ:"uniqueIndex;not null"`
	Password  string `germ:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Rentals []Rental `germ:"foreignKey:BorrowerID"`
}

// Book represents a book in the library
type Book struct {
	ID          uint   `germ:"primaryKey"`
	Title       string `germ:"not null"`
	AuthorID    uint   `germ:"not null;index"`
	PublishedAt time.Time
	CategoryID  uint `germ:"not null;index"`
	Stock       uint `germ:"not null"`

	Author   Author   `germ:"foreignKey:AuthorID"`
	Category Category `germ:"foreignKey:CategoryID"`
	Rentals  []Rental `germ:"foreignKey:BookID"`
}

// Author represents an author of a book
type Author struct {
	ID   uint   `germ:"primaryKey"`
	Name string `germ:"not null"`

	Books []Book `germ:"foreignKey:AuthorID"`
}

// Rental represents a book rental transaction
type Rental struct {
	ID         uint      `germ:"primaryKey"`
	BorrowerID uint      `germ:"not null;index"`
	BookID     uint      `germ:"not null;index"`
	RentedAt   time.Time `germ:"not null"`
	DueAt      time.Time `germ:"not null"`
	ReturnedAt *time.Time

	Borrower Borrower `germ:"foreignKey:BorrowerID"`
	Book     Book     `germ:"foreignKey:BookID"`
}

// Category represents categories for books
type Category struct {
	ID   uint   `germ:"primaryKey"`
	Name string `germ:"uniqueIndex;not null"`

	Books []Book `germ:"foreignKey:CategoryID"`
}
