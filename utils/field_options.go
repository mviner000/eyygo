package utils

import (
	"fmt"
	"strings"
)

type FieldOptions struct {
	Name          string
	Choices       map[string]string
	DBIndex       bool
	Editable      bool
	HelpText      string
	VerboseName   string
	ErrorMessages map[string]string
	UniqueForDate string
	Unique        bool
	Null          bool
	Blank         bool
	DefaultValue  interface{}
}

type Field interface {
	GetOptions() FieldOptions
	SQLType() string
}

func WithChoices(choices map[string]string) func(*FieldOptions) {
	return func(f *FieldOptions) {
		f.Choices = choices
	}
}

func WithDBIndex(index bool) func(*FieldOptions) {
	return func(f *FieldOptions) {
		f.DBIndex = index
	}
}

func WithEditable(editable bool) func(*FieldOptions) {
	return func(f *FieldOptions) {
		f.Editable = editable
	}
}

func WithHelpText(text string) func(*FieldOptions) {
	return func(f *FieldOptions) {
		f.HelpText = text
	}
}

func WithVerboseName(name string) func(*FieldOptions) {
	return func(f *FieldOptions) {
		f.VerboseName = name
	}
}

func WithErrorMessages(messages map[string]string) func(*FieldOptions) {
	return func(f *FieldOptions) {
		f.ErrorMessages = messages
	}
}

func WithUniqueForDate(dateField string) func(*FieldOptions) {
	return func(f *FieldOptions) {
		f.UniqueForDate = dateField
	}
}

func WithUnique(unique bool) func(*FieldOptions) {
	return func(f *FieldOptions) {
		f.Unique = unique
	}
}

func WithNull(null bool) func(*FieldOptions) {
	return func(f *FieldOptions) {
		f.Null = null
	}
}

func WithBlank(blank bool) func(*FieldOptions) {
	return func(f *FieldOptions) {
		f.Blank = blank
	}
}

func WithDefaultValue(value interface{}) func(*FieldOptions) {
	return func(f *FieldOptions) {
		f.DefaultValue = value
	}
}

func FieldDefinition(field Field) string {
	options := field.GetOptions()
	def := []string{options.Name, field.SQLType()}

	if options.Unique {
		def = append(def, "UNIQUE")
	}
	if options.Null {
		def = append(def, "NULL")
	} else {
		def = append(def, "NOT NULL")
	}
	if options.DBIndex {
		def = append(def, "INDEX")
	}
	if options.DefaultValue != nil {
		def = append(def, fmt.Sprintf("DEFAULT %v", options.DefaultValue))
	}

	return strings.Join(def, " ")
}
