package model

type Database interface {
	Run(func(tx Transaction) error) error
}

type ReadOnlyTransaction interface {
	Query(query string, params map[string]interface{}) (Rows, error)
}

type WriteOnlyTransaction interface {
	Exec(query string, params map[string]interface{}) (Result, error)
}

type Transaction interface {
	ReadOnlyTransaction
	WriteOnlyTransaction
}

type Rows interface {
	Close() error
	Err() error
	Next() bool
	Scan(dest ...interface{}) error
	StructScan(dest interface{}) error
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}
