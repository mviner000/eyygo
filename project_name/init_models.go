package project_name

import (
	"fmt"

	models "github.com/mviner000/eyymi/eyygo/admin/models"
	"github.com/mviner000/eyymi/eyygo/registry"
	notes "github.com/mviner000/eyymi/project_name/notes"
)

// Register models in a single call
func RegisterModels() {
	registry.Model.Register(

		&models.AuthUser{},
		&models.AuthGroup{},
		&models.AdminLog{},
		&models.AuthPermission{},
		&models.EyygoContentType{},
		&models.Session{},

		&notes.Note{},
	)
	fmt.Println(registry.GetRegisteredModels())
}
