package cmd

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/mviner000/eyymi/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var CreateSuperuserCmd = &cobra.Command{
	Use:   "createsuperuser",
	Short: "Create a superuser with full access to the admin interface",
	Run: func(cmd *cobra.Command, args []string) {
		createSuperuser()
	},
}

type User struct {
	Username    string
	Email       string
	Password    string
	IsSuperuser bool
	IsStaff     bool
	IsActive    bool
	DateJoined  string
}

func createSuperuser() {
	dbURL := config.GetDatabaseURL()

	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}
	defer db.Close()

	user, err := promptUserDetails(db)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Debug: Log the created user details
	fmt.Printf("Debug: Created user - Username: %s, Email: %s, Password: %s\n", user.Username, user.Email, user.Password)

	if err := saveUser(user, db); err != nil {
		fmt.Printf("Error saving user: %v\n", err)
		return
	}

	fmt.Println("Superuser created successfully.")
}

func promptUserDetails(db *sql.DB) (*User, error) {
	user := &User{
		IsSuperuser: true,
		IsStaff:     true,
		IsActive:    true,
	}

	reader := bufio.NewReader(os.Stdin)

	// Prompt for username
	if err := promptUsername(user, reader, db); err != nil {
		return nil, err
	}

	// Prompt for email
	if err := promptEmail(user, reader, db); err != nil {
		return nil, err
	}

	// Prompt for password
	if err := promptPassword(user); err != nil {
		return nil, err
	}

	return user, nil
}

func promptUsername(user *User, reader *bufio.Reader, db *sql.DB) error {
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

func promptEmail(user *User, reader *bufio.Reader, db *sql.DB) error {
	// Regular expression for basic email validation
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

func promptPassword(user *User) error {
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

		hashedPassword, err := config.HashPassword(string(password))
		if err != nil {
			return fmt.Errorf("error hashing password: %w", err)
		}

		user.Password = hashedPassword
		return nil
	}
}

func saveUser(user *User, db *sql.DB) error {
	// The password is already hashed, so we can use it directly
	hashedPassword := user.Password

	// Set the current time as the date joined
	user.DateJoined = time.Now().Format("2006-01-02 15:04:05")

	// Insert the user into the database with the hashed password and date joined
	_, err := db.Exec(`
        INSERT INTO auth_user (username, email, password, is_superuser, is_staff, is_active, date_joined)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `, user.Username, user.Email, hashedPassword, user.IsSuperuser, user.IsStaff, user.IsActive, user.DateJoined)
	return err
}

func checkUsernameExists(username string, db *sql.DB) error {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM auth_user WHERE username=?)", username).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("username already exists")
	}
	return nil
}

func checkEmailExists(email string, db *sql.DB) error {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM auth_user WHERE email=?)", email).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
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
