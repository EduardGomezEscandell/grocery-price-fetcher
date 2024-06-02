package mysql

import (
	"errors"

	driver "github.com/go-sql-driver/mysql"
)

const (
	errTableExists uint16 = 1050
)

//nolint:unparam // This function is used to check if a MySQL error is of a certain type, but this may change in the future.
func errorIs(err error, target uint16) bool {
	x := &driver.MySQLError{}
	if errors.As(err, &x) {
		return x.Number == target
	}
	return false
}
