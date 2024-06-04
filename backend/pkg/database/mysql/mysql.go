package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"time"

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

type Settings struct {
	User            string
	Password        string
	Host            string
	Port            string
	ConnectTimeout  time.Duration
	ConnectCooldown time.Duration
}

func DefaultSettings() Settings {
	return Settings{
		User:            "root",
		Password:        "example",
		Host:            "localhost",
		Port:            "3306",
		ConnectTimeout:  time.Minute,
		ConnectCooldown: 5 * time.Second,
	}
}

func New(ctx context.Context, log logger.Logger, sett Settings) (*SQL, error) {
	datasource, err := getDatasource(sett)
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
		s.createPantries(tx),
		s.createShoppingLists(tx),
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

func getDatasource(s Settings) (string, error) {
	return fmt.Sprintf("%s:%s@tcp(%s)/grocery-price-fetcher", s.User, s.Password, net.JoinHostPort(s.Host, s.Port)), nil
}
