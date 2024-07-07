package mysql

import (
	"database/sql"
	"fmt"
	"io/fs"
	"slices"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/recipe"
)

var menuTables = []tableDef{
	{
		name: "menus",
		columns: []string{
			"user VARCHAR(255)",
			"name VARCHAR(255)",
			"FOREIGN KEY (user) REFERENCES users(id) ON DELETE CASCADE",
			"PRIMARY KEY (user, name)",
		},
	},
	{
		name: "menu_days",
		columns: []string{
			"user VARCHAR(255) NOT NULL",
			"menu VARCHAR(255) NOT NULL",
			"pos INT NOT NULL",
			"name VARCHAR(255) NOT NULL",
			"FOREIGN KEY (user) REFERENCES users(id) ON DELETE CASCADE",
			"FOREIGN KEY (user, menu) REFERENCES menus(user, name) ON DELETE CASCADE",
			"PRIMARY KEY (user, menu, pos)",
		},
	},
	{
		name: "menu_meals",
		columns: []string{
			"user VARCHAR(255) NOT NULL",
			"menu VARCHAR(255) NOT NULL",
			"day INT NOT NULL",
			"pos INT NOT NULL",
			"name VARCHAR(255) NOT NULL",
			"FOREIGN KEY (user) REFERENCES users(id) ON DELETE CASCADE",
			"FOREIGN KEY (user, menu, day) REFERENCES menu_days(user, menu, pos) ON DELETE CASCADE",
			"PRIMARY KEY (user, menu, day, pos)",
		},
	},
	{
		name: "menu_dishes",
		columns: []string{
			"user VARCHAR(255)",
			"menu VARCHAR(255)",
			"day INT NOT NULL",
			"meal INT NOT NULL",
			"pos INT NOT NULL",
			"recipe INT UNSIGNED NOT NULL",
			"amount FLOAT NOT NULL",
			"FOREIGN KEY (user) REFERENCES users(id) ON DELETE CASCADE",
			"FOREIGN KEY (user, menu, day, meal) REFERENCES menu_meals(user, menu, day, pos) ON DELETE CASCADE",
			"FOREIGN KEY (recipe) REFERENCES recipes(id) ON DELETE CASCADE",
			"PRIMARY KEY (user, menu, day, meal, pos)",
		},
	},
}

func (s *SQL) Menus(user string) ([]dbtypes.Menu, error) {
	if user == "" {
		return nil, nil
	}

	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	menus, err := s.queryMenus(tx, user)
	if err != nil {
		return nil, fmt.Errorf("could not query menus: %v", err)
	}

	m, err := s.queryMenuContents(tx, user, menus)
	if err != nil {
		return nil, fmt.Errorf("could not query menu contents: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %v", err)
	}

	return m, nil
}

func (s *SQL) queryMenus(tx *sql.Tx, user string) ([]string, error) {
	rows, err := tx.QueryContext(s.ctx, "SELECT name FROM menus WHERE user = ?", user)
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

func (s *SQL) queryMenuContents(tx *sql.Tx, user string, names []string) ([]dbtypes.Menu, error) {
	if len(names) == 0 {
		return nil, nil
	}

	builder := newMenuBuilder(user, names)

	days, err := s.queryMenuDays(tx, user)
	if err != nil {
		return nil, fmt.Errorf("could not query menu days: %v", err)
	}
	builder.setDays(days)

	meals, err := s.queryMenuMeals(tx, user)
	if err != nil {
		return nil, fmt.Errorf("could not query menu meals: %v", err)
	}
	builder.setMeals(meals)

	items, err := s.queryMealItems(tx, user)
	if err != nil {
		return nil, fmt.Errorf("could not query meal items: %v", err)
	}
	builder.setItems(items)

	return builder.menus, nil
}

type menuDayRow struct {
	User string
	Menu string
	Pos  int
	Name string
}

func (s *SQL) queryMenuDays(tx *sql.Tx, user string) ([]menuDayRow, error) {
	rows, err := tx.QueryContext(s.ctx, "SELECT menu, pos, name FROM menu_days WHERE user = ?", user)
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
	User string
	Menu string
	Day  int
	Pos  int
	Name string
}

func (s *SQL) queryMenuMeals(tx *sql.Tx, user string) ([]menuMealRow, error) {
	query := `
		SELECT
			menu, day, pos, name
		FROM
			menu_meals
		WHERE
			user = ?`

	rows, err := tx.QueryContext(s.ctx, query, user)
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
	User   string
	Menu   string
	Day    int
	Meal   int
	Pos    int
	Recipe recipe.ID
	Amount float32
}

func (s *SQL) queryMealItems(tx *sql.Tx, user string) ([]menuDishRow, error) {
	query := `
	SELECT
		menu, day, meal, pos, recipe, amount 
	FROM 
		menu_dishes
	WHERE user = ?`

	r, err := tx.QueryContext(s.ctx, query, user)
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

func newMenuBuilder(user string, names []string) menuBuilder {
	p := menuBuilder{
		menus: make([]dbtypes.Menu, 0, len(names)),
	}
	for _, n := range names {
		p.menus = append(p.menus, dbtypes.Menu{
			User: user,
			Name: n,
		})
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

func (s *SQL) LookupMenu(user, name string) (dbtypes.Menu, error) {
	if user == "" {
		return dbtypes.Menu{}, fs.ErrNotExist
	}

	tx, err := s.db.BeginTx(s.ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return dbtypes.Menu{}, fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	q := `SELECT name FROM menus WHERE name = ? AND user = ?`
	s.log.Trace(q)

	row := tx.QueryRowContext(s.ctx, q, name, user)
	if err := row.Scan(&name); errorIs(err, errKeyNotFound) {
		return dbtypes.Menu{}, fs.ErrNotExist
	} else if err != nil {
		return dbtypes.Menu{}, fmt.Errorf("could not scan menu name: %v", err)
	}

	if err := row.Err(); err != nil {
		return dbtypes.Menu{}, fmt.Errorf("could not scan menu name: %v", err)
	}

	m, err := s.queryMenuContents(tx, user, []string{name})
	if err != nil {
		return dbtypes.Menu{}, fmt.Errorf("could not query menu contents: %v", err)
	}

	if len(m) == 0 {
		return dbtypes.Menu{}, fs.ErrNotExist
	}

	if err := tx.Commit(); err != nil {
		return dbtypes.Menu{}, fmt.Errorf("could not commit transaction: %v", err)
	}

	return m[0], nil
}

func (s *SQL) SetMenu(m dbtypes.Menu) error {
	if m.User == "" {
		return fmt.Errorf("user cannot be empty")
	}

	if m.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	var dc struct {
		days  []menuDayRow
		meals []menuMealRow
		items []menuDishRow
	}

	for dayIdx, day := range m.Days {
		dc.days = append(dc.days, menuDayRow{
			User: m.User,
			Menu: m.Name,
			Pos:  dayIdx,
			Name: day.Name,
		})

		for mealIdx, meal := range day.Meals {
			dc.meals = append(dc.meals, menuMealRow{
				User: m.User,
				Menu: m.Name,
				Day:  dayIdx,
				Pos:  mealIdx,
				Name: meal.Name,
			})

			for k, dish := range meal.Dishes {
				dc.items = append(dc.items, menuDishRow{
					User:   m.User,
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
	if _, err := tx.ExecContext(s.ctx, `
		DELETE FROM
			menu_days
		WHERE 
			menu = ? AND user = ?
		`, m.Name, m.User); err != nil {
		return fmt.Errorf("could not delete extra meal items: %v", err)
	}

	// Insert new menu from top to bottom
	if err := s.setMenu(tx, m.User, m.Name); err != nil {
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

func (s *SQL) DeleteMenu(user, name string) error {
	if user == "" {
		return nil
	}

	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	if _, err := tx.ExecContext(s.ctx, `
		DELETE FROM
			menus
		WHERE
			name = ? AND user = ?`, name, user); err != nil {
		return fmt.Errorf("could not delete menu: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}

func (s *SQL) setMenu(tx *sql.Tx, user, name string) error {
	// Insert new menu
	_, err := tx.ExecContext(s.ctx, `
		REPLACE INTO
			menus (user, name)
		VALUES (?, ?)`, user, name)
	if err != nil {
		return fmt.Errorf("could not insert menu: %v", err)
	}

	return nil
}

func (s *SQL) setDays(tx *sql.Tx, rows []menuDayRow) error {
	return bulkInsert(s, tx,
		"menu_days (user, menu, pos, name)", rows,
		func(row menuDayRow) []any {
			return []any{row.User, row.Menu, row.Pos, row.Name}
		})
}

func (s *SQL) setMeals(tx *sql.Tx, rows []menuMealRow) error {
	return bulkInsert(s, tx,
		"menu_meals (user, menu, day, pos, name)", rows,
		func(row menuMealRow) []any {
			return []any{row.User, row.Menu, row.Day, row.Pos, row.Name}
		})
}

func (s *SQL) setItems(tx *sql.Tx, rows []menuDishRow) error {
	return bulkInsert(s, tx,
		"menu_dishes (user, menu, day, meal, pos, recipe, amount)", rows,
		func(row menuDishRow) []any {
			return []any{row.User, row.Menu, row.Day, row.Meal, row.Pos, row.Recipe, row.Amount}
		})
}
