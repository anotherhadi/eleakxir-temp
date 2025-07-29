package leak

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
)

type Result struct {
	Columns      []string
	Content      []string
	DataleakName string
}

type Query struct {
	Terms      []string
	ExactMatch bool
}

func (d *Dataleaks) GetDataleakFromPath(path string) (*Dataleak, error) {
	for _, leak := range d.Dataleaks {
		if leak.Path == path {
			return &leak, nil
		}
	}
	return nil, fmt.Errorf("dataleak not found for path: %s", path)
}

func (d *Dataleaks) Search(parquetFile string, columns []string, query Query, fulltext bool) ([]Result, error) {
	// Get the dataleak from the path
	dataleak, err := d.GetDataleakFromPath(parquetFile)
	if err != nil {
		return nil, fmt.Errorf("dataleak not found for path %s: %w", parquetFile, err)
	}

	// Filtrer les colonnes valides
	var validColumns []string
	for _, col := range columns {
		if slices.Contains(dataleak.Columns, col) {
			validColumns = append(validColumns, col)
		}
	}
	if len(validColumns) == 0 && !fulltext {
		return []Result{}, nil
	}

	// Build condition
	var condition string
	if fulltext {
		var termClauses []string
		for _, term := range query.Terms {
			termEscaped := strings.ReplaceAll(term, "'", "''")
			termEscaped = strings.ReplaceAll(termEscaped, "_", "\\_")
			termEscaped = strings.ReplaceAll(termEscaped, "%", "\\%")

			concatCols := "lower(" + strings.Join(dataleak.Columns, " || ' ' || ") + ")"
			termClauses = append(termClauses,
				fmt.Sprintf("%s ILIKE '%%%s%%' ESCAPE '\\'", concatCols, strings.ToLower(termEscaped)))
		}
		condition = strings.Join(termClauses, " AND ")
	} else {
		var orClauses []string
		for _, col := range validColumns {
			var andClauses []string
			for _, term := range query.Terms {
				termEscaped := strings.ReplaceAll(term, "'", "''")
				if query.ExactMatch {
					andClauses = append(andClauses, fmt.Sprintf("%s = '%s'", col, termEscaped))
				} else {
					termEscaped := strings.ReplaceAll(termEscaped, "_", "\\_")
					termEscaped = strings.ReplaceAll(termEscaped, "%", "\\%")
					andClauses = append(andClauses, fmt.Sprintf("lower(%s) ILIKE '%%%s%%' ESCAPE '\\'", col, strings.ToLower(termEscaped)))
				}
			}
			orClauses = append(orClauses, "("+strings.Join(andClauses, " AND ")+")")
		}
		condition = strings.Join(orClauses, " OR ")
	}

	// Start the search
	var results []Result

	parquetPath := filepath.Join(d.DataleaksDirectory, parquetFile)
	sqlQuery := fmt.Sprintf("SELECT * FROM '%s' WHERE %s LIMIT 200", parquetPath, condition)
	rows, err := d.Duckdb.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		values := make([]sql.NullString, len(cols))
		ptrs := make([]any, len(cols))
		for i := range values {
			ptrs[i] = &values[i]
		}

		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}

		var content []string
		for _, val := range values {
			if val.Valid {
				content = append(content, val.String)
			} else {
				content = append(content, "")
			}
		}

		results = append(results, Result{
			Columns:      cols,
			Content:      content,
			DataleakName: dataleak.Name,
		})
	}

	results = removeDups(results)

	return results, nil
}

func removeDups(results []Result) []Result {
	seen := make(map[string]bool)
	var uniqueResults []Result

	for _, result := range results {
		key := result.DataleakName + strings.Join(result.Content, "")
		if !seen[key] {
			seen[key] = true
			uniqueResults = append(uniqueResults, result)
		}
	}

	return uniqueResults
}
