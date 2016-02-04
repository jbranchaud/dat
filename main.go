package main

import (
	"fmt"
	"github.com/jackc/pgx"
	"os"
)

var conn *pgx.Conn

func main() {
	var err error
	conn, err = pgx.Connect(extractConfig())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to the database: %v\n", err)
		os.Exit(1)
	}

	rows, _ := conn.Query("select 1")

	var num int32
	rows.Next()
	err = rows.Scan(&num)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to scan query result: %v\n", err)
		os.Exit(1)
	}
	rows.Close()
	fmt.Printf("num: %d\n", num)

	database_name := db_name()
	fmt.Printf("database name: %s\n", database_name)
	fmt.Printf("database size: %s\n", db_size(database_name))
	public_tables := tables_by_schema("public")
	fmt.Printf("public tables: %v\n", public_tables)
}

func extractConfig() pgx.ConnConfig {
	var config pgx.ConnConfig

	config.Host = "localhost"
	// move to a config file or a command-line argument
	config.Database = "hr_hotels"

	return config
}

func db_name() string {
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

func db_size(database_name string) string {
	query := fmt.Sprintf("select pg_size_pretty(pg_database_size('%s'))", database_name)
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

func tables_by_schema(schema_name string) []string {
	query := fmt.Sprintf("select table_name from information_schema.tables where table_schema = '%s'", schema_name)
	rows, err := conn.Query(query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to query database: %v\n", err)
		os.Exit(1)
	}

	var tables []string
	for rows.Next() {
		var table_name string
		err = rows.Scan(&table_name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to scan query result: %v\n", err)
			os.Exit(1)
		}
		tables = append(tables, table_name)
	}

	return tables
}
