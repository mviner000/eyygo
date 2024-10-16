// migrations/detect.go
package migrations

import (
	"database/sql"
	"fmt"

	"github.com/mviner000/eyygo/src/operations"
)

func DetectChanges(model *operations.Model, db *sql.DB) ([]operations.Operation, error) {
	var ops []operations.Operation

	// Check if the table exists
	var tableName string
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?;", model.TableName).Scan(&tableName)
	if err == sql.ErrNoRows {
		// Table does not exist, create it
		ops = append(ops, &operations.AddTable{Model: model})
	} else if err != nil {
		return nil, err
	} else {
		// Table exists, check for field differences
		for _, field := range model.Fields {
			columnExists := false
			rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s);", model.TableName))
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			for rows.Next() {
				var cid int
				var name, ctype string
				var notnull, dfltValue, pk interface{}
				if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err != nil {
					return nil, err
				}
				if name == field.GetOptions().Name {
					columnExists = true
					break
				}
			}

			if !columnExists {
				ops = append(ops, &operations.AddField{
					ModelName: model.TableName,
					FieldName: field.GetOptions().Name,
					FieldType: field.SQLType(),
				})
			}
		}
	}

	return ops, nil
}
