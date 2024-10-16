package project_name

import (
	"fmt"

	notes "github.com/mviner000/eyygo/project_name/notes"
	models "github.com/mviner000/eyygo/src/admin/models"
	"github.com/mviner000/eyygo/src/registry"
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
