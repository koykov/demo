package ddl

import (
	"bytes"
	"database/sql"
	"os"
)

func ApplyMysql(db *sql.DB) error {
	raw, err := os.ReadFile("ddl/mysql.sql")
	if err != nil {
		return err
	}
	scripts := bytes.Split(raw, []byte(";"))
	for _, script := range scripts {
		script = bytes.Trim(script, " \n\t")
		if len(script) == 0 {
			continue
		}
		_, err = db.Exec(string(script))
		if err != nil {
			return err
		}
	}
	return nil
}
