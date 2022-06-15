package dmp

import "strconv"

type DMP struct{}

func (d DMP) GetUserName(id int32) string {
	return "u" + strconv.Itoa(int(id))
}
