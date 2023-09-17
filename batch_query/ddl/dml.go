package ddl

import (
	"database/sql"
	"math/rand"
)

const chars = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789"

func ApplyMysqlDML(db *sql.DB, maxKey int64) error {
	for i := int64(0); i < maxKey; i++ {
		name := string(randbyte(32))
		status := rand.Intn(1_000_000)
		bio := randbyte(512)
		balance := rand.Float32()
		if _, err := db.Exec("insert into users(name, status, bio, balance) values(?,?,?,?)", name, status, bio, balance); err != nil {
			return err
		}
	}
	return nil
}

func randbyte(len_ int) []byte {
	buf := make([]byte, 0, len_)
	for i := 0; i < len_; i++ {
		buf = append(buf, chars[rand.Intn(len(chars))])
	}
	return buf
}
