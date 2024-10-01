package types

import (
    "time"

    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username    string `gorm:"unique;not null"`
    Email       string `gorm:"unique;not null"`
    Password    string `gorm:"not null"`
    DateJoined  time.Time
    IsActive    bool `gorm:"default:true"`
    IsStaff     bool `gorm:"default:false"`
    IsSuperuser bool `gorm:"default:false"`
}