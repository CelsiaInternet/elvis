package jdb

import (
	"context"
	"database/sql"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/logs"
	"github.com/celsiainternet/elvis/msg"
)

/**
* Tx wraps a *sql.Tx exposing the same query methods as DB so that
* linq operations can be executed inside a transaction transparently.
**/
type Tx struct {
	tx *sql.Tx
}

/**
* BeginTx starts a new transaction on the DB connection.
* @param ctx context.Context
* @return *Tx, error
**/
func (d *DB) BeginTx(ctx context.Context) (*Tx, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &Tx{tx: tx}, nil
}

/**
* Commit commits the transaction.
* @return error
**/
func (t *Tx) Commit() error {
	return t.tx.Commit()
}

/**
* Rollback aborts the transaction.
* @return error
**/
func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}

func (t *Tx) exec(sql string, args ...any) (*sql.Rows, error) {
	if t == nil {
		return nil, logs.Alertf(msg.NOT_CONNECT_DB)
	}

	return t.tx.QueryContext(context.Background(), sql, args...)
}

/**
* Query executes a SELECT inside the transaction.
* @param sql string, args ...any
* @return et.Items, error
**/
func (t *Tx) Query(sql string, args ...any) (et.Items, error) {
	rows, err := t.exec(sql, args...)
	if err != nil {
		return et.Items{}, err
	}
	defer rows.Close()

	return rowsItems(rows), nil
}

/**
* Command executes a DML statement inside the transaction.
* @param sql string, args ...any
* @return et.Items, error
**/
func (t *Tx) Command(sql string, args ...any) (et.Items, error) {
	rows, err := t.exec(sql, args...)
	if err != nil {
		return et.Items{}, err
	}
	defer rows.Close()

	return rowsItems(rows), nil
}

/**
* CommandSource executes a DML statement inside the transaction and
* extracts the JSONB sourceField from each returned row.
* @param sourceField string, sql string, args ...any
* @return et.Items, error
**/
func (t *Tx) CommandSource(sourceField, sql string, args ...any) (et.Items, error) {
	rows, err := t.exec(sql, args...)
	if err != nil {
		return et.Items{}, err
	}
	defer rows.Close()

	return sourceItems(rows, sourceField), nil
}
