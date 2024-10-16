package admin

import (
	"log"

	"github.com/mviner000/eyymi/src/auth"
	"github.com/mviner000/eyymi/src/config"
)

var tokenGenerator *auth.PasswordResetTokenGenerator

func init() {
	// Initialize the database connection in the config package
	db := config.GetDB()
	if db == nil {
		log.Fatalf("Failed to connect to database")
	}
	log.Println("Successfully connected to the database")

	// Pass *sql.DB to auth.InitDB
	auth.InitDB(db)

	// Initialize the token generator
	tokenGenerator = auth.NewPasswordResetTokenGenerator()
}
