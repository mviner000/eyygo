package sqlite

import (
	"encoding/json"

	"github.com/mviner000/eyymi/src/germ"
)

// The error codes to map sqlite errors to gorm errors, here is a reference about error codes for sqlite https://www.sqlite.org/rescode.html.
var errCodes = map[int]error{
	1555: germ.ErrDuplicatedKey,
	2067: germ.ErrDuplicatedKey,
	787:  germ.ErrForeignKeyViolated,
}

type ErrMessage struct {
	Code         int `json:"Code"`
	ExtendedCode int `json:"ExtendedCode"`
	SystemErrno  int `json:"SystemErrno"`
}

// Translate it will translate the error to native germ errors.
// We are not using go-sqlite3 error type intentionally here because it will need the CGO_ENABLED=1 and cross-C-compiler.
func (dialector Dialector) Translate(err error) error {
	parsedErr, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		return err
	}

	var errMsg ErrMessage
	unmarshalErr := json.Unmarshal(parsedErr, &errMsg)
	if unmarshalErr != nil {
		return err
	}

	if translatedErr, found := errCodes[errMsg.ExtendedCode]; found {
		return translatedErr
	}
	return err
}
