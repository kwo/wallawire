package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"wallawire/model"
)

func NewDatabase(db *sqlx.DB) model.Database {
	return &database{
		db: db,
	}
}

type database struct {
	db *sqlx.DB
}

func (z *database) Run(fn func(tx model.Transaction) error) error {

	ctx := context.Background()
	tx, errBegin := z.db.BeginTxx(ctx, nil)
	if errBegin != nil {
		return errBegin
	}

	if err := fn(&transaction{tx: tx}); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()

}

type transaction struct {
	tx *sqlx.Tx
}

func (z *transaction) Exec(query string, params map[string]interface{}) (model.Result, error) {
	return z.tx.NamedExec(query, params)
}

func (z *transaction) Query(query string, params map[string]interface{}) (model.Rows, error) {
	return z.tx.NamedQuery(query, params)
}

func toNullString(value string) sql.NullString {
	if len(strings.TrimSpace(value)) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: value,
		Valid:  true,
	}
}

// converts a time pointer to sql null integer, returning a null sql integer if the time is nil or IsZero
func toNullTimeInteger(value *time.Time) sql.NullInt64 {
	if value == nil || value.IsZero() {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: value.Unix(),
		Valid: true,
	}
}

// converts an sql null integer to a time struct, returning a zero time if the sql value is null
func toTime(value sql.NullInt64) time.Time {
	if !value.Valid {
		return time.Time{}
	}
	return time.Unix(value.Int64, 0).UTC()
}

// converts a time to a time pointer, returning nil if the time IsZero
func toTimePointer(value time.Time) *time.Time {
	if value.IsZero() {
		return nil
	}
	return &value
}
