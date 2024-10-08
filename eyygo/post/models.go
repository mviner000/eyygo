package models

import (
	"time"
)

// User represents a user in the social media platform
type User struct {
	ID        uint   `germ:"primaryKey"`
	Username  string `germ:"uniqueIndex;not null"`
	Email     string `germ:"uniqueIndex;not null"`
	Password  string `germ:"not null"`
	RoleID    uint   `germ:"not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Role      Role       `germ:"foreignKey:RoleID"`
	Posts     []Post     `germ:"foreignKey:UserID"`
	Following []Follower `germ:"foreignKey:FollowerUserID"`
	Followers []Follower `germ:"foreignKey:FollowedUserID"`
}

// Post represents a post made by a user
type Post struct {
	ID        uint   `germ:"primaryKey"`
	Title     string `germ:"not null"`
	Content   string `germ:"not null"`
	UserID    uint   `germ:"not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time

	User       User       `germ:"foreignKey:UserID"`
	Categories []Category `germ:"many2many:post_categories"`
}

// Follower represents a follower relationship between users
type Follower struct {
	ID             uint `germ:"primaryKey"`
	FollowerUserID uint `germ:"not null;index"`
	FollowedUserID uint `germ:"not null;index"`
	CreatedAt      time.Time

	FollowerUser User `germ:"foreignKey:FollowerUserID"`
	FollowedUser User `germ:"foreignKey:FollowedUserID"`
}

// Role represents user roles such as admin, user, etc.
type Role struct {
	ID   uint   `germ:"primaryKey"`
	Name string `germ:"uniqueIndex;not null"`

	Users []User `germ:"foreignKey:RoleID"`
}

// Category represents categories for posts
type Category struct {
	ID   uint   `germ:"primaryKey"`
	Name string `germ:"uniqueIndex;not null"`

	Posts []Post `germ:"many2many:post_categories"`
}
