package project_name

import (
	"fmt"

	models "github.com/mviner000/eyymi/eyygo/admin/models"
	"github.com/mviner000/eyymi/eyygo/registry"
)

// Register models in a single call
func RegisterModels() {
	registry.Model.Register(
		// &auth.AuthGroup{},
		// &auth.AuthPermission{},
		// &auth.AuthUser{},
		// &auth.EyygoContentType{},
		// &models.Role{},
		// &models.Account{},
		// &models.Post{},
		// &models.Comment{},
		// &models.Follower{},
		// &models.Like{},
		&models.AuthUser{},
		&models.AuthGroup{},
		&models.AdminLog{},
		&models.AuthPermission{},
		&models.EyygoContentType{},
		&models.Session{},
	)
	fmt.Println(registry.GetRegisteredModels())
}
