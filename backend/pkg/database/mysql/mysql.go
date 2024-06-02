package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	// Import the MySQL driver.
	_ "github.com/go-sql-driver/mysql"
)

type SQL struct {
	ctx    context.Context
	cancel context.CancelFunc

	log logger.Logger

	db *sql.DB
}

func DefaultSettings() map[string]interface{} {
	return map[string]interface{}{
		"user":     "root",
		"password": "example",
		"host":     "localhost",
		"port":     "3306",
	}
}

func New(ctx context.Context, log logger.Logger, options map[string]interface{}) (*SQL, error) {
	datasource, err := getDatasource(options)
	if err != nil {
		return nil, fmt.Errorf("could not parse options: %w", err)
	}

	log = log.WithField("database", "mysql")
	log.Tracef("connecting to %s", datasource)

	db, err := sql.Open("mysql", datasource)
	if err != nil {
		return nil, fmt.Errorf("could not open mysql: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	sql := &SQL{
		ctx:    ctx,
		cancel: cancel,
		log:    log,
		db:     db,
	}

	if err := sql.createTables(); err != nil {
		sql.Close()
		return nil, fmt.Errorf("could not create tables: %v", err)
	}

	return sql, nil
}

func (s *SQL) createTables() error {
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	err = errors.Join(
		s.createProducts(tx),
		s.createRecipes(tx),
		s.createMenus(tx),
	)

	if err != nil {
		defer s.Close()
		return fmt.Errorf("could not ensure tables exist: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func (s *SQL) Close() error {
	s.cancel()
	return s.db.Close()
}

func getDatasource(options map[string]any) (string, error) {
	user, err := getStringOption(options, "user")
	if err != nil {
		return "", err
	}

	pass, err := getStringOption(options, "password")
	if err != nil {
		return "", err
	}

	host, err := getStringOption(options, "host")
	if err != nil {
		return "", err
	}

	port, err := getStringOption(options, "port")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s@tcp(%s)/grocery-price-fetcher", user, pass, net.JoinHostPort(host, port)), nil
}

func getStringOption(options map[string]any, key string) (string, error) {
	p, ok := options[key]
	if !ok {
		def, ok := DefaultSettings()[key].(string)
		if !ok {
			return "", fmt.Errorf("option %q not found", key)
		}
		return def, nil
	}

	path, ok := p.(string)
	if !ok {
		return "", fmt.Errorf("option %q is not a string", key)
	}

	return path, nil
}
