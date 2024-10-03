package fields

import (
	"fmt"
)

type TextField struct {
	Options
	MaxLength  int
	Validators []func(string) error
}

func (f TextField) GetOptions() Options {
	return f.Options
}

func (f TextField) SQLType() string {
	if f.MaxLength > 0 {
		return fmt.Sprintf("VARCHAR(%d)", f.MaxLength)
	}
	return "TEXT"
}

type CharField struct {
	TextField
}

func (f CharField) SQLType() string {
	return fmt.Sprintf("VARCHAR(%d)", f.MaxLength)
}

type SlugField struct {
	CharField
	AllowUnicode bool
}

type EmailField struct {
	CharField
}

type URLField struct {
	CharField
}

type FilePathField struct {
	CharField
	Path            string
	Match           string
	RecursiveSearch bool
	AllowFiles      bool
	AllowFolders    bool
}

type FileField struct {
	CharField
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

func CreateTextField(name string, fieldType string, options ...func(*Options)) Field {
	field := TextField{
		Options: Options{
			Name:     name,
			Editable: true,
			Unique:   false,
			Null:     false,
			Blank:    false,
		},
	}

	for _, option := range options {
		option(&field.Options)
	}

	switch fieldType {
	case "text":
		return field
	case "char":
		return CharField{TextField: field}
	case "slug":
		return SlugField{CharField: CharField{TextField: field}}
	case "email":
		return EmailField{CharField: CharField{TextField: field}}
	case "url":
		return URLField{CharField: CharField{TextField: field}}
	case "filepath":
		return FilePathField{CharField: CharField{TextField: field}}
	case "file":
		return FileField{CharField: CharField{TextField: field}}
	case "image":
		return ImageField{FileField: FileField{CharField: CharField{TextField: field}}}
	case "json":
		return JSONField{TextField: field}
	default:
		return nil
	}
}

func WithMaxLength(length int) func(*TextField) {
	return func(tf *TextField) {
		tf.MaxLength = length
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
