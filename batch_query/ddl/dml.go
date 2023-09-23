package ddl

import (
	"database/sql"
	"fmt"
	"math/rand"

	bqsql "github.com/koykov/batch_query/mods/sql"
)

const chars = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789"

func ApplyDML(db *sql.DB, maxKey int64, pt bqsql.PlaceholderType) error {
	for i := int64(0); i < maxKey; i++ {
		name := string(randbyte(32))
		status := rand.Intn(1_000_000)
		bio := randbyte(512)
		balance := rand.Float32()
		var pts, pfx string
		switch pt {
		case bqsql.PlaceholderMySQL:
			pts, pfx = "?,?,?,?", ""
		case bqsql.PlaceholderPgSQL:
			pts, pfx = "$1,$2,$3,$4", "bq."
		}
		if _, err := db.Exec(fmt.Sprintf("insert into %susers(name, status, bio, balance) values(%s)", pfx, pts), name, status, bio, balance); err != nil {
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
