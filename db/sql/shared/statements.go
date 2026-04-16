package shared

import (
	_ "embed"
)

//go:embed create-indexes.sql
var CreateIndexesQuery string

//go:embed create-tables.sql
var CreateTablesQuery string

//go:embed drop-tables.sql
var DropTablesQuery string
