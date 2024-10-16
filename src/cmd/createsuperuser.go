package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"syscall"
	"time"

	models "github.com/mviner000/eyygo/src/admin/models"
	"github.com/mviner000/eyygo/src/config"
	"github.com/mviner000/eyygo/src/shared"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
	"gorm.io/gorm"
)

var CreateSuperuserCmd = &cobra.Command{
	Use:   "createsuperuser",
	Short: "Create a superuser with full access to the admin interface",
	Run: func(cmd *cobra.Command, args []string) {
		createSuperuser()
	},
}

func createSuperuser() {
	db := config.GetDB()
	if db == nil {
		fmt.Println("Failed to connect to database")
		return
	}

	user, err := promptUserDetails(db)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Debug: Log the created user details
	fmt.Printf("Debug: Created user - Username: %s, Email: %s\n", user.Username, user.Email)

	if err := saveUser(user, db); err != nil {
		fmt.Printf("Error saving user: %v\n", err)
		return
	}

	fmt.Println("Superuser created successfully.")
}

func promptUserDetails(db *gorm.DB) (*models.AuthUser, error) {
	user := &models.AuthUser{
		IsSuperuser: true,
		IsStaff:     true,
		IsActive:    true,
		DateJoined:  time.Now(),
	}

	reader := bufio.NewReader(os.Stdin)

	if err := promptUsername(user, reader, db); err != nil {
		return nil, err
	}

	if err := promptEmail(user, reader, db); err != nil {
		return nil, err
	}

	if err := promptPassword(user); err != nil {
		return nil, err
	}

	return user, nil
}

func promptUsername(user *models.AuthUser, reader *bufio.Reader, db *gorm.DB) error {
	for {
		fmt.Print("Username: ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)

		if username == "" {
			fmt.Println("Username cannot be blank.")
			continue
		}

		if err := checkUsernameExists(username, db); err != nil {
			fmt.Println(err)
			continue
		}

		user.Username = username
		return nil
	}
}

func promptEmail(user *models.AuthUser, reader *bufio.Reader, db *gorm.DB) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	for {
		fmt.Print("Email address: ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)

		if email == "" {
			fmt.Println("Email cannot be blank.")
			continue
		}

		if !emailRegex.MatchString(email) {
			fmt.Println("Invalid email format.")
			continue
		}

		if err := checkEmailExists(email, db); err != nil {
			fmt.Println(err)
			continue
		}

		user.Email = email
		return nil
	}
}

func promptPassword(user *models.AuthUser) error {
	for {
		fmt.Print("Password: ")
		password, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		fmt.Println()

		fmt.Print("Password (again): ")
		passwordConfirm, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		fmt.Println()

		if string(password) != string(passwordConfirm) {
			fmt.Println("Passwords do not match.")
			continue
		}

		if err := validatePassword(string(password)); err != nil {
			fmt.Println(err)
			continue
		}

		// Use the shared secret key for password hashing
		secretKey := []byte(shared.GetSecretKey())
		hashedPassword, err := bcrypt.GenerateFromPassword(append(password, secretKey...), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("error hashing password: %w", err)
		}

		user.Password = string(hashedPassword)
		return nil
	}
}

func saveUser(user *models.AuthUser, db *gorm.DB) error {
	result := db.Create(user)
	return result.Error
}

func checkUsernameExists(username string, db *gorm.DB) error {
	var count int64
	db.Model(&models.AuthUser{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		return fmt.Errorf("username already exists")
	}
	return nil
}

func checkEmailExists(email string, db *gorm.DB) error {
	var count int64
	db.Model(&models.AuthUser{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return fmt.Errorf("email already exists")
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	// Add more validation rules as needed, such as checking for numbers, special characters, etc.
	return nil
}
