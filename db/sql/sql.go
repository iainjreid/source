package sql

import (
	"database/sql"
	"log"

	"github.com/iainjreid/source/db/sql/shared"
	_ "github.com/jackc/pgx/v5/stdlib"
)

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

	if err = createTables(); err != nil {
		panic(err)
	}

	if err = createIndexes(); err != nil {
		panic(err)
	}
}

func createTables() error {
	_, err := DB.Exec(shared.CreateTablesQuery)
	return err
}

func createIndexes() error {
	_, err := DB.Exec(shared.CreateIndexesQuery)
	return err
}

func dropTables() error {
	_, err := DB.Exec(shared.DropTablesQuery)
	return err
}

func HardReset() error {
	if err := dropTables(); err != nil {
		return err
	}

	if err := createTables(); err != nil {
		return err
	}

	if err := createIndexes(); err != nil {
		return err
	}

	return nil
}
