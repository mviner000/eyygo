package callbacks

import (
	"github.com/mviner000/eyymi/eyygo/germ"
)

var (
	createClauses = []string{"INSERT", "VALUES", "ON CONFLICT"}
	queryClauses  = []string{"SELECT", "FROM", "WHERE", "GROUP BY", "ORDER BY", "LIMIT", "FOR"}
	updateClauses = []string{"UPDATE", "SET", "WHERE"}
	deleteClauses = []string{"DELETE", "FROM", "WHERE"}
)

type Config struct {
	LastInsertIDReversed bool
	CreateClauses        []string
	QueryClauses         []string
	UpdateClauses        []string
	DeleteClauses        []string
}

func RegisterDefaultCallbacks(db *germ.DB, config *Config) {
	enableTransaction := func(db *germ.DB) bool {
		return !db.SkipDefaultTransaction
	}

	if len(config.CreateClauses) == 0 {
		config.CreateClauses = createClauses
	}
	if len(config.QueryClauses) == 0 {
		config.QueryClauses = queryClauses
	}
	if len(config.DeleteClauses) == 0 {
		config.DeleteClauses = deleteClauses
	}
	if len(config.UpdateClauses) == 0 {
		config.UpdateClauses = updateClauses
	}

	createCallback := db.Callback().Create()
	createCallback.Match(enableTransaction).Register("germ:begin_transaction", BeginTransaction)
	createCallback.Register("germ:before_create", BeforeCreate)
	createCallback.Register("germ:save_before_associations", SaveBeforeAssociations(true))
	createCallback.Register("germ:create", Create(config))
	createCallback.Register("germ:save_after_associations", SaveAfterAssociations(true))
	createCallback.Register("germ:after_create", AfterCreate)
	createCallback.Match(enableTransaction).Register("germ:commit_or_rollback_transaction", CommitOrRollbackTransaction)
	createCallback.Clauses = config.CreateClauses

	queryCallback := db.Callback().Query()
	queryCallback.Register("germ:query", Query)
	queryCallback.Register("germ:preload", Preload)
	queryCallback.Register("germ:after_query", AfterQuery)
	queryCallback.Clauses = config.QueryClauses

	deleteCallback := db.Callback().Delete()
	deleteCallback.Match(enableTransaction).Register("germ:begin_transaction", BeginTransaction)
	deleteCallback.Register("germ:before_delete", BeforeDelete)
	deleteCallback.Register("germ:delete_before_associations", DeleteBeforeAssociations)
	deleteCallback.Register("germ:delete", Delete(config))
	deleteCallback.Register("germ:after_delete", AfterDelete)
	deleteCallback.Match(enableTransaction).Register("germ:commit_or_rollback_transaction", CommitOrRollbackTransaction)
	deleteCallback.Clauses = config.DeleteClauses

	updateCallback := db.Callback().Update()
	updateCallback.Match(enableTransaction).Register("germ:begin_transaction", BeginTransaction)
	updateCallback.Register("germ:setup_reflect_value", SetupUpdateReflectValue)
	updateCallback.Register("germ:before_update", BeforeUpdate)
	updateCallback.Register("germ:save_before_associations", SaveBeforeAssociations(false))
	updateCallback.Register("germ:update", Update(config))
	updateCallback.Register("germ:save_after_associations", SaveAfterAssociations(false))
	updateCallback.Register("germ:after_update", AfterUpdate)
	updateCallback.Match(enableTransaction).Register("germ:commit_or_rollback_transaction", CommitOrRollbackTransaction)
	updateCallback.Clauses = config.UpdateClauses

	rowCallback := db.Callback().Row()
	rowCallback.Register("germ:row", RowQuery)
	rowCallback.Clauses = config.QueryClauses

	rawCallback := db.Callback().Raw()
	rawCallback.Register("germ:raw", RawExec)
	rawCallback.Clauses = config.QueryClauses
}
