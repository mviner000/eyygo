package handlers

import (
	"gorm.io/gorm"
)

type DBHandler struct {
	DB *gorm.DB
}

func NewDBHandler(db *gorm.DB) *DBHandler {
	return &DBHandler{
		DB: db,
	}
}

// Example of a database operation
func (h *DBHandler) CheckHealth() error {
	sqlDB, err := h.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
