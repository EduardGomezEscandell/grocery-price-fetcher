package mysql

import (
	"database/sql"
	"fmt"
	"slices"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/recipe"
)

func (s *SQL) clearMenus(tx *sql.Tx) error {
	tables := []string{"menus", "menu_days", "menu_meals", "menu_dishes"}

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

func (s *SQL) createMenus(tx *sql.Tx) error {
	queries := []struct {
		name  string
		query string
	}{
		{
			"menus",
			`CREATE TABLE menus (
				name VARCHAR(255) PRIMARY KEY
			)`,
		},
		{
			"menu_days",
			`CREATE TABLE menu_days (
				menu VARCHAR(255) NOT NULL,
				pos INT NOT NULL,
				name VARCHAR(255) NOT NULL,
				FOREIGN KEY (menu) REFERENCES menus(name) ON DELETE CASCADE,
				PRIMARY KEY (menu, pos)
			)`,
		},
		{
			"menu_meals",
			`CREATE TABLE menu_meals (
				menu VARCHAR(255) NOT NULL,
				day INT NOT NULL,
				pos INT NOT NULL,
				name VARCHAR(255) NOT NULL,
				FOREIGN KEY (menu, day) REFERENCES menu_days(menu, pos) ON DELETE CASCADE,
				PRIMARY KEY (menu, day, pos)
			)`,
		},
		{
			"menu_dishes",
			`CREATE TABLE menu_dishes (
				menu VARCHAR(255),
				day INT NOT NULL,
				meal INT NOT NULL,
				pos INT NOT NULL,
				recipe INT UNSIGNED NOT NULL,
				amount FLOAT NOT NULL,
				FOREIGN KEY (menu, day, meal) REFERENCES menu_meals(menu, day, pos) ON DELETE CASCADE,
				FOREIGN KEY (recipe) REFERENCES recipes(id) ON DELETE CASCADE,
				PRIMARY KEY (menu, day, meal, pos)
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

func (s *SQL) Menus() ([]dbtypes.Menu, error) {
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	menus, err := s.queryMenus(tx)
	if err != nil {
		return nil, fmt.Errorf("could not query menus: %v", err)
	}

	m, err := s.queryMenuContents(tx, menus)
	if err != nil {
		return nil, fmt.Errorf("could not query menu contents: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %v", err)
	}

	return m, nil
}

func (s *SQL) queryMenus(tx *sql.Tx) ([]string, error) {
	rows, err := tx.QueryContext(s.ctx, "SELECT name FROM menus")
	if err != nil {
		return nil, fmt.Errorf("could not query menus: %v", err)
	}
	defer rows.Close()

	var menus []string
	for rows.Next() {
		var m string
		if err := rows.Scan(&m); err != nil {
			return nil, fmt.Errorf("could not scan menu: %v", err)
		}
		menus = append(menus, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate over menus: %v", err)
	}

	return menus, nil
}

func (s *SQL) queryMenuContents(tx *sql.Tx, names []string) ([]dbtypes.Menu, error) {
	if len(names) == 0 {
		return nil, nil
	}

	builder := newMenuBuilder(names)

	days, err := s.queryMenuDays(tx)
	if err != nil {
		return nil, fmt.Errorf("could not query menu days: %v", err)
	}
	builder.setDays(days)

	meals, err := s.queryMenuMeals(tx)
	if err != nil {
		return nil, fmt.Errorf("could not query menu meals: %v", err)
	}
	builder.setMeals(meals)

	items, err := s.queryMealItems(tx)
	if err != nil {
		return nil, fmt.Errorf("could not query meal items: %v", err)
	}
	builder.setItems(items)

	return builder.menus, nil
}

type menuDayRow struct {
	Menu string
	Pos  int
	Name string
}

func (s *SQL) queryMenuDays(tx *sql.Tx) ([]menuDayRow, error) {
	rows, err := tx.QueryContext(s.ctx, "SELECT menu, pos, name FROM menu_days")
	if err != nil {
		return nil, fmt.Errorf("could not query menu days: %v", err)
	}
	defer rows.Close()

	var days []menuDayRow
	for rows.Next() {
		var d menuDayRow
		if err := rows.Scan(&d.Menu, &d.Pos, &d.Name); err != nil {
			return nil, fmt.Errorf("could not scan menu day: %v", err)
		}
		days = append(days, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate over menu days: %v", err)
	}

	return days, nil
}

type menuMealRow struct {
	Menu string
	Day  int
	Pos  int
	Name string
}

func (s *SQL) queryMenuMeals(tx *sql.Tx) ([]menuMealRow, error) {
	query := `
		SELECT
			menu, day, pos, name
		FROM
			menu_meals
		`

	rows, err := tx.QueryContext(s.ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not query menu meals: %v", err)
	}
	defer rows.Close()

	var meals []menuMealRow
	for rows.Next() {
		var m menuMealRow
		if err := rows.Scan(&m.Menu, &m.Day, &m.Pos, &m.Name); err != nil {
			return nil, fmt.Errorf("could not scan menu meal: %v", err)
		}
		meals = append(meals, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate over menu meals: %v", err)
	}

	return meals, nil
}

type menuDishRow struct {
	Menu   string
	Day    int
	Meal   int
	Pos    int
	Recipe recipe.ID
	Amount float32
}

func (s *SQL) queryMealItems(tx *sql.Tx) ([]menuDishRow, error) {
	query := `
		SELECT
			menu, day, meal, pos, recipe, amount
		FROM 
			menu_dishes
		`

	r, err := tx.QueryContext(s.ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not query meal items: %v", err)
	}
	defer r.Close()

	var items []menuDishRow
	for r.Next() {
		var i menuDishRow
		if err := r.Scan(&i.Menu, &i.Day, &i.Meal, &i.Pos, &i.Recipe, &i.Amount); err != nil {
			return nil, fmt.Errorf("could not scan meal item: %v", err)
		}
		items = append(items, i)
	}

	if err := r.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate over meal items: %v", err)
	}

	return items, nil
}

type menuBuilder struct {
	menus []dbtypes.Menu
}

func newMenuBuilder(names []string) menuBuilder {
	p := menuBuilder{
		menus: make([]dbtypes.Menu, 0, len(names)),
	}
	for _, n := range names {
		p.menus = append(p.menus, dbtypes.Menu{Name: n})
	}
	return p
}

func (p *menuBuilder) setDays(d []menuDayRow) {
	for _, row := range d {
		menu, ok := getMenu(p.menus, row.Menu)
		if !ok {
			continue
		}

		*at(&menu.Days, row.Pos) = dbtypes.Day{Name: row.Name}
	}
}

func (p *menuBuilder) setMeals(m []menuMealRow) {
	for _, row := range m {
		menu, ok := getMenu(p.menus, row.Menu)
		if !ok {
			continue
		}

		if len(menu.Days) <= row.Day {
			continue
		}

		*at(&menu.Days[row.Day].Meals, row.Pos) = dbtypes.Meal{Name: row.Name}
	}
}

func (p *menuBuilder) setItems(i []menuDishRow) {
	for _, row := range i {
		menu, ok := getMenu(p.menus, row.Menu)
		if !ok {
			continue
		}

		if len(menu.Days) <= row.Day {
			continue
		}

		if len(menu.Days[row.Day].Meals) <= row.Meal {
			continue
		}

		*at(&menu.Days[row.Day].Meals[row.Meal].Dishes, row.Pos) = dbtypes.Dish{
			ID:     row.Recipe,
			Amount: row.Amount,
		}
	}
}

func getMenu(s []dbtypes.Menu, name string) (*dbtypes.Menu, bool) {
	idx := slices.IndexFunc(s, func(v dbtypes.Menu) bool { return v.Name == name })
	if idx == -1 {
		return nil, false
	}
	return &s[idx], true
}

// at returns a pointer to the i-th element of slice, growing it if necessary.
func at[T any](slice *[]T, i int) *T {
	if len(*slice) <= i {
		*slice = append(*slice, make([]T, i-len(*slice)+1)...)
	}
	return &(*slice)[i]
}

func (s *SQL) LookupMenu(name string) (dbtypes.Menu, bool) {
	tx, err := s.db.BeginTx(s.ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		s.log.Warnf("could not begin transaction: %v", err)
		return dbtypes.Menu{}, false
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	q := `SELECT name FROM menus WHERE name = ?`
	s.log.Trace(q)

	row := tx.QueryRowContext(s.ctx, q, name)
	if err := row.Scan(&name); err != nil {
		s.log.Warnf("could not scan menu name: %v", err)
		return dbtypes.Menu{}, false
	}

	if err := row.Err(); err != nil {
		s.log.Warnf("could not scan menu name: %v", err)
		return dbtypes.Menu{}, false
	}

	m, err := s.queryMenuContents(tx, []string{name})
	if err != nil {
		s.log.Warnf("could not query menu contents: %v", err)
		return dbtypes.Menu{}, false
	}

	if len(m) == 0 {
		return dbtypes.Menu{}, false
	}

	if err := tx.Commit(); err != nil {
		s.log.Warnf("could not commit transaction: %v", err)
		return dbtypes.Menu{}, false
	}

	return m[0], true
}

func (s *SQL) SetMenu(m dbtypes.Menu) error {
	var dc struct {
		days  []menuDayRow
		meals []menuMealRow
		items []menuDishRow
	}

	for dayIdx, day := range m.Days {
		dc.days = append(dc.days, menuDayRow{
			Menu: m.Name,
			Pos:  dayIdx,
			Name: day.Name,
		})

		for mealIdx, meal := range day.Meals {
			dc.meals = append(dc.meals, menuMealRow{
				Menu: m.Name,
				Day:  dayIdx,
				Pos:  mealIdx,
				Name: meal.Name,
			})

			for k, dish := range meal.Dishes {
				dc.items = append(dc.items, menuDishRow{
					Menu:   m.Name,
					Day:    dayIdx,
					Meal:   mealIdx,
					Pos:    k,
					Recipe: dish.ID,
					Amount: dish.Amount,
				})
			}
		}
	}

	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	// Delete old menu dependencies via cascade
	if _, err := tx.ExecContext(s.ctx, `DELETE FROM menu_days WHERE menu = ?`, m.Name); err != nil {
		return fmt.Errorf("could not delete extra meal items: %v", err)
	}

	// Insert new menu from top to bottom
	if err := s.setMenu(tx, m.Name); err != nil {
		return fmt.Errorf("could not set menu: %v", err)
	}

	if err := s.setDays(tx, dc.days); err != nil {
		return fmt.Errorf("could not set menu days: %v", err)
	}

	if err := s.setMeals(tx, dc.meals); err != nil {
		return fmt.Errorf("could not set menu meals: %v", err)
	}

	if err := s.setItems(tx, dc.items); err != nil {
		return fmt.Errorf("could not set meal items: %v", err)
	}

	return tx.Commit()
}

func (s *SQL) DeleteMenu(name string) error {
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	if _, err := tx.ExecContext(s.ctx, `DELETE FROM menus WHERE name = ?`, name); err != nil {
		return fmt.Errorf("could not delete menu: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}

func (s *SQL) setMenu(tx *sql.Tx, name string) error {
	// Insert new menu
	_, err := tx.ExecContext(s.ctx, "REPLACE INTO menus (name) VALUES (?)", name)
	if err != nil {
		return fmt.Errorf("could not insert menu: %v", err)
	}

	return nil
}

func (s *SQL) setDays(tx *sql.Tx, rows []menuDayRow) error {
	return bulkInsert(s, tx,
		"menu_days (menu, pos, name)", rows,
		func(row menuDayRow) []any {
			return []any{row.Menu, row.Pos, row.Name}
		})
}

func (s *SQL) setMeals(tx *sql.Tx, rows []menuMealRow) error {
	return bulkInsert(s, tx,
		"menu_meals (menu, day, pos, name)", rows,
		func(row menuMealRow) []any {
			return []any{row.Menu, row.Day, row.Pos, row.Name}
		})
}

func (s *SQL) setItems(tx *sql.Tx, rows []menuDishRow) error {
	return bulkInsert(s, tx,
		"menu_dishes (menu, day, meal, pos, recipe, amount)", rows,
		func(row menuDishRow) []any {
			return []any{row.Menu, row.Day, row.Meal, row.Pos, row.Recipe, row.Amount}
		})
}
