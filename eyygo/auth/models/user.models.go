package models

import (
	"time"
)

// AuthUser represents the auth_user table
type AuthUser struct {
	ID                uint           `germ:"primaryKey;autoIncrement"`
	Password          string         `germ:"not null"`
	LastLogin         *time.Time     `germ:"type:datetime"`
	IsSuperuser       bool           `germ:"not null;check:is_superuser IN (0, 1)"`
	Username          string         `germ:"unique;not null"`
	FirstName         string         `germ:"default:''"`
	LastName          string         `germ:"default:''"`
	Email             string         `germ:"default:''"`
	IsStaff           bool           `germ:"not null;check:is_staff IN (0, 1)"`
	IsActive          bool           `germ:"not null;check:is_active IN (0, 1)"`
	DateJoined        time.Time      `germ:"not null;type:datetime"`
	GroupsID          *uint          `germ:"default:null"`
	UserPermissionsID *uint          `germ:"default:null"`
	Groups            AuthGroup      `germ:"foreignKey:GroupsID"`
	UserPermissions   AuthPermission `germ:"foreignKey:UserPermissionsID"`
}
