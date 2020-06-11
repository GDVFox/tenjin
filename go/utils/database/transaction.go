package database

import (
	"github.com/gocraft/dbr/v2"
)

// TransactionProcessor function that works with transaction
type TransactionProcessor func(tx *dbr.Tx) error

// InTransaction creates transcations and calls TransactionProcessor
func InTransaction(sess *Session, tp TransactionProcessor) error {
	tx, err := sess.Begin()
	defer tx.RollbackUnlessCommitted()
	if err != nil {
		return err
	}

	if err := tp(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
