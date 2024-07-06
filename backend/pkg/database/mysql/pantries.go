package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/recipe"
)

func (s *SQL) clearPantries(tx *sql.Tx) error {
	tables := []string{"pantries", "pantry_items"}

	// Remove tables from bottom to top to avoid foreign key constraints
	for i := range tables {
		table := tables[len(tables)-i-1]
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
				user VARCHAR(255) NOT NULL,
				name VARCHAR(255) NOT NULL,
				PRIMARY KEY (user, name)
			)`,
		},
		{
			name: "pantry_items",
			query: `
			CREATE TABLE pantry_items (
				user VARCHAR(255),
				pantry VARCHAR(255),
				product INT UNSIGNED,
				amount FLOAT,
				FOREIGN KEY (user, pantry) REFERENCES pantries(user, name) ON DELETE CASCADE,
				FOREIGN KEY (product) REFERENCES products(id) ON DELETE CASCADE,
				PRIMARY KEY (user, pantry, product)
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

func (s *SQL) Pantries(user string) ([]dbtypes.Pantry, error) {
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	names, err := s.queryPantries(tx, user)
	if err != nil {
		return nil, fmt.Errorf("could not query pantries: %v", err)
	}

	pantries := make([]dbtypes.Pantry, 0, len(names))
	for _, name := range names {
		contents, err := s.queryPantryContents(tx, user, name)
		if err != nil {
			return nil, fmt.Errorf("could not query pantry items: %v", err)
		}
		pantries = append(pantries, dbtypes.Pantry{
			User:     user,
			Name:     name,
			Contents: contents,
		})
	}

	return pantries, nil
}

func (s *SQL) queryPantries(tx *sql.Tx, user string) ([]string, error) {
	r, err := tx.QueryContext(s.ctx, `
		SELECT 
			name
		FROM 
			pantries
		WHERE
			user = ?`, user)
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

func (s *SQL) queryPantryContents(tx *sql.Tx, user, name string) ([]recipe.Ingredient, error) {
	r, err := tx.QueryContext(s.ctx, `
		SELECT
			product, amount
		FROM
			pantry_items
		WHERE
			user = ? AND pantry = ?`, user, name)
	if err != nil {
		return nil, fmt.Errorf("could not query pantry items: %v", err)
	}

	items := make([]recipe.Ingredient, 0)
	for r.Next() {
		var item recipe.Ingredient
		if err := r.Scan(&item.ProductID, &item.Amount); err != nil {
			return nil, fmt.Errorf("could not scan pantry item: %v", err)
		}
		items = append(items, item)
	}

	if err := r.Err(); err != nil {
		return nil, fmt.Errorf("could not get pantry items: %v", err)
	}

	return items, nil
}

func (s *SQL) LookupPantry(user, name string) (dbtypes.Pantry, error) {
	if user == "" {
		return dbtypes.Pantry{}, errors.New("user cannot be empty")
	} else if name == "" {
		return dbtypes.Pantry{}, errors.New("name cannot be empty")
	}

	tx, err := s.db.BeginTx(s.ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return dbtypes.Pantry{}, fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	row := tx.QueryRowContext(s.ctx, `
		SELECT
			COUNT(*)
		FROM
			pantries
		WHERE
			user = ? AND name = ?`, user, name)
	var count int
	if err = row.Scan(&count); err != nil {
		return dbtypes.Pantry{}, fmt.Errorf("could not query pantry %s: %v", name, err)
	} else if count == 0 {
		return dbtypes.Pantry{}, fs.ErrNotExist
	}

	if err := row.Err(); err != nil {
		return dbtypes.Pantry{}, fmt.Errorf("could not get pantry %s: %v", name, err)
	}

	contents, err := s.queryPantryContents(tx, user, name)
	if err != nil {
		return dbtypes.Pantry{}, err
	}

	return dbtypes.Pantry{
		User:     user,
		Name:     name,
		Contents: contents,
	}, nil
}

func (s *SQL) SetPantry(p dbtypes.Pantry) error {
	if p.User == "" {
		return errors.New("user cannot be empty")
	} else if p.Name == "" {
		return errors.New("name cannot be empty")
	}

	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	_, err = tx.ExecContext(s.ctx, "INSERT INTO pantries (user, name) VALUES (?, ?)", p.User, p.Name)
	if err != nil && !errorIs(err, errKeyExists) {
		return fmt.Errorf("could not insert pantry: %v", err)
	}

	// Remove all items from the pantry
	_, err = tx.ExecContext(s.ctx, `
		DELETE FROM
			pantry_items
		WHERE
			user = ?  AND pantry = ?`, p.User, p.Name)
	if err != nil {
		return fmt.Errorf("could not delete old pantry items: %v", err)
	}

	err = bulkInsert(s, tx,
		"pantry_items (user, pantry, product, amount)",
		p.Contents,
		func(i recipe.Ingredient) []any {
			return []any{p.User, p.Name, i.ProductID, i.Amount}
		})
	if err != nil {
		return fmt.Errorf("could not insert new pantry items: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}

func (s *SQL) DeletePantry(user, name string) error {
	if user == "" {
		return errors.New("user cannot be empty")
	} else if name == "" {
		return errors.New("name cannot be empty")
	}

	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	_, err = tx.ExecContext(s.ctx, `
		DELETE FROM
			pantries
		WHERE 
			user = ? AND name = ?`, user, name)
	if err != nil {
		return fmt.Errorf("could not delete pantry: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}
