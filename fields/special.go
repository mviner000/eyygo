package fields

import (
	"fmt"
)

type BooleanField struct {
	Options
}

func (f BooleanField) SQLType() string {
	return "BOOLEAN"
}

func (f BooleanField) GetOptions() Options {
	return f.Options
}

type DateField struct {
	Options
	AutoNow    bool
	AutoNowAdd bool
}

func (f DateField) SQLType() string {
	return "DATE"
}

func (f DateField) GetOptions() Options {
	return f.Options
}

type TimeField struct {
	Options
	AutoNow    bool
	AutoNowAdd bool
}

func (f TimeField) SQLType() string {
	return "TIME"
}

func (f TimeField) GetOptions() Options {
	return f.Options
}

type DateTimeField struct {
	Options
	AutoNow    bool
	AutoNowAdd bool
}

func (f DateTimeField) SQLType() string {
	return "DATETIME"
}

func (f DateTimeField) GetOptions() Options {
	return f.Options
}

type UUIDField struct {
	Options
}

func (f UUIDField) SQLType() string {
	return "UUID"
}

func (f UUIDField) GetOptions() Options {
	return f.Options
}

type BinaryField struct {
	Options
	MaxLength int
}

func (f BinaryField) SQLType() string {
	if f.MaxLength > 0 {
		return fmt.Sprintf("VARBINARY(%d)", f.MaxLength)
	}
	return "BLOB"
}

func (f BinaryField) GetOptions() Options {
	return f.Options
}

type IPAddressField struct {
	Options
}

func (f IPAddressField) SQLType() string {
	return "VARCHAR(39)" // IPv6 addresses can be up to 39 characters long
}

func (f IPAddressField) GetOptions() Options {
	return f.Options
}

type GenericIPAddressField struct {
	Options
	Protocol string // 'both', 'IPv4', or 'IPv6'
}

func (f GenericIPAddressField) SQLType() string {
	return "VARCHAR(39)" // IPv6 addresses can be up to 39 characters long
}

func (f GenericIPAddressField) GetOptions() Options {
	return f.Options
}

type AutoField struct {
	Options
}

func (f AutoField) SQLType() string {
	return "SERIAL PRIMARY KEY"
}

func (f AutoField) GetOptions() Options {
	return f.Options
}

type ForeignKey struct {
	Options
	RelatedModel string
	OnDelete     string
}

func (f ForeignKey) SQLType() string {
	return fmt.Sprintf("INTEGER REFERENCES %s(id)", f.RelatedModel)
}

func (f ForeignKey) GetOptions() Options {
	return f.Options
}

type ManyToManyField struct {
	Options
	RelatedModel string
	ThroughModel string
}

func (f ManyToManyField) SQLType() string {
	// This is a special case and doesn't directly translate to a SQL type
	return "MANYTOMANY"
}

func (f ManyToManyField) GetOptions() Options {
	return f.Options
}

type OneToOneField struct {
	Options
	RelatedModel string
	OnDelete     string
}

func (f OneToOneField) SQLType() string {
	return fmt.Sprintf("INTEGER UNIQUE REFERENCES %s(id)", f.RelatedModel)
}

func (f OneToOneField) GetOptions() Options {
	return f.Options
}

func CreateSpecialField(name string, fieldType string, options ...func(*Options)) Field {
	field := Options{
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
		return BooleanField{Options: field}
	case "date":
		return DateField{Options: field}
	case "time":
		return TimeField{Options: field}
	case "datetime":
		return DateTimeField{Options: field}
	case "uuid":
		return UUIDField{Options: field}
	case "binary":
		return BinaryField{Options: field}
	case "ip":
		return IPAddressField{Options: field}
	case "genericip":
		return GenericIPAddressField{Options: field}
	case "auto":
		return AutoField{Options: field}
	case "foreignkey":
		return ForeignKey{Options: field}
	case "manytomany":
		return ManyToManyField{Options: field}
	case "onetoone":
		return OneToOneField{Options: field}
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

func WithRelatedModel(relatedModel string) func(interface{}) {
	return func(f interface{}) {
		switch field := f.(type) {
		case *ForeignKey:
			field.RelatedModel = relatedModel
		case *ManyToManyField:
			field.RelatedModel = relatedModel
		case *OneToOneField:
			field.RelatedModel = relatedModel
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
