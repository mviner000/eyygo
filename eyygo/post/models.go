package models

import (
	"time"
)

// User represents a user in the social media platform
type User struct {
	ID        uint   `germ:"primaryKey"`
	Username  string `germ:"unique;not null"`
	Email     string `germ:"unique;not null"`
	Password  string `germ:"not null"`
	RoleID    uint   `germ:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Post represents a post made by a user
type Post struct {
	ID        uint   `germ:"primaryKey"`
	Title     string `germ:"not null"`
	Content   string `germ:"not null"`
	UserID    uint   `germ:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Follower represents a follower relationship between users
type Follower struct {
	ID             uint `germ:"primaryKey"`
	FollowerUserID uint `germ:"not null"`
	FollowedUserID uint `germ:"not null"`
	CreatedAt      time.Time
}

// Role represents user roles such as admin, user, etc.
type Role struct {
	ID   uint   `germ:"primaryKey"`
	Name string `germ:"unique;not null"`
}

// Category represents categories for posts
type Category struct {
	ID   uint   `germ:"primaryKey"`
	Name string `germ:"unique;not null"`
}
