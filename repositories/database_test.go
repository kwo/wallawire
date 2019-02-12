package repositories_test

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"wallawire/ctxutil"
	"wallawire/repositories"
)

const (
	postgresTestURL  = "postgresql://wallawire@localhost:5432/wallawire"
	postgresUserCert = "../walladata/certs/dbclient/client.wallawire.crt"
	postgresUserKey  = "../walladata/certs/dbclient/client.wallawire.key"
	postgresCACert   = "../walladata/certs/dbclient/ca.crt"
)

var (
	db           *sqlx.DB
	dbLogger     = ctxutil.NewLogger("DatabaseTest", "", nil)
	stmtLock     sync.Mutex
	stmtSetup    []string
	stmtTeardown []string
)

func TestMain(m *testing.M) {
	var rc int
	if err := setup(); err != nil {
		rc = 1
		dbLogger.Error().Err(err).Msg("setup failed")
	} else {
		rc = m.Run()
		teardown()
	}
	os.Exit(rc)
}

func TestFoo(t *testing.T) {
	t.SkipNow()
	if err := db.Ping(); err != nil {
		t.Error(err)
	}
}

func setup() error {
	url := repositories.BuildPostgresURL(postgresTestURL, postgresUserCert, postgresUserKey, postgresCACert)
	x, errOpen := sqlx.Open("postgres", url)
	if errOpen != nil {
		return errOpen
	}
	db = x
	execStatements(stmtTeardown)
	return execStatements(stmtSetup)
}

func teardown() {
	if err := execStatements(stmtTeardown); err != nil {
		dbLogger.Error().Err(err).Msg("teardown statements failed")
	}
	if err := db.Close(); err != nil {
		dbLogger.Error().Err(err).Msg("database close failed")
	}
}

func addTestStatements(setup, teardown []string) {
	stmtLock.Lock()
	defer stmtLock.Unlock()
	for _, stmt := range setup {
		stmtSetup = append(stmtSetup, stmt)
	}
	for _, stmt := range teardown {
		stmtTeardown = append(stmtTeardown, stmt)
	}
}

func execStatements(stmts []string) error {
	for _, stmt := range stmts {
		// dbLogger.Debug().Str("stmt", stmt).Msg("exec")
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}

func compareTimes(t *testing.T, a, b *time.Time) {

	if a == nil && b != nil {
		t.Errorf("nil date, expected %s", b)
	} else if a != nil && b == nil {
		t.Errorf("bad date %s, expected nil", a)
	} else if a != nil && b != nil {
		if !a.Equal(*b) {
			if a.Unix() != b.Unix() {
				t.Errorf("bad unix %d, expected %d", a.Unix(), b.Unix())
			} else if a.Location() != b.Location() {
				t.Errorf("bad tz %v, expected %v", a.Location(), b.Location())
			} else {
				t.Errorf("bad date %s, expected %s", a, b)
			}
		}
	}

	// otherwise, both nil, no error

}
