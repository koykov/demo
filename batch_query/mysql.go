package main

import "database/sql"

type MysqlRecord struct {
	//
}

func (r MysqlRecord) Scan(rows *sql.Rows) (any, error) {
	return nil, nil
}

func (r MysqlRecord) Match(key, value any) bool {
	return false
}
