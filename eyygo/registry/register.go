package registry

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	mu               sync.Mutex
	registeredModels = make(map[string]interface{})
)

// RegisterModels registers multiple models at once.
func RegisterModels(models map[string]interface{}) {
	mu.Lock()
	defer mu.Unlock()

	for name, model := range models {
		registeredModels[name] = model
		fmt.Printf("Registered model: %s\n", name)
	}
}

// GetRegisteredModelNames returns the names of all registered models.
func GetRegisteredModelNames() []string {
	mu.Lock()
	defer mu.Unlock()

	names := make([]string, 0, len(registeredModels))
	for name := range registeredModels {
		names = append(names, name)
	}
	return names
}

// GetRegisteredModels returns a formatted string with the structure of all registered models.
func GetRegisteredModels() string {
	mu.Lock()
	defer mu.Unlock()

	modelInfo := ""
	for name, model := range registeredModels {
		modelInfo += fmt.Sprintf("Model: %s\n", name)
		modelType := reflect.TypeOf(model)

		for i := 0; i < modelType.NumField(); i++ {
			field := modelType.Field(i)
			modelInfo += fmt.Sprintf("  Field: %s, Type: %s, Tags: %s\n", field.Name, field.Type, field.Tag)
		}
		modelInfo += fmt.Sprintf("  Total Fields: %d\n\n", modelType.NumField()) // Add total fields count
	}
	return modelInfo
}

// GetRegisteredModel returns the registered model instance by name.
func GetRegisteredModel(name string) (interface{}, bool) {
	mu.Lock()
	defer mu.Unlock()

	model, exists := registeredModels[name]
	return model, exists
}
