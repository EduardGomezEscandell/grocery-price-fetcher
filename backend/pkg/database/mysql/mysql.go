package mysql

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
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

	allowInsertNewID bool
}

type Settings struct {
	User            string
	PasswordFile    string `yaml:"password-file"`
	Host            string
	Port            string
	ConnectTimeout  time.Duration `yaml:"connect-timeout"`
	ConnectCooldown time.Duration `yaml:"connect-cooldown"`

	// AllowInsertNewID allows the Set* functions to insert new objects with a user-provided ID.
	AllowInsertNewID bool
}

func DefaultSettings() Settings {
	return Settings{
		User:             "root",
		PasswordFile:     "",
		Host:             "localhost",
		Port:             "3306",
		ConnectTimeout:   time.Minute,
		ConnectCooldown:  5 * time.Second,
		AllowInsertNewID: false,
	}
}

func New(ctx context.Context, log logger.Logger, sett Settings) (*SQL, error) {
	datasource, err := getDatasource(sett)
	if err != nil {
		return nil, fmt.Errorf("could not parse options: %w", err)
	}

	log = log.WithField("database", "mysql")

	db, err := sql.Open("mysql", datasource)
	if err != nil {
		return nil, fmt.Errorf("could not open mysql: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	sql := &SQL{
		ctx:              ctx,
		cancel:           cancel,
		log:              log,
		db:               db,
		allowInsertNewID: sett.AllowInsertNewID,
	}

	if err := sql.waitConnection(sett); err != nil {
		sql.Close()
		return nil, fmt.Errorf("could not connect to database: %w", err)
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

func (s *SQL) waitConnection(sett Settings) error {
	ctx, cancel := context.WithTimeout(s.ctx, sett.ConnectTimeout)
	defer cancel()

	tk := time.NewTicker(sett.ConnectCooldown)
	defer tk.Stop()

	// First tick must not wait
	tick := func() <-chan time.Time {
		ch := make(chan time.Time)
		close(ch)
		return ch
	}()

	var i int
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out after %d attempts", i)
		case <-tick:
		}

		// Further ticks will wait
		tick = tk.C
		i++

		s.log.Info("Connecting to database")

		if err := s.db.PingContext(ctx); err != nil {
			s.log.Infof("Attempt %d: could not connect to database: %v", i, err)
			continue
		}

		s.log.Info("Connected to database")
		return nil
	}
}

func (s *SQL) Close() error {
	s.cancel()
	return s.db.Close()
}

func getDatasource(s Settings) (string, error) {
	if s.User == "" {
		return "", errors.New("user is empty")
	}
	if s.PasswordFile == "" {
		return "", errors.New("password is empty")
	}
	if s.Host == "" {
		return "", errors.New("host is empty")
	}
	if s.Port == "" {
		return "", errors.New("port is empty")
	}

	pass, err := os.ReadFile(s.PasswordFile)
	if err != nil {
		return "", fmt.Errorf("could not read password file: %w", err)
	} else if len(pass) == 0 {
		return "", errors.New("password file is empty")
	}

	return fmt.Sprintf(
			"%s:%s@tcp(%s)/grocery-price-fetcher",
			s.User,
			string(bytes.TrimSpace(pass)),
			net.JoinHostPort(s.Host, s.Port)),
		nil
}

// bulkInsert inserts multiple rows into a table.
// INSERT INTO ${into} VALUES (?, ?, ...), (?, ?, ...), ...
// The data is flattened into a single slice of arguments.
// The get function is used to extract the arguments from each element.
//
// DO NOT USE THIS FUNCTION WITH USER INPUT.
func bulkInsert[T any](s *SQL, tx *sql.Tx, into string, data []T, get func(t T) []any) error {
	nRows := len(data)
	if nRows == 0 {
		return nil
	} else if nRows > 1000 {
		// This is a safety measure to prevent accidental large inserts
		return fmt.Errorf("too many rows: %d", nRows)
	}

	//nolint:gosec // The query is constructed by the code, not user input
	query := fmt.Sprintf(`
		INSERT INTO
			%s
		VALUES
			`, into)

	// Get first element to know how many columns to expect
	first := get(data[0])
	nCols := len(first)

	// Build (?, ?, ...) string
	valueStr := fmt.Sprintf("(%s)", repeatStringWithSeparator("?", ", ", nCols))

	// Build (?, ?, ...), (?, ?, ...), ... string
	query += repeatStringWithSeparator(valueStr, ", ", nRows)

	// Flatten data
	argv := make([]any, 0, nCols*nRows)
	argv = append(argv, first...)
	for _, item := range data[1:] {
		argv = append(argv, get(item)...)
	}

	// Insert
	s.log.Trace(query)
	if _, err := tx.ExecContext(s.ctx, query, argv...); err != nil {
		return fmt.Errorf("could not insert: %v", err)
	}

	return nil
}

func repeatStringWithSeparator(str string, sep string, n int) string {
	switch n {
	case 0:
		return ""
	case 1:
		return str
	}

	var b strings.Builder
	b.Grow(len(str)*n + len(sep)*(n-1))
	b.WriteString(str)
	for range n - 1 {
		b.WriteString(sep)
		b.WriteString(str)
	}

	return b.String()
}
