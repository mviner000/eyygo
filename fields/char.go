// fields/char.go
package fields

import (
	"fmt"
)

type CharFieldType struct {
	TextField
}

func (f CharFieldType) GetOptions() FieldOptions {
	return f.FieldOptions
}

func (f CharFieldType) SQLType() string {
	return fmt.Sprintf("VARCHAR(%d)", f.MaxLength)
}

func CharField(name string, maxLength int, options ...func(*FieldOptions)) Field {
	field := TextField{
		FieldOptions: FieldOptions{
			Name:     name,
			Editable: true,
			Unique:   false,
			Null:     false,
			Blank:    false,
		},
		MaxLength: maxLength,
	}

	for _, option := range options {
		option(&field.FieldOptions)
	}

	return CharFieldType{TextField: field}
}
