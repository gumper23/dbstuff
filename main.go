package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/fatih/color"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gumper23/dbstuff/dbhelper"
	_ "github.com/lib/pq"
)

func main() {
	hasOutput := false

	dsn, ok := os.LookupEnv("POSTGRES_DSN")
	if ok {
		hasOutput = true
		err := printAllTables("mysql", dsn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}
	}

	dsn, ok = os.LookupEnv("MYSQL_DSN")
	if ok {
		hasOutput = true
		err := printAllTables("mysql", dsn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}
	}

	if !hasOutput {
		fmt.Println("Please set POSTGRES_DSN and/or MYSQL_DSN environment variables")
	}
}

func printAllTables(dbtype, dsn string) (err error) {
	db, err := sql.Open(dbtype, dsn)
	if err != nil {
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return
	}

	tables, _, err := dbhelper.QueryRows(db, "select table_schema as table_schema, table_name as table_name from information_schema.tables")
	if err != nil {
		return
	}

	for _, table := range tables {
		query := fmt.Sprintf("select * from %s.%s limit 10", table["table_schema"], table["table_name"])
		rows, cols, err := dbhelper.QueryRows(db, query)
		if err != nil {
			break
		}
		color.New(color.Bold).Printf("Results of %s:\n", query)
		dbhelper.PrintRows(rows, cols)
		fmt.Println("")
	}
	return
}
