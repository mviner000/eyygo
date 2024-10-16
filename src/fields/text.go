package fields

import (
	"fmt"
)

type TextField struct {
	FieldOptions
	MaxLength  int
	Validators []func(string) error
}

func (f TextField) GetOptions() FieldOptions {
	return f.FieldOptions
}

func (f TextField) SQLType() string {
	if f.MaxLength > 0 {
		return fmt.Sprintf("VARCHAR(%d)", f.MaxLength)
	}
	return "TEXT"
}

type SlugField struct {
	CharFieldType
	AllowUnicode bool
}

type EmailField struct {
	CharFieldType
}

type URLField struct {
	CharFieldType
}

type FilePathField struct {
	CharFieldType
	Path            string
	Match           string
	RecursiveSearch bool
	AllowFiles      bool
	AllowFolders    bool
}

type FileField struct {
	CharFieldType
	UploadTo string
	Storage  string
}

type ImageField struct {
	FileField
	Width  int
	Height int
}

type JSONField struct {
	TextField
	Encoder func(interface{}) ([]byte, error)
	Decoder func([]byte, interface{}) error
}

func (f JSONField) SQLType() string {
	return "JSON"
}

func CreateTextField(name string, fieldType string, options ...func(*FieldOptions)) Field {
	field := TextField{
		FieldOptions: FieldOptions{
			Name:     name,
			Editable: true,
			Unique:   false,
			Null:     false,
			Blank:    false,
		},
	}

	for _, option := range options {
		option(&field.FieldOptions)
	}

	switch fieldType {
	case "text":
		return field
	case "char":
		return CharFieldType{TextField: field}
	case "slug":
		return SlugField{CharFieldType: CharFieldType{TextField: field}}
	case "email":
		return EmailField{CharFieldType: CharFieldType{TextField: field}}
	case "url":
		return URLField{CharFieldType: CharFieldType{TextField: field}}
	case "filepath":
		return FilePathField{CharFieldType: CharFieldType{TextField: field}}
	case "file":
		return FileField{CharFieldType: CharFieldType{TextField: field}}
	case "image":
		return ImageField{FileField: FileField{CharFieldType: CharFieldType{TextField: field}}}
	case "json":
		return JSONField{TextField: field}
	default:
		return nil
	}
}

func WithTextValidators(validators ...func(string) error) func(*TextField) {
	return func(tf *TextField) {
		tf.Validators = validators
	}
}

func WithAllowUnicode(allow bool) func(*SlugField) {
	return func(sf *SlugField) {
		sf.AllowUnicode = allow
	}
}

func WithUploadTo(path string) func(*FileField) {
	return func(ff *FileField) {
		ff.UploadTo = path
	}
}

func WithStorage(storage string) func(*FileField) {
	return func(ff *FileField) {
		ff.Storage = storage
	}
}

func WithDimensions(width, height int) func(*ImageField) {
	return func(imageField *ImageField) {
		imageField.Width = width
		imageField.Height = height
	}
}
