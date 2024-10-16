package callbacks

import (
	"github.com/mviner000/eyygo/src/germ"
)

func RawExec(db *germ.DB) {
	if db.Error == nil && !db.DryRun {
		result, err := db.Statement.ConnPool.ExecContext(db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...)
		if err != nil {
			db.AddError(err)
			return
		}

		db.RowsAffected, _ = result.RowsAffected()
	}
}
