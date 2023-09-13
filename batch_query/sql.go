package main

import "database/sql"

type SQLRecord struct {
	ID      int32
	Name    string
	Status  int
	Bio     []byte
	Balance float64
}

func (r SQLRecord) Scan(rows *sql.Rows) (any, error) {
	var rec SQLRecord
	if err := rows.Scan(&rec.ID, &rec.Name, &rec.Status, &rec.Bio, &rec.Balance); err != nil {
		return nil, err
	}
	return rec, nil
}

func (r SQLRecord) Match(key, value any) bool {
	var rec *SQLRecord
	switch value.(type) {
	case SQLRecord:
		t := value.(SQLRecord)
		rec = &t
	case *SQLRecord:
		rec = value.(*SQLRecord)
	default:
		return false
	}

	var id int32
	switch key.(type) {
	case int32:
		id = key.(int32)

		// ...
	}

	return rec.ID == id
}
