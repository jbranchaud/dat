package main

import (
	"fmt"
	"github.com/jackc/pgx"
	"os"
)

var conn *pgx.Conn

func main() {
	// get the database as the first command-line argument
	if len(os.Args) <= 1 {
		fmt.Fprintf(os.Stderr, "Expected the name of a database as a command-line argument, but got nothing\n")
		os.Exit(1)
	}
	database_name := os.Args[1]

	pgxConfig := extractConfig()
	pgxConfig.Database = database_name

	var err error
	conn, err = pgx.Connect(pgxConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to the database: %v\n", err)
		os.Exit(1)
	}

	dbName := selectDatabaseName()
	fmt.Printf("database name: %s\n", dbName)
	fmt.Printf("database size: %s\n", selectDatabaseSize(dbName))
	publicTables := selectTablesBySchema("public")
	fmt.Printf("public tables: %v\n", publicTables)
}

func extractConfig() pgx.ConnConfig {
	var config pgx.ConnConfig

	config.Host = "localhost"

	return config
}

func selectDatabaseName() string {
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

func selectDatabaseSize(databaseName string) string {
	query := fmt.Sprintf("select pg_size_pretty(pg_database_size('%s'))", databaseName)
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

func selectTablesBySchema(schemaName string) []string {
	query := fmt.Sprintf("select table_name from information_schema.tables where table_schema = '%s'", schemaName)
	rows, err := conn.Query(query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to query database: %v\n", err)
		os.Exit(1)
	}

	var tables []string
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to scan query result: %v\n", err)
			os.Exit(1)
		}
		tables = append(tables, tableName)
	}

	return tables
}
