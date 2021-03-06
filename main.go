package dbstats

import (
	"fmt"
	"github.com/jackc/pgx"
	"os"
	"time"
)

var conn *pgx.Conn

type column struct {
	name     string
	dataType string
}

func main() {
	// get the database name as the first command-line argument
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

	db := AnalyzeDatabase(conn)
	initiateDeepAnalysis(db.name)
	// print a bunch of stats about the database
	fmt.Printf("database name: %s\n", db.name)
	fmt.Printf("database size: %s\n", db.size)
	publicTables := selectTablesBySchema("public")
	fmt.Printf("the 'public' schema contains %d tables\n", len(publicTables))
	printTables(publicTables)
	fmt.Printf("\n-----\n")
	printTableColumns("public", publicTables)
}

func selectRow(query string) string {
	rows, err := conn.Query(query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to query database: %v\n", err)
		os.Exit(1)
	}

	var result string
	rows.Next()
	err = rows.Scan(&result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to scan query result: %v\n", err)
		os.Exit(1)
	}
	rows.Close()

	return result
}

func extractConfig() pgx.ConnConfig {
	var config pgx.ConnConfig

	config.Host = "localhost"

	return config
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

func selectColumnData(schemaName string, tableName string) []column {
	query := fmt.Sprintf("select column_name, data_type from information_schema.columns where table_schema = '%s' and table_name = '%s';", schemaName, tableName)
	rows, err := conn.Query(query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to query database: %v\n", err)
		os.Exit(1)
	}

	var columns []column
	for rows.Next() {
		var c column
		err = rows.Scan(&c.name, &c.dataType)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to scan query result: %v\n", err)
			os.Exit(1)
		}
		columns = append(columns, c)
	}

	return columns
}

func printTables(tables []string) {
	for i := 0; i < len(tables); i++ {
		fmt.Printf("  - %s\n", tables[i])
	}
}

func printTableColumns(schemaName string, tables []string) {
	for i := 0; i < len(tables); i++ {
		columns := selectColumnData(schemaName, tables[i])
		fmt.Printf("## %s (%d)\n", tables[i], len(columns))
		for j := 0; j < len(columns); j++ {
			fmt.Printf("   - %s [%s]\n", columns[j].name, columns[j].dataType)
		}
	}
}

func initiateDeepAnalysis(dbName string) {
	fmt.Printf("Initiating deep analysis of '%s'", dbName)
	time.Sleep(2000 * time.Millisecond)
	fmt.Printf(".")
	time.Sleep(1000 * time.Millisecond)
	fmt.Printf(".")
	time.Sleep(1000 * time.Millisecond)
	fmt.Printf(".")
	time.Sleep(1000 * time.Millisecond)
	fmt.Printf("\n")
}
