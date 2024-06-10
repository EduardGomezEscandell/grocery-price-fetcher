package mysql

import (
	"database/sql"
	"fmt"
	"slices"
	"strings"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
)

func (s *SQL) clearShoppingLists(tx *sql.Tx) error {
	tables := []string{"shopping_list_items"}

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
			name: "shopping_list_items",
			query: `
			CREATE TABLE shopping_list_items (
				menu_name VARCHAR(255) REFERENCES menus(name),
				pantry_name VARCHAR(255) REFERENCES pantries(name),
				product_name VARCHAR(255) REFERENCES products(name),
				PRIMARY KEY (menu_name, pantry_name, product_name)
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

func (s *SQL) ShoppingLists() ([]dbtypes.ShoppingList, error) {
	query := `SELECT * FROM shopping_list_items`
	s.log.Trace(query)

	r, err := s.db.QueryContext(s.ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not query shopping lists: %v", err)
	}

	type item struct {
		Menu, Pantry, Product string
	}

	items := make([]item, 0)
	for r.Next() {
		var i item
		if err := r.Scan(&i.Menu, i.Pantry, i.Product); err != nil {
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
		if i.Menu != j.Menu {
			return strings.Compare(i.Menu, j.Menu)
		}
		if i.Pantry != j.Pantry {
			return strings.Compare(i.Pantry, j.Pantry)
		}
		return strings.Compare(i.Product, j.Product)
	})

	lists := []dbtypes.ShoppingList{
		{
			Menu:     items[0].Menu,
			Pantry:   items[0].Pantry,
			Contents: []string{items[0].Product},
		},
	}

	for i := 1; i < len(items); i++ {
		if items[i].Menu == items[i-1].Menu && items[i].Pantry == items[i-1].Pantry {
			lists[len(lists)-1].Contents = append(lists[len(lists)-1].Contents, items[i].Product)
			continue
		}

		lists = append(lists, dbtypes.ShoppingList{
			Menu:     items[i].Menu,
			Pantry:   items[i].Pantry,
			Contents: []string{items[i].Product},
		})
	}

	return lists, nil
}

func (s *SQL) LookupShoppingList(menu, pantry string) (dbtypes.ShoppingList, bool) {
	sl := dbtypes.ShoppingList{
		Menu:     menu,
		Pantry:   pantry,
		Contents: []string{},
	}

	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		s.log.Warningf("could not begin transaction: %v", err)
		return sl, false
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	query := `
	SELECT product_name
	FROM shopping_list_items
	WHERE menu_name = ? AND pantry_name = ?
	`
	s.log.Trace(query)

	r, err := tx.QueryContext(s.ctx, query, menu, pantry)
	if err != nil {
		s.log.Warningf("could not query shopping list items: %v", err)
		return sl, false
	}

	for r.Next() {
		var item string
		if err := r.Scan(&item); err != nil {
			s.log.Warningf("could not scan: %v", err)
			return sl, false
		}

		sl.Contents = append(sl.Contents, item)
	}

	if err := tx.Commit(); err != nil {
		return sl, false
	}

	return sl, true
}

func (s *SQL) SetShoppingList(list dbtypes.ShoppingList) error {
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	_, err = tx.ExecContext(s.ctx, `
		DELETE FROM
			shopping_list_items
		WHERE 
			menu_name = ? 
			AND pantry_name = ?
	`, list.Menu, list.Pantry)
	if err != nil {
		return fmt.Errorf("could not delete old shopping list items: %v", err)
	}

	err = bulkInsert(s, tx, "shopping_list_items(menu_name, pantry_name, product_name)", list.Contents, func(s string) []any {
		return []any{list.Menu, list.Pantry, s}
	})
	if err != nil {
		return fmt.Errorf("could not insert shopping list items: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}

func (s *SQL) DeleteShoppingList(menu, pantry string) error {
	query := `DELETE FROM shopping_list_items WHERE menu_name = ? AND pantry_name = ?`
	s.log.Trace(query)

	if _, err := s.db.ExecContext(s.ctx, query, menu, pantry); err != nil {
		return fmt.Errorf("could not delete shopping list items: %v", err)
	}

	return nil
}
