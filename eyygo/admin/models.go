package admin

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mviner000/eyymi/eyygo/config"
	"github.com/mviner000/eyymi/eyygo/fields"
	"github.com/mviner000/eyymi/eyygo/operations"
)

// GetModels defines and returns all the models for the application.
func GetModels() []*operations.Model {
	// Define a new model for the User
	userModel := operations.NewModel("users")
	userModel.AddField(fields.CharField("username", 255, fields.WithRequired(true)))
	userModel.AddField(fields.CharField("lastname", 255, fields.WithRequired(true)))
	userModel.AddField(fields.CharField("firstname", 255, fields.WithRequired(true)))
	userModel.AddField(fields.CharField("middlename", 255, fields.WithRequired(true)))
	userModel.AddField(fields.CharField("email", 255, fields.WithRequired(true), fields.WithUnique(true)))

	// Generate SQL for creating the table (for debugging purposes)
	sql := userModel.CreateTableSQL()
	fmt.Println(sql)

	// Return the models
	return []*operations.Model{userModel}
}

// GetAllUsers returns a list of users from the database
func GetAllUsers() []User {
	dbURL := config.GetDatabaseURL()
	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT username, email FROM users")
	if err != nil {
		log.Fatalf("Failed to query users: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Username, &user.Email); err != nil {
			log.Fatalf("Failed to scan user row: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Row iteration error: %v", err)
	}

	log.Printf("Retrieved %d users from the database", len(users))
	return users
}
