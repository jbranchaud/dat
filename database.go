package dbstats

import (
	"github.com/jackc/pgx"
)

type Database struct {
	name string
	size string
}

func AnalyzeDatabase(conn *pgx.Conn) *Database {
	db := new(Database)

	db.name = selectRow("select current_database()")
	db.size = selectRow("select pg_size_pretty(pg_database_size(select current_database()))")

	return db
}
