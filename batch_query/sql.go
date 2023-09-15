package main

import (
	"database/sql"
	"strconv"
)

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
	case int:
		id = int32(key.(int))
	case int8:
		id = int32(key.(int8))
	case int16:
		id = int32(key.(int16))
	case int32:
		id = key.(int32)
	case int64:
		id = int32(key.(int64))
	case uint:
		id = int32(key.(uint))
	case uint8:
		id = int32(key.(uint8))
	case uint16:
		id = int32(key.(uint16))
	case uint32:
		id = int32(key.(uint32))
	case uint64:
		id = int32(key.(uint64))
	case string:
		id64, err := strconv.ParseInt(key.(string), 10, 64)
		if err != nil {
			return false
		}
		id = int32(id64)
	case []byte:
		id64, err := strconv.ParseInt(string(key.([]byte)), 10, 64)
		if err != nil {
			return false
		}
		id = int32(id64)
	default:
		return false
	}

	return rec.ID == id
}
