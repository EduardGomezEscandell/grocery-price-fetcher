package mysql

import (
	"database/sql"
	"fmt"
	"slices"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
)

func (s *SQL) clearMenus(tx *sql.Tx) error {
	tables := []string{"menus", "menu_days", "menu_day_meals", "menu_day_meal_recipes"}

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
				menu_name VARCHAR(255) REFERENCES menu(name),
				pos INT,
				name VARCHAR(255) NOT NULL,
				PRIMARY KEY (menu_name, pos)
			)`,
		},
		{
			"menu_day_meals",
			`CREATE TABLE menu_day_meals (
				menu_name VARCHAR(255) REFERENCES menu(name),
				day_pos INT REFERENCES menu_days(pos),
				pos INT,
				name VARCHAR(255) NOT NULL,
				PRIMARY KEY (menu_name, day_pos, pos)
			)`,
		},
		{
			"menu_meal_recipes",
			`CREATE TABLE menu_day_meal_recipes (
				menu_name VARCHAR(255) REFERENCES menu(name),
				day_pos INT REFERENCES menu_days(pos),
				meal_pos INT REFERENCES menu_day_meals(pos),
				pos INT,
				recipe_name VARCHAR(255) REFERENCES recipes(name),
				amount FLOAT NOT NULL,
				PRIMARY KEY (menu_name, day_pos, meal_pos, pos)
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

type dayMenuRow struct {
	Menu string
	Pos  int
	Name string
}

func (s *SQL) queryMenuDays(tx *sql.Tx) ([]dayMenuRow, error) {
	rows, err := tx.QueryContext(s.ctx, "SELECT menu_name, pos, name FROM menu_days")
	if err != nil {
		return nil, fmt.Errorf("could not query menu days: %v", err)
	}
	defer rows.Close()

	var days []dayMenuRow
	for rows.Next() {
		var d dayMenuRow
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

type mealMenuRow struct {
	Menu   string
	DayPos int
	Pos    int
	Name   string
}

func (s *SQL) queryMenuMeals(tx *sql.Tx) ([]mealMenuRow, error) {
	rows, err := tx.QueryContext(s.ctx, "SELECT menu_name, day_pos, pos, name FROM menu_day_meals")
	if err != nil {
		return nil, fmt.Errorf("could not query menu meals: %v", err)
	}
	defer rows.Close()

	var meals []mealMenuRow
	for rows.Next() {
		var m mealMenuRow
		if err := rows.Scan(&m.Menu, &m.DayPos, &m.Pos, &m.Name); err != nil {
			return nil, fmt.Errorf("could not scan menu meal: %v", err)
		}
		meals = append(meals, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate over menu meals: %v", err)
	}

	return meals, nil
}

type mealItemRow struct {
	Menu    string
	DayPos  int
	MealPos int
	Pos     int
	Recipe  string
	Amount  float32
}

func (s *SQL) queryMealItems(tx *sql.Tx) ([]mealItemRow, error) {
	query := `
	SELECT
		menu_name, day_pos, meal_pos, pos, recipe_name, amount 
	FROM 
		menu_day_meal_recipes`

	r, err := tx.QueryContext(s.ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not query meal items: %v", err)
	}
	defer r.Close()

	var items []mealItemRow
	for r.Next() {
		var i mealItemRow
		if err := r.Scan(&i.Menu, &i.DayPos, &i.MealPos, &i.Pos, &i.Recipe, &i.Amount); err != nil {
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

func (p *menuBuilder) setDays(d []dayMenuRow) {
	for _, row := range d {
		menu, ok := getMenu(p.menus, row.Menu)
		if !ok {
			continue
		}

		*at(&menu.Days, row.Pos) = dbtypes.Day{Name: row.Name}
	}
}

func (p *menuBuilder) setMeals(m []mealMenuRow) {
	for _, row := range m {
		menu, ok := getMenu(p.menus, row.Menu)
		if !ok {
			continue
		}

		if len(menu.Days) <= row.DayPos {
			continue
		}

		*at(&menu.Days[row.DayPos].Meals, row.Pos) = dbtypes.Meal{Name: row.Name}
	}
}

func (p *menuBuilder) setItems(i []mealItemRow) {
	for _, row := range i {
		menu, ok := getMenu(p.menus, row.Menu)
		if !ok {
			continue
		}

		if len(menu.Days) <= row.DayPos {
			continue
		}

		if len(menu.Days[row.DayPos].Meals) <= row.MealPos {
			continue
		}

		*at(&menu.Days[row.DayPos].Meals[row.MealPos].Dishes, row.Pos) = dbtypes.Dish{
			Name:   row.Recipe,
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
		days  []dayMenuRow
		meals []mealMenuRow
		items []mealItemRow
	}

	for i, day := range m.Days {
		dc.days = append(dc.days, dayMenuRow{
			Menu: m.Name,
			Pos:  i,
			Name: day.Name,
		})

		for j, meal := range day.Meals {
			dc.meals = append(dc.meals, mealMenuRow{
				Menu:   m.Name,
				DayPos: i,
				Pos:    j,
				Name:   meal.Name,
			})

			for k, dish := range meal.Dishes {
				dc.items = append(dc.items, mealItemRow{
					Menu:    m.Name,
					DayPos:  i,
					MealPos: j,
					Pos:     k,
					Recipe:  dish.Name,
					Amount:  dish.Amount,
				})
			}
		}
	}

	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	if err := s.deleteMenuDependencies(tx, m.Name); err != nil {
		return fmt.Errorf("could not delete old menu dependencies: %v", err)
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

	if err := s.deleteMenuDependencies(tx, name); err != nil {
		return fmt.Errorf("could not delete menu dependencies: %v", err)
	}

	if _, err := tx.ExecContext(s.ctx, `DELETE FROM menus WHERE name = ?`, name); err != nil {
		return fmt.Errorf("could not delete menu: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}

func (s *SQL) deleteMenuDependencies(tx *sql.Tx, name string) error {
	// Remove old menu from bottom to top
	if _, err := tx.ExecContext(s.ctx, `DELETE FROM menu_day_meal_recipes WHERE menu_name = ?`, name); err != nil {
		return fmt.Errorf("could not delete extra meal items: %v", err)
	}

	if _, err := tx.ExecContext(s.ctx, `DELETE FROM menu_day_meals WHERE menu_name = ?`, name); err != nil {
		return fmt.Errorf("could not delete extra menu meals: %v", err)
	}

	if _, err := tx.ExecContext(s.ctx, `DELETE FROM menu_days WHERE menu_name = ?`, name); err != nil {
		return fmt.Errorf("could not delete old menu days: %v", err)
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

func (s *SQL) setDays(tx *sql.Tx, rows []dayMenuRow) error {
	return bulkInsert(s, tx,
		"menu_days (menu_name, pos, name)", rows,
		func(row dayMenuRow) []any {
			return []any{row.Menu, row.Pos, row.Name}
		})
}

func (s *SQL) setMeals(tx *sql.Tx, rows []mealMenuRow) error {
	return bulkInsert(s, tx,
		"menu_day_meals (menu_name, day_pos, pos, name)", rows,
		func(row mealMenuRow) []any {
			return []any{row.Menu, row.DayPos, row.Pos, row.Name}
		})
}

func (s *SQL) setItems(tx *sql.Tx, rows []mealItemRow) error {
	return bulkInsert(s, tx,
		"menu_day_meal_recipes (menu_name, day_pos, meal_pos, pos, recipe_name, amount)", rows,
		func(row mealItemRow) []any {
			return []any{row.Menu, row.DayPos, row.MealPos, row.Pos, row.Recipe, row.Amount}
		})
}
