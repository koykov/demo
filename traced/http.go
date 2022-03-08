package traced

import "net/http"

func CheckMethod(r *http.Request, must string) bool {
	return r.Method == must
}
