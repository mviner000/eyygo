package models

import (
	"fmt"

	"github.com/mviner000/eyymi/eyygo/registry"
)

// Register models in a single call
func RegisterModels() {
	registry.Model.Register(
		&Role{},
		&Account{},
		&Post{},
		&Comment{},
		&Follower{},
		&Like{},
	)
	fmt.Println(registry.GetRegisteredModels())
}
