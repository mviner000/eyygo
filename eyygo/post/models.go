package models

import (
	"time"
)

// Role represents user roles such as admin, user, etc.
type Role struct {
	ID   uint   `germ:"primaryKey"`
	Name string `germ:"uniqueIndex;not null"`

	Users []User `germ:"foreignKey:RoleID"`
}

// User represents a user in the social media platform
type User struct {
	ID        uint   `germ:"primaryKey"`
	Username  string `germ:"uniqueIndex;not null"`
	Email     string `germ:"uniqueIndex;not null"`
	Password  string `germ:"not null"`
	RoleID    uint   `germ:"not null;index;foreignKey:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Role Role `germ:"foreignKey:RoleID"`
}

// Post represents a post made by a user
type Post struct {
	ID        uint   `germ:"primaryKey"`
	UserID    uint   `germ:"not null;index;foreignKey:ID"`
	Content   string `germ:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	User User `germ:"foreignKey:UserID"`
}

// Comment represents a comment made on a post
type Comment struct {
	ID        uint   `germ:"primaryKey"`
	PostID    uint   `germ:"not null;index;foreignKey:ID"`
	UserID    uint   `germ:"not null;index;foreignKey:ID"`
	Content   string `germ:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Post Post `germ:"foreignKey:PostID"`
	User User `germ:"foreignKey:UserID"`
}

// Follower represents a user who is following another user
type Follower struct {
	ID         uint `germ:"primaryKey"`
	UserID     uint `germ:"not null;index;foreignKey:ID"`
	FollowerID uint `germ:"not null;index;foreignKey:ID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	User     User `germ:"foreignKey:UserID"`
	Follower User `germ:"foreignKey:FollowerID"`
}

// Like represents a like made on a post
type Like struct {
	ID        uint `germ:"primaryKey"`
	PostID    uint `germ:"not null;index;foreignKey:ID"`
	UserID    uint `germ:"not null;index;foreignKey:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Post Post `germ:"foreignKey:PostID"`
	User User `germ:"foreignKey:UserID"`
}
