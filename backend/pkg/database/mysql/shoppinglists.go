package mysql

import (
	"database/sql"
	"fmt"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/types"
)

func (s *SQL) clearShoppingLists(tx *sql.Tx) error {
	tables := []string{"shopping_lists", "shopping_list_items"}

	for _, table := range tables {
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
			name: "shopping_lists",
			query: `
			CREATE TABLE shopping_lists (
				name VARCHAR(255) PRIMARY KEY,
				timestamp VARCHAR(255) NOT NULL
			)`,
		},
		{
			name: "shopping_list_items",
			query: `
			CREATE TABLE shopping_list_items (
				shopping_list_name VARCHAR(255) REFERENCES shopping_lists(name),
				product_name VARCHAR(255) REFERENCES products(name),
				PRIMARY KEY (shopping_list_name, product_name)
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

func (s *SQL) ShoppingLists() ([]types.ShoppingList, error) {
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not start transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	lists, err := s.queryShoppingLists(tx)
	if err != nil {
		return nil, fmt.Errorf("could not query shopping lists: %v", err)
	}

	for i := range lists {
		lists[i].Items, err = s.shoppingListItems(tx, lists[i].Name)
		if err != nil {
			return nil, fmt.Errorf("could not get items for shopping list %s: %v", lists[i].Name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %v", err)
	}

	return lists, nil
}

func (s *SQL) queryShoppingLists(tx *sql.Tx) ([]types.ShoppingList, error) {
	query := `SELECT name, timestamp FROM shopping_lists`
	s.log.Trace(query)

	r, err := tx.QueryContext(s.ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not query shopping lists: %v", err)
	}

	lists := make([]types.ShoppingList, 0)
	for r.Next() {
		var sl types.ShoppingList
		if err := r.Scan(&sl.Name, &sl.TimeStamp); err != nil {
			return nil, fmt.Errorf("could not scan shopping list: %v", err)
		}

		lists = append(lists, sl)
	}

	return lists, nil
}

func (s *SQL) shoppingListItems(tx *sql.Tx, name string) ([]string, error) {
	query := `
	SELECT product_name
	FROM shopping_list_items
	WHERE shopping_list_name = ?
	`
	s.log.Trace(query)

	r, err := tx.QueryContext(s.ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("could not query shopping list items: %v", err)
	}

	var items []string
	for r.Next() {
		var item string
		if err := r.Scan(&item); err != nil {
			return nil, fmt.Errorf("could not scan shopping list item: %v", err)
		}

		items = append(items, item)
	}

	return items, nil
}

func (s *SQL) LookupShoppingList(name string) (types.ShoppingList, bool) {
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return types.ShoppingList{}, false
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	query := `SELECT name, timestamp FROM shopping_lists WHERE name = ?`
	s.log.Trace(query)

	var sl types.ShoppingList
	row := tx.QueryRowContext(s.ctx, query, name)
	if err := row.Scan(&sl.Name, &sl.TimeStamp); err != nil {
		return types.ShoppingList{}, false
	}

	if sl.Items, err = s.shoppingListItems(tx, name); err != nil {
		return types.ShoppingList{}, false
	}

	if err := tx.Commit(); err != nil {
		return types.ShoppingList{}, false
	}

	return sl, true
}

func (s *SQL) SetShoppingList(list types.ShoppingList) error {
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	query := `REPLACE INTO shopping_lists (name, timestamp) VALUES (?, ?)`
	s.log.Trace(query)

	_, err = tx.ExecContext(s.ctx, query, list.Name, list.TimeStamp)
	if err != nil {
		return fmt.Errorf("could not insert shopping list: %v", err)
	}

	for _, item := range list.Items {
		query := `REPLACE INTO shopping_list_items (shopping_list_name, product_name) VALUES (?, ?)`
		s.log.Trace(query)

		_, err = tx.ExecContext(s.ctx, query, list.Name, item)
		if err != nil {
			return fmt.Errorf("could not insert shopping list item: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}

func (s *SQL) DeleteShoppingList(name string) error {
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	query := `DELETE FROM shopping_lists WHERE name = ?`
	s.log.Trace(query)

	_, err = tx.ExecContext(s.ctx, query, name)
	if err != nil {
		return fmt.Errorf("could not delete shopping list: %v", err)
	}

	query = `DELETE FROM shopping_list_items WHERE shopping_list_name = ?`
	s.log.Trace(query)

	_, err = tx.ExecContext(s.ctx, query, name)
	if err != nil {
		return fmt.Errorf("could not delete shopping list items: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}
