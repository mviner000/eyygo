package fields

import (
	"fmt"
)

type NumericField struct {
	Options    FieldOptions
	Validators []func(int64) error
}

func (f NumericField) GetOptions() FieldOptions {
	return f.Options
}

type IntegerField struct {
	NumericField
}

func (f IntegerField) SQLType() string {
	return "INTEGER"
}

type BigIntegerField struct {
	NumericField
}

func (f BigIntegerField) SQLType() string {
	return "BIGINT"
}

type PositiveIntegerField struct {
	NumericField
}

func (f PositiveIntegerField) SQLType() string {
	return "INTEGER UNSIGNED"
}

type PositiveSmallIntegerField struct {
	NumericField
}

func (f PositiveSmallIntegerField) SQLType() string {
	return "SMALLINT UNSIGNED"
}

type SmallIntegerField struct {
	NumericField
}

func (f SmallIntegerField) SQLType() string {
	return "SMALLINT"
}

type FloatField struct {
	NumericField
}

func (f FloatField) SQLType() string {
	return "FLOAT"
}

type DecimalField struct {
	NumericField
	MaxDigits     int
	DecimalPlaces int
}

func (f DecimalField) SQLType() string {
	return fmt.Sprintf("DECIMAL(%d,%d)", f.MaxDigits, f.DecimalPlaces)
}

type DurationField struct {
	NumericField
}

func (f DurationField) SQLType() string {
	return "BIGINT" // Storing microseconds as a 64-bit integer
}

func CreateNumericField(name string, fieldType string, options ...func(*FieldOptions)) Field {
	field := NumericField{
		Options: FieldOptions{
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
	case "integer":
		return IntegerField{NumericField: field}
	case "biginteger":
		return BigIntegerField{NumericField: field}
	case "positiveinteger":
		return PositiveIntegerField{NumericField: field}
	case "positivesmallinteger":
		return PositiveSmallIntegerField{NumericField: field}
	case "smallinteger":
		return SmallIntegerField{NumericField: field}
	case "float":
		return FloatField{NumericField: field}
	case "decimal":
		return DecimalField{NumericField: field}
	case "duration":
		return DurationField{NumericField: field}
	default:
		return nil
	}
}

func WithNumericValidators(validators ...func(int64) error) func(*NumericField) {
	return func(nf *NumericField) {
		nf.Validators = validators
	}
}

func WithPrecision(maxDigits int, decimalPlaces int) func(*DecimalField) {
	return func(df *DecimalField) {
		df.MaxDigits = maxDigits
		df.DecimalPlaces = decimalPlaces
	}
}
