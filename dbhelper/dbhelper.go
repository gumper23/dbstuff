package dbhelper

import (
	"database/sql"
	"fmt"
	"os"
	"text/tabwriter"
)

// QueryRows executes a query and returns n rows of results
// Each row is a map of column_names to query values as strings
func QueryRows(db *sql.DB, query string) (results []map[string]string, cols []string, err error) {
	rows, err := db.Query(query)
	if err != nil {
		return
	}
	defer rows.Close()

	cols, err = rows.Columns()
	if err != nil {
		return
	}

	rawResult := make([][]byte, len(cols))
	dest := make([]interface{}, len(cols))
	for i := range rawResult {
		dest[i] = &rawResult[i]
	}

	for rows.Next() {
		row := make(map[string]string, len(cols))
		err = rows.Scan(dest...)
		if err != nil {
			return
		}

		for i, raw := range rawResult {
			if raw == nil {
				row[cols[i]] = "NULL"
			} else {
				row[cols[i]] = string(raw)
			}
		}
		results = append(results, row)
	}
	err = rows.Err()
	return
}

// QueryRow executes a query and returns the first row
// The row is a map of column_names to query values
// 0 rows returns sql.ErrNoRows
func QueryRow(db *sql.DB, query string) (row map[string]string, cols []string, err error) {
	rows, cols, err := QueryRows(db, query)
	if err != nil {
		return
	} else if len(rows) == 0 {
		err = sql.ErrNoRows
		return
	} else {
		return rows[0], cols, nil
	}
}

// PrintRows prints all rows in order of cols
func PrintRows(rows []map[string]string, cols []string) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', tabwriter.Debug)
	// Header
	for i, col := range cols {
		if i == 0 {
			col = "|" + col
		}
		fmt.Fprintf(w, "%s\t", col)
	}
	fmt.Fprintf(w, "\n")

	// Data
	for _, row := range rows {
		for i, col := range cols {
			if i == 0 {
				row[col] = "|" + row[col]
			}
			fmt.Fprintf(w, "%s\t", row[col])
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
}
