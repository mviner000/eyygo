package fields

import (
	"fmt"
)

type BooleanField struct {
	FieldOptions
}

func (f BooleanField) SQLType() string {
	return "BOOLEAN"
}

func (f BooleanField) GetOptions() FieldOptions {
	return f.FieldOptions
}

// Common field options
var (
	RequiredChar   = []func(*FieldOptions){WithRequired(true)}
	RequiredUnique = []func(*FieldOptions){WithRequired(true), WithUnique(true)}
	NullableChar   = []func(*FieldOptions){WithNull(true), WithBlank(true)}
)

type DateField struct {
	FieldOptions
	AutoNow    bool
	AutoNowAdd bool
}

func (f DateField) SQLType() string {
	return "DATE"
}

func (f DateField) GetOptions() FieldOptions {
	return f.FieldOptions
}

type TimeField struct {
	FieldOptions
	AutoNow    bool
	AutoNowAdd bool
}

func (f TimeField) SQLType() string {
	return "TIME"
}

func (f TimeField) GetOptions() FieldOptions {
	return f.FieldOptions
}

type DateTimeField struct {
	FieldOptions
	AutoNow    bool
	AutoNowAdd bool
}

func (f DateTimeField) SQLType() string {
	return "DATETIME"
}

func (f DateTimeField) GetOptions() FieldOptions {
	return f.FieldOptions
}

type UUIDField struct {
	FieldOptions
}

func (f UUIDField) SQLType() string {
	return "UUID"
}

func (f UUIDField) GetOptions() FieldOptions {
	return f.FieldOptions
}

type BinaryField struct {
	FieldOptions
	MaxLength int
}

func (f BinaryField) SQLType() string {
	if f.MaxLength > 0 {
		return fmt.Sprintf("VARBINARY(%d)", f.MaxLength)
	}
	return "BLOB"
}

func (f BinaryField) GetOptions() FieldOptions {
	return f.FieldOptions
}

type IPAddressField struct {
	FieldOptions
}

func (f IPAddressField) SQLType() string {
	return "VARCHAR(39)" // IPv6 addresses can be up to 39 characters long
}

func (f IPAddressField) GetOptions() FieldOptions {
	return f.FieldOptions
}

type GenericIPAddressField struct {
	FieldOptions
	Protocol string // 'both', 'IPv4', or 'IPv6'
}

func (f GenericIPAddressField) SQLType() string {
	return "VARCHAR(39)" // IPv6 addresses can be up to 39 characters long
}

func (f GenericIPAddressField) GetOptions() FieldOptions {
	return f.FieldOptions
}

type AutoField struct {
	FieldOptions
}

func (f AutoField) SQLType() string {
	return "SERIAL PRIMARY KEY"
}

func (f AutoField) GetOptions() FieldOptions {
	return f.FieldOptions
}

type ForeignKey struct {
	FieldOptions
	RelatedModel string
	OnDelete     string
}

func (f ForeignKey) SQLType() string {
	return fmt.Sprintf("INTEGER REFERENCES %s(id)", f.RelatedModel)
}

func (f ForeignKey) GetOptions() FieldOptions {
	return f.FieldOptions
}

type ManyToManyField struct {
	FieldOptions
	RelatedModel string
	ThroughModel string
}

func (f ManyToManyField) SQLType() string {
	// This is a special case and doesn't directly translate to a SQL type
	return "MANYTOMANY"
}

func (f ManyToManyField) GetOptions() FieldOptions {
	return f.FieldOptions
}

type OneToOneField struct {
	FieldOptions
	RelatedModel string
	OnDelete     string
}

func (f OneToOneField) SQLType() string {
	return fmt.Sprintf("INTEGER UNIQUE REFERENCES %s(id)", f.RelatedModel)
}

func (f OneToOneField) GetOptions() FieldOptions {
	return f.FieldOptions
}

func CreateSpecialField(name string, fieldType string, options ...func(*FieldOptions)) Field {
	field := FieldOptions{
		Name:     name,
		Editable: true,
		Unique:   false,
		Null:     false,
		Blank:    false,
	}

	for _, option := range options {
		option(&field)
	}

	switch fieldType {
	case "boolean":
		return BooleanField{FieldOptions: field}
	case "date":
		return DateField{FieldOptions: field}
	case "time":
		return TimeField{FieldOptions: field}
	case "datetime":
		return DateTimeField{FieldOptions: field}
	case "uuid":
		return UUIDField{FieldOptions: field}
	case "binary":
		return BinaryField{FieldOptions: field}
	case "ip":
		return IPAddressField{FieldOptions: field}
	case "genericip":
		return GenericIPAddressField{FieldOptions: field}
	case "auto":
		return AutoField{FieldOptions: field}
	case "foreignkey":
		return ForeignKey{FieldOptions: field}
	case "manytomany":
		return ManyToManyField{FieldOptions: field}
	case "onetoone":
		return OneToOneField{FieldOptions: field}
	default:
		return nil
	}
}

func WithAutoNow(autoNow bool) func(interface{}) {
	return func(f interface{}) {
		switch field := f.(type) {
		case *DateField:
			field.AutoNow = autoNow
		case *TimeField:
			field.AutoNow = autoNow
		case *DateTimeField:
			field.AutoNow = autoNow
		}
	}
}

func WithAutoNowAdd(autoNowAdd bool) func(interface{}) {
	return func(f interface{}) {
		switch field := f.(type) {
		case *DateField:
			field.AutoNowAdd = autoNowAdd
		case *TimeField:
			field.AutoNowAdd = autoNowAdd
		case *DateTimeField:
			field.AutoNowAdd = autoNowAdd
		}
	}
}

func WithOnDelete(onDelete string) func(interface{}) {
	return func(f interface{}) {
		switch field := f.(type) {
		case *ForeignKey:
			field.OnDelete = onDelete
		case *OneToOneField:
			field.OnDelete = onDelete
		}
	}
}

func WithThroughModel(throughModel string) func(interface{}) {
	return func(f interface{}) {
		if field, ok := f.(*ManyToManyField); ok {
			field.ThroughModel = throughModel
		}
	}
}

func WithProtocol(protocol string) func(interface{}) {
	return func(f interface{}) {
		if field, ok := f.(*GenericIPAddressField); ok {
			field.Protocol = protocol
		}
	}
}
