package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	models "github.com/mviner000/eyymi/project_name/posts" // Update with the correct import path
)

var db *gorm.DB

func InitDB(database *gorm.DB) {
	db = database
}

// GetUserByUsername retrieves a user from the database by username.
func GetUserByUsername(username string) (*models.AuthUser, error) {
	var user models.AuthUser
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		log.Printf("Error retrieving user %s from database: %v", username, err)
		return nil, err
	}

	log.Printf("User %s retrieved successfully from database", username)
	return &user, nil
}

// UpdateLastLogin updates the last login timestamp for a user in the database.
func UpdateLastLogin(userID uint) error {
	err := db.Model(&models.AuthUser{}).Where("id = ?", userID).Update("last_login", time.Now()).Error
	if err != nil {
		log.Printf("Error updating last_login for user ID %d: %v", userID, err)
	} else {
		log.Printf("Last login updated successfully for user ID %d", userID)
	}
	return err
}

// GetAllUsers retrieves all users from the database.
func GetAllUsers() ([]models.AuthUser, error) {
	var users []models.AuthUser
	err := db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetAllGroups retrieves all groups from the database.
func GetAllGroups() ([]string, error) {
	var groups []struct {
		Name string `gorm:"column:name"`
	}
	err := db.Table("auth_group").Select("name").Find(&groups).Error
	if err != nil {
		return nil, err
	}

	var groupNames []string
	for _, group := range groups {
		groupNames = append(groupNames, group.Name)
	}
	return groupNames, nil
}

// GetAllPermissions retrieves all permissions from the database.
func GetAllPermissions() ([]string, error) {
	var permissions []struct {
		Name string `gorm:"column:name"`
	}
	err := db.Table("auth_permission").Select("name").Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	var permissionNames []string
	for _, permission := range permissions {
		permissionNames = append(permissionNames, permission.Name)
	}
	return permissionNames, nil
}

// GetSessionFromDB retrieves session details from the database.
func GetSessionFromDB(c *fiber.Ctx) (uint, string, error) {
	sessionID := c.Cookies("hey_sesion")
	if sessionID == "" {
		return 0, "", fmt.Errorf("session ID not found in cookie")
	}

	var session struct {
		UserID     uint      `gorm:"column:user_id"`
		AuthToken  string    `gorm:"column:auth_token"`
		ExpireDate time.Time `gorm:"column:expire_date"`
	}
	err := db.Table("eyygo_session").Where("session_key = ?", sessionID).First(&session).Error
	if err != nil {
		return 0, "", fmt.Errorf("session not found")
	}

	if session.ExpireDate.Before(time.Now()) {
		return 0, "", fmt.Errorf("session expired")
	}

	return session.UserID, session.AuthToken, nil
}

// DeleteSessionFromDB deletes a session from the database.
func DeleteSessionFromDB(sessionID string) error {
	err := db.Where("session_key = ?", sessionID).Delete(&struct{}{}).Error
	if err != nil {
		log.Printf("Error deleting session %s from database: %v", sessionID, err)
		return err
	}
	log.Printf("Session %s deleted successfully from database", sessionID)
	return nil
}
