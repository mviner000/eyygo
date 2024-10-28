// models/user.go
package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Username  string `gorm:"uniqueIndex;size:50;not null"`
	Email     string `gorm:"uniqueIndex;size:255;not null"`
	Password  string `gorm:"size:255;not null"`
	FirstName string `gorm:"size:50"`
	LastName  string `gorm:"size:50"`

	// Role flags
	IsSuperUser bool `gorm:"default:false"`
	IsStaff     bool `gorm:"default:false"`
	IsActive    bool `gorm:"default:true"`

	// Timestamps for user management
	LastLogin  *time.Time
	DateJoined time.Time `gorm:"autoCreateTime"`
}

// BeforeCreate hook to hash password before saving
func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifies the provided password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// CreateSuperUser creates a new superuser
func CreateSuperUser(db *gorm.DB, username, email, password string) error {
	user := &User{
		Username:    username,
		Email:       email,
		Password:    password,
		IsSuperUser: true,
		IsStaff:     true,
		IsActive:    true,
	}
	return db.Create(user).Error
}
