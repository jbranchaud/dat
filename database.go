package dbstats

import (
	"fmt"
	"github.com/jackc/pgx"
	"os"
)

type Database struct {
	name string
	size string
}

func AnalyzeDatabase(conn *pgx.Conn) *Database {
	db := new(Database)

	db.name = selectDatabaseName(conn)
	db.size = selectDatabaseSize(conn)

	return db
}

func selectDatabaseName(conn *pgx.Conn) string {
	rows, _ := conn.Query("select current_database()")

	var name string
	rows.Next()
	err := rows.Scan(&name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to scan query result: %v\n", err)
		os.Exit(1)
	}
	rows.Close()

	return name
}

func selectDatabaseSize(conn *pgx.Conn) string {
	query := "select pg_size_pretty(pg_database_size(select current_database()))"
	rows, err := conn.Query(query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to query database: %v\n", err)
		os.Exit(1)
	}

	var size string
	rows.Next()
	err = rows.Scan(&size)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to scan query result: %v\n", err)
		os.Exit(1)
	}
	rows.Close()

	return size
}
