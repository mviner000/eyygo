package settings

import (
	"fmt"
	"log"

	"github.com/mviner000/eyygo/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDBConnection(config *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch config.DBDriver {
	case "sqlite":
		dialector = sqlite.Open(fmt.Sprintf("%s.db", config.DBName))

	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.DBUser,
			config.DBPassword,
			config.DBHost,
			config.DBPort,
			config.DBName,
		)
		dialector = mysql.Open(dsn)

	case "postgresql":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			config.DBHost,
			config.DBUser,
			config.DBPassword,
			config.DBName,
			config.DBPort,
			config.DBSSLMode,
		)
		dialector = postgres.Open(dsn)

	default:
		return nil, fmt.Errorf("unsupported database driver: %s", config.DBDriver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Printf("Connected to %s database", config.DBDriver)
	return db, nil
}
