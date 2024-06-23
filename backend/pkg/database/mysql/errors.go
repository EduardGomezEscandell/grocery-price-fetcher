package mysql

import (
	"errors"

	driver "github.com/go-sql-driver/mysql"
)

const (
	errTableExists uint16 = 1050
	errKeyExists   uint16 = 1062
	errKeyNotFound uint16 = 1216
)

func errorIs(err error, target uint16) bool {
	x := &driver.MySQLError{}
	if errors.As(err, &x) {
		return x.Number == target
	}
	return false
}
