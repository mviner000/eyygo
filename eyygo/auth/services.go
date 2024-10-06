package auth

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
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

func GetAllUsers() ([]User, error) {
	query := "SELECT id, username, email, last_login FROM auth_user"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var lastLogin sql.NullTime
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &lastLogin)
		if err != nil {
			return nil, err
		}
		if lastLogin.Valid {
			user.LastLogin = lastLogin.Time
		}
		users = append(users, user)
	}
	return users, nil
}

func GetAllGroups() ([]string, error) {
	query := "SELECT name FROM auth_group"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []string
	for rows.Next() {
		var groupName string
		err := rows.Scan(&groupName)
		if err != nil {
			return nil, err
		}
		groups = append(groups, groupName)
	}
	return groups, nil
}

func GetAllPermissions() ([]string, error) {
	query := "SELECT name FROM auth_permission"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var permissionName string
		err := rows.Scan(&permissionName)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permissionName)
	}
	return permissions, nil
}

func GetSessionFromDB(c *fiber.Ctx) (int, string, error) {
	sessionID := c.Cookies("session_id")
	if sessionID == "" {
		return 0, "", fmt.Errorf("session ID not found in cookie")
	}

	var userIDStr string
	var authToken string
	var expireDate time.Time

	query := `SELECT user_id, auth_token, expire_date FROM eyygo_session WHERE session_key = ?`
	err := db.QueryRow(query, sessionID).Scan(&userIDStr, &authToken, &expireDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", fmt.Errorf("session not found")
		}
		return 0, "", err
	}

	if expireDate.Before(time.Now()) {
		return 0, "", fmt.Errorf("session expired")
	}

	// Convert userID from string to int
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, "", fmt.Errorf("invalid user ID in session")
	}

	return userID, authToken, nil
}
