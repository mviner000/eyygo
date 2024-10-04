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

func (m *Model) AddField(field fields.Field) {
	if m.Fields == nil {
		m.Fields = make(map[string]fields.Field)
	}
	m.Fields[field.GetOptions().Name] = field
}

func NewModel(tableName string) *Model {
	return &Model{
		TableName: tableName,
		Fields:    make(map[string]fields.Field),
	}
}

func (m *Model) CreateTableSQL() string {
	fieldDefs := []string{}
	for _, field := range m.Fields {
		fieldDefs = append(fieldDefs, fields.FieldDefinition(field))
	}
	return fmt.Sprintf("CREATE TABLE %s (%s);", m.TableName, fieldDefs)
}

func (m *Model) GenerateAddFieldOperations(existingSchema map[string]map[string]string) []AddField {
	var ops []AddField
	for name, field := range m.Fields {
		if _, exists := existingSchema[m.TableName][name]; !exists {
			ops = append(ops, AddField{
				ModelName: m.TableName,
				FieldName: name,
				FieldType: field.SQLType(),
			})
		}
	}
	return ops
}
