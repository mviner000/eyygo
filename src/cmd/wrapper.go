package cmd

import (
	"log"
	"sync"

	"github.com/mviner000/eyymi/src/germ"
)

type GERMWrapper struct {
	DB   *germ.DB
	once sync.Once
}

func NewGERMWrapper(db *germ.DB) *GERMWrapper {
	return &GERMWrapper{DB: db}
}

func (w *GERMWrapper) Close() {
	w.once.Do(func() {
		sqlDB, err := w.DB.DB()
		if err != nil {
			log.Printf("Error getting underlying SQL database: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	})
}
