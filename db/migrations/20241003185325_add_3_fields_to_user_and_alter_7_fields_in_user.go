package migrations

import (
	"gorm.io/gorm"
)

type add_3_fields_to_user_and_alter_7_fields_in_user struct{}

func (m *add_3_fields_to_user_and_alter_7_fields_in_user) Up(db *gorm.DB) error {
	// Add column last_name to users
	if err := db.Exec("your SQL here").Error; err != nil {
		return err
	}

	// Alter column is_staff in users
	if err := db.Exec("your SQL here").Error; err != nil {
		return err
	}

	// Add column first_name to users
	if err := db.Exec("your SQL here").Error; err != nil {
		return err
	}

	// Alter column is_active in users
	if err := db.Exec("your SQL here").Error; err != nil {
		return err
	}

	// Alter column username in users
	if err := db.Exec("your SQL here").Error; err != nil {
		return err
	}

	// Alter column email in users
	if err := db.Exec("your SQL here").Error; err != nil {
		return err
	}

	// Alter column is_superuser in users
	if err := db.Exec("your SQL here").Error; err != nil {
		return err
	}

	// Alter column id in users
	if err := db.Exec("your SQL here").Error; err != nil {
		return err
	}

	// Alter column password in users
	if err := db.Exec("your SQL here").Error; err != nil {
		return err
	}

	// Add column  to users
	if err := db.Exec("your SQL here").Error; err != nil {
		return err
	}

	return nil
}

func (m *add_3_fields_to_user_and_alter_7_fields_in_user) Down(db *gorm.DB) error {
	// Reverse: Add column  to users
	if err := db.Exec("your reverse SQL here").Error; err != nil {
		return err
	}

	// Reverse: Alter column password in users
	if err := db.Exec("your reverse SQL here").Error; err != nil {
		return err
	}

	// Reverse: Alter column id in users
	if err := db.Exec("your reverse SQL here").Error; err != nil {
		return err
	}

	// Reverse: Alter column is_superuser in users
	if err := db.Exec("your reverse SQL here").Error; err != nil {
		return err
	}

	// Reverse: Alter column email in users
	if err := db.Exec("your reverse SQL here").Error; err != nil {
		return err
	}

	// Reverse: Alter column username in users
	if err := db.Exec("your reverse SQL here").Error; err != nil {
		return err
	}

	// Reverse: Alter column is_active in users
	if err := db.Exec("your reverse SQL here").Error; err != nil {
		return err
	}

	// Reverse: Add column first_name to users
	if err := db.Exec("your reverse SQL here").Error; err != nil {
		return err
	}

	// Reverse: Alter column is_staff in users
	if err := db.Exec("your reverse SQL here").Error; err != nil {
		return err
	}

	// Reverse: Add column last_name to users
	if err := db.Exec("your reverse SQL here").Error; err != nil {
		return err
	}

	return nil
}