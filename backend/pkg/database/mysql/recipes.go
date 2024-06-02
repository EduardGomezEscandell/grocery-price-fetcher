package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/types"
)

func (s *SQL) clearRecipes(tx *sql.Tx) (err error) {
	tables := []string{"recipe_ingredients", "recipes"}

	for _, table := range tables {
		q := fmt.Sprintf("DROP TABLE %s", table)
		s.log.Tracef(q)

		_, err = tx.ExecContext(s.ctx, q)
		if err != nil {
			return fmt.Errorf("could not drop table: %v", err)
		}
	}

	return nil
}

func (s *SQL) createRecipes(tx *sql.Tx) error {
	queries := []struct {
		name  string
		query string
	}{
		{
			"recipes",
			`CREATE TABLE recipes (
				name VARCHAR(255) PRIMARY KEY
			)`,
		},
		{
			"recipe_ingredients",
			`CREATE TABLE recipe_ingredients (
				recipe_name VARCHAR(255) REFERENCES recipes(name),
				ingredient_name VARCHAR(255) REFERENCES ingredients(name),
				amount FLOAT NOT NULL,
				PRIMARY KEY (recipe_name, ingredient_name)
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

func (s *SQL) Recipes() (recipes []types.Recipe, err error) {
	tx, err := s.db.BeginTx(s.ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	recs, err := s.queryRecipes(tx)
	if err != nil {
		return nil, fmt.Errorf("could not query recipes: %v", err)
	}

	for _, name := range recs {
		rec, err := s.queryIngredients(tx, name)
		if err != nil {
			s.log.Warningf("could not get recipe %s: %v", name, err)
			continue
		}
		recipes = append(recipes, rec)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %v", err)
	}

	return recipes, nil
}

func (s *SQL) LookupRecipe(name string) (types.Recipe, bool) {
	tx, err := s.db.BeginTx(s.ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		s.log.Errorf("could not begin transaction: %v", err)
		return types.Recipe{}, false
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	row := tx.QueryRowContext(s.ctx, "SELECT name FROM recipes WHERE name = ?", name)
	if err := row.Scan(&name); err != nil {
		return types.Recipe{}, false
	}

	rec, err := s.queryIngredients(tx, name)
	if err != nil {
		s.log.Warningf("could not get recipe %s: %v", name, err)
		return types.Recipe{}, false
	}

	if err := tx.Commit(); err != nil {
		s.log.Errorf("could not commit transaction: %v", err)
		return types.Recipe{}, false
	}

	return rec, true
}

func (s *SQL) queryRecipes(tx *sql.Tx) ([]string, error) {
	query := `SELECT name FROM recipes`

	s.log.Tracef(query)
	r, err := tx.QueryContext(s.ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not query: %v", err)
	}
	defer r.Close()

	var recs []string
	for r.Next() {
		var rec types.Recipe
		if err := r.Scan(&rec.Name); err != nil {
			s.log.Warnf("could not scan: %v", err)
			continue
		}
		recs = append(recs, rec.Name)
	}

	return recs, nil
}

func (s *SQL) queryIngredients(tx *sql.Tx, recipe string) (types.Recipe, error) {
	rec := types.Recipe{Name: recipe}

	query := `
	SELECT
		recipe_name, ingredient_name, amount
	FROM
		recipe_ingredients
	WHERE
		recipe_name = ?
	`

	s.log.Tracef(query)
	ingr, err := tx.QueryContext(s.ctx, query, recipe)
	if err != nil {
		return rec, fmt.Errorf("could not query ingredients: %v", err)
	}
	defer ingr.Close()

	for ingr.Next() {
		var i types.Ingredient
		var dummy string
		if err := ingr.Scan(&dummy, &i.Name, &i.Amount); err != nil {
			return rec, fmt.Errorf("could not scan ingredients: %v", err)
		}
		rec.Ingredients = append(rec.Ingredients, i)
	}

	return rec, nil
}

func (s *SQL) SetRecipe(r types.Recipe) error {
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	queryRecipe := `
	REPLACE INTO
		recipes (name)
	VALUES
		(?)
	`
	s.log.Tracef(queryRecipe)
	if _, err := tx.ExecContext(s.ctx, queryRecipe, r.Name); err != nil {
		return fmt.Errorf("could not insert recipe: %v", err)
	}

	if len(r.Ingredients) > 1000 {
		// This is an arbitrary limit to prevent abuse
		return fmt.Errorf("too many ingredients")
	}

	if len(r.Ingredients) != 0 {
		//nolint:gosec // The query is constructed by the code, not user input
		queryIngredients := `
		REPLACE INTO
		recipe_ingredients (recipe_name, ingredient_name, amount)
		VALUES
		` + repeatString("(?, ?, ?)", ", ", len(r.Ingredients))

		s.log.Tracef(queryIngredients)

		var argv []any
		for _, i := range r.Ingredients {
			argv = append(argv, r.Name, i.Name, i.Amount)
		}

		if _, err := tx.ExecContext(s.ctx, queryIngredients, argv...); err != nil {
			return fmt.Errorf("could not insert ingredients: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}

func (s *SQL) DeleteRecipe(name string) error {
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	if err := s.deleteRecipe(tx, name); err != nil {
		return fmt.Errorf("could not delete recipe: %v", err)
	}

	if err := s.deleteRecipeIngredients(tx, name); err != nil {
		return fmt.Errorf("could not delete ingredients: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}

func (s *SQL) deleteRecipe(tx *sql.Tx, name string) error {
	query := `DELETE FROM recipes WHERE name = ?`
	s.log.Tracef(query)

	_, err := tx.ExecContext(s.ctx, query, name)
	if err != nil {
		return fmt.Errorf("could not delete product: %v", err)
	}

	return nil
}

func (s *SQL) deleteRecipeIngredients(tx *sql.Tx, name string) error {
	query := `DELETE FROM recipe_ingredients WHERE recipe_name = ?`
	s.log.Tracef(query)

	_, err := tx.ExecContext(s.ctx, query, name)
	if err != nil {
		return fmt.Errorf("could not delete ingredients: %v", err)
	}

	return nil
}

func repeatString(str string, sep string, n int) string {
	switch n {
	case 0:
		return ""
	case 1:
		return str
	}

	var b strings.Builder
	b.Grow(len(str)*n + len(sep)*(n-1))
	b.WriteString(str)
	for range n - 1 {
		b.WriteString(sep)
		b.WriteString(str)
	}

	return b.String()
}
