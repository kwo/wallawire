package schema

//go:generate go run ../tools/schema/generate.go

import (
	"database/sql"
	"strings"

	_ "github.com/lib/pq"
	"github.com/rubenv/sql-migrate"
)

func Migrate(db *sql.DB, revertLast bool) (int, error) {

	migrations := &migrate.AssetMigrationSource{
		Asset:    getAsset,
		AssetDir: getAssetNames,
	}

	direction := migrate.Up
	numMigrations := 0
	if revertLast {
		direction = migrate.Down
		numMigrations = 1
	}

	return migrate.ExecMax(db, "postgres", migrations, direction, numMigrations)

}

func getAsset(path string) ([]byte, error) {
	return Asset("/" + path), nil
}

func getAssetNames(path string) ([]string, error) {
	var result []string
	for _, name := range AssetNames() {
		result = append(result, strings.TrimPrefix(name, "/"))
	}
	return result, nil
}
