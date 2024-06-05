package mysql

import (
	"database/sql"
	"fmt"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/blank"
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
		name VARCHAR(255) PRIMARY KEY,
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

	var products []product.Product
	for r.Next() {
		p, err := parseProduct(s.log, r)
		if err != nil {
			s.log.Warningf("could not parse product: %v", err)
			continue
		}

		products = append(products, p)
	}

	return products, nil
}

func (s *SQL) LookupProduct(name string) (product.Product, bool) {
	query := `
	SELECT
		name,
		batch_size,
		price,
		provider,
		provider_id0,
		provider_id1,
		provider_id2
	FROM products
	WHERE name = ?
	`
	s.log.Trace(query)

	r := s.db.QueryRowContext(s.ctx, query, name)
	p, err := parseProduct(s.log, r)
	if err != nil {
		return product.Product{}, false
	}

	return p, true
}

func (s *SQL) SetProduct(p product.Product) error {
	query := `
	REPLACE INTO
		products
		(name, batch_size, price, provider, provider_id0, provider_id1, provider_id2)
	VALUES 
		(?, ?, ?, ?, ?, ?, ?)
	`
	s.log.Trace(query)

	argv := []any{p.Name, p.BatchSize, p.Price, p.Provider.Name(), p.ProductID[0], p.ProductID[1], p.ProductID[2]}

	if _, err := s.db.ExecContext(s.ctx, query, argv...); err != nil {
		return fmt.Errorf("could not insert into table products: %v", err)
	}

	return nil
}

func (s *SQL) DeleteProduct(name string) error {
	query := `DELETE FROM products WHERE name = ?`
	s.log.Trace(query)

	stmt, err := s.db.PrepareContext(s.ctx, query)
	if err != nil {
		return fmt.Errorf("could not prepare removal of product %s: %v", name, err)
	}

	if _, err = stmt.ExecContext(s.ctx, name); err != nil {
		return fmt.Errorf("could not remove product %s: %v", name, err)
	}

	return nil
}

func parseProduct(log logger.Logger, r interface{ Scan(...any) error }) (p product.Product, err error) {
	var provider string
	var providerID [3]string

	err = r.Scan(&p.Name, &p.BatchSize, &p.Price, &provider, &providerID[0], &providerID[1], &providerID[2])
	if err != nil {
		return p, fmt.Errorf("could not scan product: %v", err)
	}

	if prov, ok := providers.Lookup(provider); !ok {
		log.Warningf("could not find provider %q", provider)
		p.Provider = blank.Provider{}
	} else {
		p.Provider = prov
	}

	if err = p.Provider.ValidateID(providerID); err != nil {
		log.Warningf("Provider %s: could not validate product ID: %v", p.Provider.Name(), err)
		p.Provider = blank.Provider{}
	} else {
		p.ProductID = providerID
	}

	return p, nil
}
