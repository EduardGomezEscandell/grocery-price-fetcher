package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"strings"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/blank"
	"github.com/go-sql-driver/mysql"
)

func (s *SQL) clearProducts(tx *sql.Tx) error {
	query := "DROP TABLE products"
	s.log.Trace(query)

	stmt, err := tx.PrepareContext(s.ctx, query)
	if err != nil {
		return fmt.Errorf("could not prepare deletion of products: %v", err)
	}

	if _, err = stmt.ExecContext(s.ctx); err != nil {
		return fmt.Errorf("could not delete table products: %v", err)
	}

	return nil
}

func (s *SQL) createProducts(tx *sql.Tx) error {
	query := `
	CREATE TABLE products (
		id INT UNSIGNED PRIMARY KEY,
		name VARCHAR(255),
		batch_size FLOAT NOT NULL,
		price FLOAT NOT NULL,
		provider VARCHAR(255) NOT NULL,
		provider_id0 VARCHAR(255) NOT NULL,
		provider_id1 VARCHAR(255) NOT NULL,
		provider_id2 VARCHAR(255) NOT NULL
	)`
	s.log.Trace(query)

	stmt, err := tx.PrepareContext(s.ctx, query)
	if err != nil {
		return fmt.Errorf("could not prepare creation of products: %v", err)
	}

	_, err = stmt.ExecContext(s.ctx)
	if err != nil && !errorIs(err, errTableExists) {
		return fmt.Errorf("could not create table products: %v", err)
	}

	return nil
}

func (s *SQL) Products() ([]product.Product, error) {
	query := `
	SELECT 
		id,
		name,
		batch_size,
		price,
		provider,
		provider_id0,
		provider_id1,
		provider_id2
	FROM products
	`
	s.log.Trace(query)

	r, err := s.db.QueryContext(s.ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not query products: %v", err)
	}

	products := make([]product.Product, 0)
	for r.Next() {
		p, err := parseProduct(s.log, r)
		if err != nil {
			s.log.Warningf("could not parse product: %v", err)
			continue
		}

		products = append(products, p)
	}

	if err := r.Err(); err != nil {
		return nil, fmt.Errorf("could not get products: %v", err)
	}

	return products, nil
}

func (s *SQL) LookupProduct(ID uint32) (product.Product, error) {
	query := `
	SELECT
		id,
		name,
		batch_size,
		price,
		provider,
		provider_id0,
		provider_id1,
		provider_id2
	FROM products
	WHERE id = ?
	`
	s.log.Trace(query)

	r := s.db.QueryRowContext(s.ctx, query, ID)
	p, err := parseProduct(s.log, r)
	if err != nil {
		return product.Product{}, err
	}

	if err := r.Err(); err != nil {
		return p, fmt.Errorf("could not get product %d: %v", ID, err)
	}

	return p, nil
}

func (s *SQL) SetProduct(p product.Product) (uint32, error) {
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("could not start transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	verb := "REPLACE"
	if p.ID == 0 {
		// Generate a new IDs until we find one that doesn't exist
		// We never expect to have anywhere near 2^32 (4.3 billion) products
		// Collisions are extremely unlikely, but taken care of with the loop
		verb = "INSERT"
		p.ID = rand.Uint32() //nolint:gosec // We don't need a secure random number here
	}

	for {
		//nolint:gosec // This is safe because both halves of the query are hardcoded
		query := verb + ` INTO
			products
			(id, name, batch_size, price, provider, provider_id0, provider_id1, provider_id2)
		VALUES 
			(?, ?, ?, ?, ?, ?, ?, ?)
		`
		s.log.Trace(query)

		argv := []any{p.ID, p.Name, p.BatchSize, p.Price, p.Provider.Name(), p.ProductCode[0], p.ProductCode[1], p.ProductCode[2]}

		_, err := tx.ExecContext(s.ctx, query, argv...)
		if err == nil {
			// Success
			break
		}

		target := (&mysql.MySQLError{})
		if errors.As(err, &target) && target.Number == errKeyExists {
			// Key conflict: generate a new ID
			//nolint:gosec // We don't need a secure random number here
			p.ID = rand.Uint32()
			continue
		}

		// Some other error
		return 0, fmt.Errorf("could not insert into table products: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("could not commit transaction: %v", err)
	}

	return p.ID, nil
}

func (s *SQL) DeleteProduct(ID uint32) error {
	query := `DELETE FROM products WHERE id = ?`
	s.log.Trace(query)

	stmt, err := s.db.PrepareContext(s.ctx, query)
	if err != nil {
		return fmt.Errorf("could not prepare removal of product %d: %v", ID, err)
	}

	if _, err = stmt.ExecContext(s.ctx, ID); err != nil {
		return fmt.Errorf("could not remove product %d: %v", ID, err)
	}

	return nil
}

func parseProduct(log logger.Logger, r interface{ Scan(...any) error }) (p product.Product, err error) {
	var provider string
	var productCode [3]string

	err = r.Scan(&p.ID, &p.Name, &p.BatchSize, &p.Price, &provider, &productCode[0], &productCode[1], &productCode[2])
	if errorIs(err, errKeyNotFound) {
		return p, fs.ErrNotExist
	} else if err != nil && strings.Contains(err.Error(), "sql: no rows in result set") { // MySQL terrible error handling
		return p, fs.ErrNotExist
	} else if err != nil {
		return p, fmt.Errorf("could not scan product: %v", err)
	}

	if prov, ok := providers.Lookup(provider); !ok {
		log.Warningf("could not find provider %q", provider)
		p.Provider = blank.Provider{}
	} else {
		p.Provider = prov
	}

	if err = p.Provider.ValidateCode(productCode); err != nil {
		log.Warningf("Provider %s: could not validate product ID: %v", p.Provider.Name(), err)
		p.Provider = blank.Provider{}
	} else {
		p.ProductCode = productCode
	}

	return p, nil
}
