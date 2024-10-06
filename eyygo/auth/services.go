package auth

import (
	"database/sql"
	"log"
	"time"
)

var db *sql.DB

func InitDB(database *sql.DB) {
	db = database
}

func (u *User) ToAuthUser() *User {
	return &User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		LastLogin: u.LastLogin,
	}
}

// getUserByUsername retrieves a user from the database by username.
func GetUserByUsername(username string) (*User, error) {
	var user User
	var lastLogin sql.NullTime
	query := "SELECT id, username, email, password, last_login FROM auth_user WHERE username = ?"
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &lastLogin)
	if err != nil {
		log.Printf("Error retrieving user %s from database: %v", username, err)
		return nil, err
	}
	if lastLogin.Valid {
		user.LastLogin = lastLogin.Time
	} else {
		user.LastLogin = time.Time{} // Set to zero time if NULL
	}
	log.Printf("User %s retrieved successfully from database", username)
	return &user, nil
}

// updateLastLogin updates the last login timestamp for a user in the database.
func UpdateLastLogin(userID int) error {
	_, err := db.Exec("UPDATE auth_user SET last_login = ? WHERE id = ?", time.Now(), userID)
	if err != nil {
		log.Printf("Error updating last_login for user ID %d: %v", userID, err)
	} else {
		log.Printf("Last login updated successfully for user ID %d", userID)
	}
	return err
}
