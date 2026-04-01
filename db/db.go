package db

import (
	"database/sql"
	_ "embed"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed create-indexes.sql
var createIndexesQuery string

//go:embed create-tables.sql
var createTablesQuery string

//go:embed drop-tables.sql
var dropTablesQuery string

var DB *sql.DB

func init() {
	connStr := "postgresql://postgres:changeme@localhost?sslmode=disable"

	var err error

	if DB, err = sql.Open("pgx", connStr); err != nil {
		log.Fatal("Invalid DB config: ", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("DB unreachable: ", err)
	}

	if err = CreateTables(); err != nil {
		panic(err)
	}

	if err = CreateIndexes(); err != nil {
		panic(err)
	}
}

func CreateTables() error {
	_, err := DB.Exec(createTablesQuery)
	return err
}

func CreateIndexes() error {
	_, err := DB.Exec(createTablesQuery)
	return err
}

func DropTables() error {
	_, err := DB.Exec(dropTablesQuery)
	return err
}

func HardReset() error {
	if err := DropTables(); err != nil {
		return err
	}

	if err := CreateTables(); err != nil {
		return err
	}

	if err := CreateIndexes(); err != nil {
		return err
	}

	return nil
}

// func GraphLookup(storage *storage.Storage, hash string) {
// 	rows, err := DB.Query(`
// 	WITH CTE AS
// 	(
// 	--initialization
// 	SELECT type, cont, hash, parent_hash
// 	FROM objects
// 	WHERE hash = $1
// 	UNION ALL
// 	--recursive execution
// 	SELECT o.type, o.cont, o.hash, o.parent_hash
// 	FROM objects o INNER JOIN objects m
// 	ON o.parent_hash = m.hash
// 	)
// 	SELECT type, cont FROM CTE LIMIT 1000;
// 	`, hash)

// 	if err != nil {
// 		panic(err)
// 	}

// 	storage.ObjectStorage.ConsumeRows(rows)
// }
