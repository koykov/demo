package ddl

import (
	"bytes"
	"database/sql"
	"os"
)

func ApplyMysqlDDL(db *sql.DB, ddlPath string) error {
	raw, err := os.ReadFile(ddlPath)
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

func ApplyPgsqlDDL(db *sql.DB, ddlPath string) error {
	return ApplyMysqlDDL(db, ddlPath)
}
