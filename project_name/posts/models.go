package models

import (
	"fmt"
	"time"

	"github.com/mviner000/eyymi/eyygo/registry"
)

// Register models in a single call
func RegisterModels() {
	registry.Model.Register(
		&Role{},
		&Account{},
		&Post{},
		&Comment{},
		&Follower{},
		&Like{},
	)
	fmt.Println(registry.GetRegisteredModels())
}

// Role represents user roles such as admin, user, etc.
type Role struct {
	ID   uint   `germ:"primaryKey"`
	Name string `germ:"uniqueIndex;not null"`

	Accounts []Account `germ:"foreignKey:RoleID"`
}

// Account represents a user account in the social media platform
type Account struct {
	ID        uint   `germ:"primaryKey"`
	Username  string `germ:"uniqueIndex;not null"`
	Email     string `germ:"uniqueIndex;not null"`
	Password  string `germ:"not null"`
	RoleID    uint   `germ:"not null;index;foreignKey:roles:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Role Role `germ:"foreignKey:RoleID"`
}

// Post represents a post made by an account
type Post struct {
	ID        uint   `germ:"primaryKey"`
	AccountID uint   `germ:"not null;index;foreignKey:accounts:ID"`
	Content   string `germ:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Account Account `germ:"foreignKey:AccountID"`
}

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

// Like represents a like made on a post
type Like struct {
	ID        uint `germ:"primaryKey"`
	PostID    uint `germ:"not null;index;foreignKey:posts:ID"`
	AccountID uint `germ:"not null;index;foreignKey:accounts:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Post    Post    `germ:"foreignKey:PostID"`
	Account Account `germ:"foreignKey:AccountID"`
}
