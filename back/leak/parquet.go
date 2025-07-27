package leak

import (
	"database/sql"
	"fmt"
)

func (d *Dataleaks) GetParquetColumns(parquetFile string) ([]string, error) {
	if d.Duckdb == nil {
		return nil, fmt.Errorf("DuckDB connection is not initialized")
	}

	query := fmt.Sprintf("DESCRIBE SELECT * FROM '%s';", parquetFile)

	rows, err := d.Duckdb.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query DESCRIBE SELECT * FROM '%s': %w", parquetFile, err)
	}
	defer rows.Close()

	var columns []string

	for rows.Next() {
		var columnName string
		var columnType string
		var nullable string

		var key sql.NullString
		var defaultValue sql.NullString
		var extra sql.NullString

		if err := rows.Scan(&columnName, &columnType, &nullable, &key, &defaultValue, &extra); err != nil {
			return nil, fmt.Errorf("failed to scan column info for '%s': %w", parquetFile, err)
		}
		columns = append(columns, columnName)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over DESCRIBE results for '%s': %w", parquetFile, err)
	}

	if len(columns) == 0 {
		return nil, fmt.Errorf("no columns found for parquet file '%s'. File might be empty, corrupted, or not found", parquetFile)
	}

	return columns, nil
}

func (d *Dataleaks) GetParquetLength(parquetFile string) (uint64, error) {
	if d.Duckdb == nil {
		return 0, fmt.Errorf("DuckDB connection is not initialized")
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM '%s';", parquetFile)

	row := d.Duckdb.QueryRow(query)

	var count uint64
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to get row count for '%s': %w", parquetFile, err)
	}

	return count, nil
}
