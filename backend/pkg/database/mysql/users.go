package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"time"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
)

var userTables = []tableDef{
	{
		name: "users",
		columns: []string{
			"id VARCHAR(255) NOT NULL PRIMARY KEY",
		},
	}, {
		name: "user_sessions",
		columns: []string{
			"id VARCHAR(255) NOT NULL",
			"user VARCHAR(255) NOT NULL",
			"access_token VARCHAR(255) NOT NULL",
			"refresh_token VARCHAR(255) NOT NULL",
			"expiration BIGINT NOT NULL",
			"PRIMARY KEY (id)",
			"FOREIGN KEY (user) REFERENCES users(id) ON DELETE CASCADE",
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

func (s *SQL) LookupSession(ID string) (dbtypes.Session, error) {
	query := `
	SELECT
		user, id, access_token, refresh_token, expiration
	FROM
		user_sessions
	WHERE
		id = ?
	`
	s.log.Trace(query)

	var session dbtypes.Session
	var timeUnix int64
	err := s.db.QueryRowContext(s.ctx, query, ID).Scan(
		&session.User,
		&session.ID,
		&session.AccessToken,
		&session.RefreshToken,
		&timeUnix,
	)
	if errorIs(err, errKeyNotFound) {
		return dbtypes.Session{}, fs.ErrNotExist
	} else if err != nil {
		return dbtypes.Session{}, fmt.Errorf("could not lookup session: %v", err)
	}

	session.NotAfter = time.Unix(timeUnix, 0)

	return session, nil
}

func (s *SQL) SetSession(session dbtypes.Session) error {
	if session.User == "" {
		return errors.New("empty user")
	} else if session.ID == "" {
		return errors.New("empty ID")
	} else if session.AccessToken == "" {
		return errors.New("empty access token")
	} else if time.Now().After(session.NotAfter) {
		return errors.New("session expired")
	}

	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	ok, err := s.insertSession(tx, session)
	if err != nil {
		return fmt.Errorf("could not insert session: %v", err)
	} else if !ok {
		err = s.updateSession(tx, session)
		if err != nil {
			return fmt.Errorf("could not update session: %v", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}

func (s *SQL) insertSession(tx *sql.Tx, session dbtypes.Session) (bool, error) {
	query := `
	INSERT INTO
		user_sessions (id, user, access_token, refresh_token, expiration)
	VALUES
		(?, ?, ?, ?, ?)
	`
	s.log.Trace(query)
	argv := []any{
		session.ID,
		session.User,
		session.AccessToken,
		session.RefreshToken,
		session.NotAfter.Unix(),
	}

	_, err := tx.ExecContext(s.ctx, query, argv...)
	if errorIs(err, errKeyExists) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("could not insert session: %v", err)
	}

	return true, nil
}

func (s *SQL) updateSession(tx *sql.Tx, session dbtypes.Session) error {
	query := `
	UPDATE
		user_sessions
	SET
		access_token = ?,
		refresh_token = ?,
		expiration = ?
	WHERE
		id = ?
		AND user = ?
	`
	s.log.Trace(query)

	_, err := tx.ExecContext(s.ctx, query,
		session.AccessToken,
		session.RefreshToken,
		session.NotAfter.Unix(),
		session.ID,
		session.User,
	)
	if err != nil {
		return fmt.Errorf("could not update session: %v", err)
	}

	return nil
}

func (s *SQL) DeleteSession(ID string) error {
	query := "DELETE FROM user_sessions WHERE id = ?"
	s.log.Trace(query)

	_, err := s.db.ExecContext(s.ctx, query, ID)
	if errorIs(err, errKeyNotFound) {
		return nil
	} else if err != nil {
		return fmt.Errorf("could not delete session: %v", err)
	}

	return nil
}

func (s *SQL) PurgeSessions() error {
	query := "DELETE FROM user_sessions WHERE expiration < ?"
	s.log.Trace(query)

	_, err := s.db.ExecContext(s.ctx, query, time.Now().Unix())
	if errorIs(err, errKeyNotFound) {
		return nil
	} else if err != nil {
		return fmt.Errorf("could not purge sessions: %v", err)
	}

	return nil
}
