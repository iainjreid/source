// Copyright 2026 Iain J. Reid
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
