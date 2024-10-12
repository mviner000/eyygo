package models

import (
	"time"
)

type AuthUser struct {
	ID               uint   `germ:"primaryKey;autoIncrement"` // Primary key with auto-increment
	Password         string `germ:"not null"`
	LastLogin        *time.Time
	IsSuperuser      bool      `germ:"not null"`
	Username         string    `germ:"unique;not null"`
	FirstName        string    `germ:"default:''"`
	LastName         string    `germ:"default:''"`
	Email            string    `germ:"default:''"`
	IsStaff          bool      `germ:"not null"`
	IsActive         bool      `germ:"not null"`
	DateJoined       time.Time `germ:"not null"`
	GroupID          *uint
	Group            *AuthGroup
	UserPermissionID *uint
	UserPermission   *AuthPermission
}

type AuthGroup struct {
	ID   uint   `germ:"primaryKey;autoIncrement"`
	Name string `germ:"unique;not null"`
}

type AuthPermission struct {
	ID            uint             `germ:"primaryKey;autoIncrement"`
	Name          string           `germ:"not null"`
	ContentTypeID uint             `germ:"not null"`
	Codename      string           `germ:"not null"`
	ContentType   EyygoContentType `germ:"foreignKey:ContentTypeID"`

	// Use a composite unique index for content_type_id and codename
	UniqueIndex string `germ:"uniqueIndex:content_type_codename,priority:1"`
}

type EyygoContentType struct {
	ID   uint   `germ:"primaryKey;autoIncrement"`
	Name string `germ:"not null"`
}

type Session struct {
	SessionKey string    `germ:"primaryKey;type:TEXT"` // Use TEXT for session key
	ExpireDate time.Time `germ:"not null;index"`
	UserID     uint      `germ:"not null"`                        // Foreign key to AuthUser.ID
	User       AuthUser  `germ:"foreignKey:UserID;references:ID"` // Define the foreign key correctly
	AuthToken  string    `germ:"not null"`
}

type AdminLog struct {
	ID            uint              `germ:"primaryKey;autoIncrement"`
	ActionTime    time.Time         `germ:"not null;default:CURRENT_TIMESTAMP"` // Auto-set timestamp
	ObjectID      string            `germ:"size:255"`                           // Optional field for object identifier
	ObjectRepr    string            `germ:"not null;size:200"`                  // Representation of the object
	ActionFlag    int               `germ:"not null"`                           // Action flag as an integer
	ChangeMessage string            `germ:"not null"`                           // Change message
	ContentTypeID *uint             // Foreign key to EyygoContentType.ID (nullable)
	ContentType   *EyygoContentType `germ:"foreignKey:ContentTypeID;references:ID;onDelete:SET NULL"`
	UserID        uint              // Foreign key to AuthUser.ID
	User          AuthUser          `germ:"foreignKey:UserID;references:ID;onDelete:CASCADE"`
}

// Table name overrides
func (EyygoContentType) TableName() string {
	return "eyygo_content_type"
}

func (Session) TableName() string {
	return "eyygo_session"
}

func (AdminLog) TableName() string {
	return "eyygo_admin_log"
}

func (AuthUser) TableName() string {
	return "auth_user"
}

func (AuthPermission) TableName() string {
	return "auth_permission"
}

func (AuthGroup) TableName() string {
	return "auth_group"
}
