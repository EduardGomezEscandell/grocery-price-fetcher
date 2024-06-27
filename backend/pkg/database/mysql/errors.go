package mysql

import (
	"errors"
	"strings"

	driver "github.com/go-sql-driver/mysql"
)

const (
	errTableExists uint16 = 1050
	errKeyExists   uint16 = 1062
	errKeyNotFound uint16 = 1216
)

func errorIs(err error, target uint16) bool {
	if err == nil {
		return false
	}

	if x := (&driver.MySQLError{}); errors.As(err, &x) {
		return x.Number == target
	}

	if target == errKeyNotFound && strings.Contains(err.Error(), "no rows in result set") {
		// MySQL terrible error handling
		return true
	}

	return false
}
