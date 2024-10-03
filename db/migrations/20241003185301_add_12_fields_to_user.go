package migrations

import (
	"gorm.io/gorm"
)

type add_12_fields_to_user struct{}

func (m *add_12_fields_to_user) Up(db *gorm.DB) error {
	// Add column created_at to users
	if err := db.Exec("your SQL here").Error; err != nil {
		return err
	}

	// Add column deleted_at to users
	if err := db.Exec("your SQL here").Error; err != nil {
		return err
	}

	// Add column password to users
	if err := db.Exec("your SQL here").Error; err != nil {
		return err
	}

	// Add column is_staff to users
	if err := db.Exec("your SQL here").Error; err != nil {
		return err
	}

	// Add column is_superuser to users
	if err := db.Exec("your SQL here").Error; err != nil {
		return err
	}

	// Add column  to users
	if err := db.Exec("your SQL here").Error; err != nil {
		return err
	}

	// Add column id to users
	if err := db.Exec("your SQL here").Error; err != nil {
		return err
	}

	// Add column updated_at to users
	if err := db.Exec("your SQL here").Error; err != nil {
		return err
	}

	// Add column username to users
	if err := db.Exec("your SQL here").Error; err != nil {
		return err
	}

	// Add column email to users
	if err := db.Exec("your SQL here").Error; err != nil {
		return err
	}

	// Add column date_joined to users
	if err := db.Exec("your SQL here").Error; err != nil {
		return err
	}

	// Add column is_active to users
	if err := db.Exec("your SQL here").Error; err != nil {
		return err
	}

	return nil
}

func (m *add_12_fields_to_user) Down(db *gorm.DB) error {
	// Reverse: Add column is_active to users
	if err := db.Exec("your reverse SQL here").Error; err != nil {
		return err
	}

	// Reverse: Add column date_joined to users
	if err := db.Exec("your reverse SQL here").Error; err != nil {
		return err
	}

	// Reverse: Add column email to users
	if err := db.Exec("your reverse SQL here").Error; err != nil {
		return err
	}

	// Reverse: Add column username to users
	if err := db.Exec("your reverse SQL here").Error; err != nil {
		return err
	}

	// Reverse: Add column updated_at to users
	if err := db.Exec("your reverse SQL here").Error; err != nil {
		return err
	}

	// Reverse: Add column id to users
	if err := db.Exec("your reverse SQL here").Error; err != nil {
		return err
	}

	// Reverse: Add column  to users
	if err := db.Exec("your reverse SQL here").Error; err != nil {
		return err
	}

	// Reverse: Add column is_superuser to users
	if err := db.Exec("your reverse SQL here").Error; err != nil {
		return err
	}

	// Reverse: Add column is_staff to users
	if err := db.Exec("your reverse SQL here").Error; err != nil {
		return err
	}

	// Reverse: Add column password to users
	if err := db.Exec("your reverse SQL here").Error; err != nil {
		return err
	}

	// Reverse: Add column deleted_at to users
	if err := db.Exec("your reverse SQL here").Error; err != nil {
		return err
	}

	// Reverse: Add column created_at to users
	if err := db.Exec("your reverse SQL here").Error; err != nil {
		return err
	}

	return nil
}