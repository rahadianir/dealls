package dbhelper

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/rahadianir/dealls/internal/pkg/xerror"
)

type KeyTransaction string

const TXKey KeyTransaction = "sql-database-transaction"

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	QueryxContext(context.Context, string, ...interface{}) (*sqlx.Rows, error)
	QueryRowxContext(context.Context, string, ...interface{}) *sqlx.Row
}

func InjectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, TXKey, tx)
}

func ExtractTx(ctx context.Context, dbConn *sqlx.DB) (tx DBTX) {
	tx, ok := ctx.Value(TXKey).(*sqlx.Tx)
	if !ok {
		return dbConn
	}

	return tx
}

func WithTransaction(ctx context.Context, dbConn *sqlx.DB, txfunc func(context.Context) error) error {
	tx, err := dbConn.BeginTx(ctx, nil)
	if err != nil {
		return xerror.ServerError{Err: err}
	}
	defer tx.Rollback()

	err = txfunc(InjectTx(ctx, tx))
	if err != nil {
		return err
	}
	return tx.Commit()
}
