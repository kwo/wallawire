package repositories

import (
	"context"
	"database/sql"
	neturl "net/url"
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

func BuildPostgresURL(url, cert, key, ca string) string {

	values := neturl.Values{}

	values.Add("sslmode", "verify-full")
	values.Add("sslcert", cert)
	values.Add("sslkey", key)
	values.Add("sslrootcert", ca)

	return url + "?" + values.Encode()

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

func toTimeInteger(value *time.Time) sql.NullInt64 {
	if value == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: value.Unix(),
		Valid: true,
	}
}

func toTimePointer(value sql.NullInt64) *time.Time {
	if !value.Valid {
		return nil
	}
	x := time.Unix(value.Int64, 0).UTC()
	return &x
}
