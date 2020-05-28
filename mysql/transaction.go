// Package mysql transaction
package mysql

import (
	"database/sql"
)

// TxDB ...
type TxDB struct {
	MDB *sql.DB
}

// Update ...
func (t *TxDB) Update(fn func(tx *sql.Tx) error) error {
	tx, err := t.MDB.Begin()
	if err != nil {
		return err
	}
	err = fn(tx)
	if err != nil {
		err = tx.Rollback()
	} else {
		err = tx.Commit()
	}
	return err
}
