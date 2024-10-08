package callbacks

import (
	"github.com/mviner000/eyymi/eyygo/germ"
)

func BeginTransaction(db *germ.DB) {
	if !db.Config.SkipDefaultTransaction && db.Error == nil {
		if tx := db.Begin(); tx.Error == nil {
			db.Statement.ConnPool = tx.Statement.ConnPool
			db.InstanceSet("germ:started_transaction", true)
		} else if tx.Error == germ.ErrInvalidTransaction {
			tx.Error = nil
		} else {
			db.Error = tx.Error
		}
	}
}

func CommitOrRollbackTransaction(db *germ.DB) {
	if !db.Config.SkipDefaultTransaction {
		if _, ok := db.InstanceGet("germ:started_transaction"); ok {
			if db.Error != nil {
				db.Rollback()
			} else {
				db.Commit()
			}

			db.Statement.ConnPool = db.ConnPool
		}
	}
}
