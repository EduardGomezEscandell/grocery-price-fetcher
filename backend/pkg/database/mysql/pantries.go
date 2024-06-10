package mysql

import (
	"database/sql"
	"fmt"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
)

func (s *SQL) clearPanties(tx *sql.Tx) error {
	tables := []string{"pantries", "pantry_items"}

	for _, table := range tables {
		q := fmt.Sprintf("DROP TABLE %s", table)
		s.log.Tracef(q)

		_, err := tx.ExecContext(s.ctx, q)
		if err != nil {
			return fmt.Errorf("could not drop table: %v", err)
		}
	}

	return nil
}

func (s *SQL) createPantries(tx *sql.Tx) error {
	queries := []struct {
		name  string
		query string
	}{
		{
			name: "pantries",
			query: `
			CREATE TABLE pantries (
				name VARCHAR(255) PRIMARY KEY
			)`,
		},
		{
			name: "pantry_items",
			query: `
			CREATE TABLE pantry_items (
				pantry_name VARCHAR(255) REFERENCES pantries(name),
				product_name VARCHAR(255) REFERENCES products(name),
				amount FLOAT,
				PRIMARY KEY (pantry_name, product_name)
			)`,
		},
	}

	for _, q := range queries {
		s.log.Trace(q.query)

		_, err := tx.ExecContext(s.ctx, q.query)
		if err != nil && !errorIs(err, errTableExists) {
			return fmt.Errorf("could not create table %s: %v", q.name, err)
		}
	}

	return nil
}

func (s *SQL) Pantries() ([]dbtypes.Pantry, error) {
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	names, err := s.queryPantries(tx)
	if err != nil {
		return nil, fmt.Errorf("could not query pantries: %v", err)
	}

	pantries := make([]dbtypes.Pantry, 0, len(names))
	for _, name := range names {
		contents, err := s.queryPantryContents(tx, name)
		if err != nil {
			return nil, fmt.Errorf("could not query pantry items: %v", err)
		}
		pantries = append(pantries, dbtypes.Pantry{
			Name:     name,
			Contents: contents,
		})
	}

	return pantries, nil
}

func (s *SQL) queryPantries(tx *sql.Tx) ([]string, error) {
	r, err := tx.QueryContext(s.ctx, "SELECT name FROM pantries")
	if err != nil {
		return nil, fmt.Errorf("could not query pantries: %v", err)
	}

	var pantries []string
	for r.Next() {
		var name string
		if err := r.Scan(&name); err != nil {
			return nil, fmt.Errorf("could not scan pantry: %v", err)
		}
		pantries = append(pantries, name)
	}

	return pantries, nil
}

func (s *SQL) queryPantryContents(tx *sql.Tx, name string) ([]dbtypes.Ingredient, error) {
	r, err := tx.QueryContext(s.ctx, "SELECT product_name, amount FROM pantry_items WHERE pantry_name = ?", name)
	if err != nil {
		return nil, fmt.Errorf("could not query pantry items: %v", err)
	}

	var items []dbtypes.Ingredient
	for r.Next() {
		var item dbtypes.Ingredient
		if err := r.Scan(&item.Name, &item.Amount); err != nil {
			return nil, fmt.Errorf("could not scan pantry item: %v", err)
		}
		items = append(items, item)
	}

	return items, nil
}

func (s *SQL) LookupPantry(name string) (dbtypes.Pantry, bool) {
	tx, err := s.db.BeginTx(s.ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return dbtypes.Pantry{}, false
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	row := tx.QueryRowContext(s.ctx, "SELECT name FROM pantries WHERE name = ?", name)
	if err = row.Scan(&name); err != nil {
		return dbtypes.Pantry{}, false
	}

	contents, err := s.queryPantryContents(tx, name)
	if err != nil {
		return dbtypes.Pantry{}, false
	}

	return dbtypes.Pantry{
		Name:     name,
		Contents: contents,
	}, true
}

func (s *SQL) SetPantry(p dbtypes.Pantry) error {
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	_, err = tx.ExecContext(s.ctx, "REPLACE INTO pantries (name) VALUES (?)", p.Name)
	if err != nil {
		return fmt.Errorf("could not insert pantry: %v", err)
	}

	// Remove all items from the pantry
	_, err = tx.ExecContext(s.ctx, "DELETE FROM pantry_items WHERE pantry_name = ?", p.Name)
	if err != nil {
		return fmt.Errorf("could not delete old pantry items: %v", err)
	}

	err = bulkInsert(s, tx,
		"pantry_items (pantry_name, product_name, amount)",
		p.Contents,
		func(i dbtypes.Ingredient) []any {
			return []any{p.Name, i.Name, i.Amount}
		})
	if err != nil {
		return fmt.Errorf("could not insert new pantry items: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}

func (s *SQL) DeletePantry(name string) error {
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	_, err = tx.ExecContext(s.ctx, "DELETE FROM pantries WHERE name = ?", name)
	if err != nil {
		return fmt.Errorf("could not delete pantry: %v", err)
	}

	_, err = tx.ExecContext(s.ctx, "DELETE FROM pantry_items WHERE pantry_name = ?", name)
	if err != nil {
		return fmt.Errorf("could not delete pantry items: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}
