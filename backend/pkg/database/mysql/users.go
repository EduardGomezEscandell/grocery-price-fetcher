package mysql

import (
	"fmt"
)

var userTables = []tableDef{
	{
		name: "users",
		columns: []string{
			"id VARCHAR(255) NOT NULL PRIMARY KEY",
		},
	},
}

func (s *SQL) LookupUser(id string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE id = ?"
	s.log.Trace(query)

	var count int
	err := s.db.QueryRowContext(s.ctx, query, id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("could not lookup user: %v", err)
	}

	return count > 0, nil
}

func (s *SQL) SetUser(id string) error {
	query := "INSERT INTO users (id) VALUES (?)"
	s.log.Trace(query)

	_, err := s.db.ExecContext(s.ctx, query, id)
	if errorIs(err, errKeyExists) {
		return nil
	} else if err != nil {
		return fmt.Errorf("could not insert user: %v", err)
	}

	return nil
}

func (s *SQL) DeleteUser(id string) error {
	query := "DELETE FROM users WHERE id = ?"
	s.log.Trace(query)

	_, err := s.db.ExecContext(s.ctx, query, id)
	if err != nil {
		return fmt.Errorf("could not delete user: %v", err)
	}

	return nil
}
