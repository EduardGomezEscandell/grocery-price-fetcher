package mysql

import (
	"cmp"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"strings"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/product"
)

func (s *SQL) clearShoppingLists(tx *sql.Tx) error {
	tables := []string{"shopping_list_items"}

	// Remove tables from bottom to top to avoid foreign key constraints
	for i := range tables {
		table := tables[len(tables)-i-1]
		q := fmt.Sprintf("DROP TABLE %s", table)
		s.log.Trace(q)

		_, err := tx.ExecContext(s.ctx, q)
		if err != nil {
			return fmt.Errorf("could not drop table: %v", err)
		}
	}

	return nil
}

func (s *SQL) createShoppingLists(tx *sql.Tx) error {
	queries := []struct {
		name  string
		query string
	}{
		{
			name: "shopping_list_items",
			query: `
			CREATE TABLE shopping_list_items (
				user VARCHAR(255) NOT NULL,
				menu VARCHAR(255) NOT NULL,
				pantry VARCHAR(255) NOT NULL REFERENCES pantries(name) ON DELETE CASCADE,
				product INT UNSIGNED NOT NULL REFERENCES products(id) ON DELETE CASCADE,
				FOREIGN KEY (user, menu) REFERENCES menus(user, name) ON DELETE CASCADE,
				PRIMARY KEY (user, menu, pantry, product)
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

func (s *SQL) ShoppingLists(user string) ([]dbtypes.ShoppingList, error) {
	if user == "" {
		return nil, errors.New("user cannot be empty")
	}

	query := `SELECT menu, pantry, product FROM shopping_list_items WHERE user = ?`
	s.log.Trace(query)

	r, err := s.db.QueryContext(s.ctx, query, user)
	if err != nil {
		return nil, fmt.Errorf("could not query shopping lists: %v", err)
	}

	type item struct {
		Menu, Pantry string
		ProductID    product.ID
	}

	items := make([]item, 0)
	for r.Next() {
		var i item
		if err := r.Scan(&i.Menu, &i.Pantry, &i.ProductID); err != nil {
			return nil, fmt.Errorf("could not scan shopping list: %v", err)
		}

		items = append(items, i)
	}

	if err := r.Err(); err != nil {
		return nil, fmt.Errorf("could not get shopping list items: %v", err)
	}

	if len(items) == 0 {
		return make([]dbtypes.ShoppingList, 0), nil
	}

	slices.SortFunc(items, func(i, j item) int {
		if r := strings.Compare(i.Menu, j.Menu); r != 0 {
			return r
		}
		if r := strings.Compare(i.Pantry, j.Pantry); r != 0 {
			return r
		}
		return cmp.Compare(i.ProductID, j.ProductID)
	})

	lists := []dbtypes.ShoppingList{
		{
			User:     user,
			Menu:     items[0].Menu,
			Pantry:   items[0].Pantry,
			Contents: []product.ID{items[0].ProductID},
		},
	}

	for i := 1; i < len(items); i++ {
		if items[i].Menu == items[i-1].Menu && items[i].Pantry == items[i-1].Pantry {
			lists[len(lists)-1].Contents = append(lists[len(lists)-1].Contents, items[i].ProductID)
			continue
		}

		lists = append(lists, dbtypes.ShoppingList{
			User:     user,
			Menu:     items[i].Menu,
			Pantry:   items[i].Pantry,
			Contents: []product.ID{items[i].ProductID},
		})
	}

	return lists, nil
}

func (s *SQL) LookupShoppingList(user, menu, pantry string) (dbtypes.ShoppingList, error) {
	if user == "" {
		return dbtypes.ShoppingList{}, errors.New("user cannot be empty")
	} else if menu == "" {
		return dbtypes.ShoppingList{}, errors.New("menu cannot be empty")
	} else if pantry == "" {
		return dbtypes.ShoppingList{}, errors.New("pantry cannot be empty")
	}

	sl := dbtypes.ShoppingList{
		User:     user,
		Menu:     menu,
		Pantry:   pantry,
		Contents: make([]product.ID, 0),
	}

	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return sl, fmt.Errorf("could not start transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	query := `
	SELECT
		product
	FROM
		shopping_list_items
	WHERE 
		user = ?
		AND menu = ? 
		AND pantry = ?
	`
	s.log.Trace(query)

	r, err := tx.QueryContext(s.ctx, query, user, menu, pantry)
	if err != nil {
		return sl, fmt.Errorf("could not query shopping list items: %v", err)
	}

	for r.Next() {
		var ID product.ID
		if err := r.Scan(&ID); err != nil {
			return sl, fmt.Errorf("could not scan shopping list item: %v", err)
		}

		sl.Contents = append(sl.Contents, ID)
	}

	if err := r.Err(); err != nil {
		return sl, fmt.Errorf("could not get shopping list items: %v", err)
	}

	if len(sl.Contents) == 0 {
		return sl, fs.ErrNotExist
	}

	if err := tx.Commit(); err != nil {
		return sl, fmt.Errorf("could not commit transaction: %v", err)
	}

	return sl, nil
}

func (s *SQL) SetShoppingList(list dbtypes.ShoppingList) error {
	if list.User == "" {
		return errors.New("user cannot be empty")
	} else if list.Menu == "" {
		return errors.New("menu cannot be empty")
	} else if list.Pantry == "" {
		return errors.New("pantry cannot be empty")
	}

	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	_, err = tx.ExecContext(s.ctx, `
		DELETE FROM
			shopping_list_items
		WHERE 
			user = ?
			AND menu = ? 
			AND pantry = ?
	`, list.User, list.Menu, list.Pantry)
	if err != nil {
		return fmt.Errorf("could not delete old shopping list items: %v", err)
	}

	err = bulkInsert(s, tx, "shopping_list_items(user, menu, pantry, product)", list.Contents, func(ID product.ID) []any {
		return []any{list.User, list.Menu, list.Pantry, ID}
	})
	if err != nil {
		return fmt.Errorf("could not insert shopping list items: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}

func (s *SQL) DeleteShoppingList(user, menu, pantry string) error {
	if user == "" {
		return errors.New("user cannot be empty")
	} else if menu == "" {
		return errors.New("menu cannot be empty")
	} else if pantry == "" {
		return errors.New("pantry cannot be empty")
	}

	query := `DELETE 
		FROM
			shopping_list_items
		WHERE 
			user = ?
			AND menu = ?
			AND pantry = ?`
	s.log.Trace(query)

	if _, err := s.db.ExecContext(s.ctx, query, user, menu, pantry); err != nil {
		return fmt.Errorf("could not delete shopping list items: %v", err)
	}

	return nil
}
