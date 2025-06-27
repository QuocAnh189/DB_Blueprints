package database

import "database/sql"

type ScanFunc[T any] func(*sql.Rows) (*T, error)

func ScanRows[T any](rows *sql.Rows, scan ScanFunc[T]) ([]*T, error) {
	defer rows.Close()

	var results []*T

	for rows.Next() {
		item, err := scan(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
