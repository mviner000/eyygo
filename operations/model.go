// operations/model.go
package operations

import (
	"fmt"

	"github.com/mviner000/eyymi/fields"
)

type Model struct {
	TableName string
	Fields    map[string]fields.Field
}

func (m *Model) AddField(name string, field fields.Field) {
	if m.Fields == nil {
		m.Fields = make(map[string]fields.Field)
	}
	m.Fields[name] = field

	// Automatically create and execute the AddField operation
	addFieldOp := AddField{
		ModelName: m.TableName,
		FieldName: name,
		FieldType: field.SQLType(),
	}

	if err := addFieldOp.Execute(); err != nil {
		fmt.Println("Error adding field:", err)
	}
}

func NewModel(tableName string) *Model {
	return &Model{
		TableName: tableName,
		Fields:    make(map[string]fields.Field),
	}
}
