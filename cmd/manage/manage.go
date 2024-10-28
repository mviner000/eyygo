package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/mviner000/eyygo/config"
	"github.com/mviner000/eyygo/models"
	"github.com/mviner000/eyygo/settings"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	green  = color.New(color.FgGreen).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
)

// Password validation rules
const (
	minLength = 12
	maxLength = 128
)

var rootCmd = &cobra.Command{
	Use:   "manage",
	Short: "Management commands for the application",
}

var createSuperUserCmd = &cobra.Command{
	Use:   "createsuperuser",
	Short: "Create a new superuser",
	Run:   createSuperUser,
}

func init() {
	rootCmd.AddCommand(createSuperUserCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(red("Error:"), err)
		os.Exit(1)
	}
}

func createSuperUser(cmd *cobra.Command, args []string) {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println(red("Error loading config:"), err)
		os.Exit(1)
	}

	// Initialize database connection
	db, err := settings.NewDBConnection(cfg)
	if err != nil {
		fmt.Println(red("Error connecting to database:"), err)
		os.Exit(1)
	}

	// Display password requirements
	fmt.Println(cyan("\nPassword Requirements:"))
	fmt.Println(yellow("• Minimum length:"), minLength, "characters")
	fmt.Println(yellow("• Maximum length:"), maxLength, "characters")
	fmt.Println(yellow("• Must contain:"))
	fmt.Println("  - At least one uppercase letter")
	fmt.Println("  - At least one lowercase letter")
	fmt.Println("  - At least one number")
	fmt.Println("  - At least one special character (!@#$%^&*(),.?\":{}|<>)")
	fmt.Println("  - No common words or patterns")
	fmt.Println(yellow("• Best practices:"))
	fmt.Println("  - Use a passphrase with random words")
	fmt.Println("  - Mix numbers and special characters within the phrase")
	fmt.Println("  - Avoid personal information")
	fmt.Println("")

	// Get username
	var username string
	for {
		username = promptString("Username")
		if err := validateUsername(username); err != nil {
			fmt.Println(red("Error:"), err)
			continue
		}
		// Check if username exists
		var count int64
		db.Model(&models.User{}).Where("username = ?", username).Count(&count)
		if count > 0 {
			fmt.Println(red("Error: Username already exists"))
			continue
		}
		break
	}

	// Get email
	var email string
	for {
		email = promptString("Email")
		if err := validateEmail(email); err != nil {
			fmt.Println(red("Error:"), err)
			continue
		}
		// Check if email exists
		var count int64
		db.Model(&models.User{}).Where("email = ?", email).Count(&count)
		if count > 0 {
			fmt.Println(red("Error: Email already exists"))
			continue
		}
		break
	}

	// Get password
	var password string
	for {
		password = promptPassword("Password")
		if err := validatePassword(password); err != nil {
			fmt.Println(red("Error:"), err)
			continue
		}

		confirmation := promptPassword("Password (confirm)")
		if password != confirmation {
			fmt.Println(red("Error: Passwords don't match"))
			continue
		}
		break
	}

	// Create superuser
	if err := models.CreateSuperUser(db, username, email, password); err != nil {
		fmt.Println(red("\nError creating superuser:"), err)
		os.Exit(1)
	}

	fmt.Println(green("\nSuperuser created successfully!"))
	fmt.Printf("Username: %s\n", cyan(username))
	fmt.Printf("Email: %s\n", cyan(email))
}

func promptString(prompt string) string {
	fmt.Printf("%s: ", prompt)
	var value string
	fmt.Scanln(&value)
	return strings.TrimSpace(value)
}

func promptPassword(prompt string) string {
	fmt.Printf("%s: ", prompt)
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println() // Add newline after password input
	if err != nil {
		fmt.Println(red("Error reading password:"), err)
		os.Exit(1)
	}
	return string(password)
}

func validateUsername(username string) error {
	if len(username) < 3 || len(username) > 30 {
		return fmt.Errorf("username must be between 3 and 30 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username) {
		return fmt.Errorf("username can only contain letters, numbers, and underscores")
	}
	return nil
}

func validateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < minLength {
		return fmt.Errorf("password must be at least %d characters long", minLength)
	}
	if len(password) > maxLength {
		return fmt.Errorf("password must not exceed %d characters", maxLength)
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one number")
	}
	if !regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one special character")
	}

	// Check for common patterns
	commonPatterns := []string{
		"password", "123", "abc", "qwerty", "admin", "letmein",
		"welcome", "monkey", "dragon", "master", "superman", "batman",
	}
	lowercasePassword := strings.ToLower(password)
	for _, pattern := range commonPatterns {
		if strings.Contains(lowercasePassword, pattern) {
			return fmt.Errorf("password contains a common pattern: %s", pattern)
		}
	}

	return nil
}
