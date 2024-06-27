package mysql

import (
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/recipe"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/utils"
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
				id INT UNSIGNED PRIMARY KEY,
				name VARCHAR(255) NOT NULL
			)`,
		},
		{
			"recipe_ingredients",
			`CREATE TABLE recipe_ingredients (
				recipe_id INT UNSIGNED REFERENCES recipes(id),
				ingredient_id INT UNSIGNED REFERENCES ingredients(id),
				amount FLOAT NOT NULL,
				PRIMARY KEY (recipe_id, ingredient_id)
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

func (s *SQL) Recipes() ([]recipe.Recipe, error) {
	tx, err := s.db.BeginTx(s.ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	recs, err := s.queryRecipes(tx)
	if err != nil {
		return nil, fmt.Errorf("could not query recipes: %v", err)
	}

	for i := range recs {
		err := s.queryIngredients(tx, &recs[i])
		if err != nil {
			s.log.Warningf("could not get recipe %d %s: %v", recs[i].ID, recs[i].Name, err)
			recs[i].ID = 0 // Mark as invalid
			continue
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %v", err)
	}

	// Filter out recipes marked as invalid
	p := utils.Partition(recs, func(r recipe.Recipe) bool { return r.ID != 0 })
	return recs[:p], nil
}

func (s *SQL) LookupRecipe(id recipe.ID) (recipe.Recipe, error) {
	tx, err := s.db.BeginTx(s.ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return recipe.Recipe{}, fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	rec := recipe.Recipe{
		ID: id,
	}

	row := tx.QueryRowContext(s.ctx, "SELECT name FROM recipes WHERE id = ?", id)
	if err := row.Scan(&rec.Name); errorIs(err, errKeyNotFound) {
		return rec, fs.ErrNotExist
	} else if err != nil {
		return rec, fmt.Errorf("could not scan: %v", err)
	}

	if err := row.Err(); err != nil {
		s.log.Warningf("could not get recipe %d: %v", id, err)
		return rec, fmt.Errorf("could not get recipe %d: %v", id, err)
	}

	if err := s.queryIngredients(tx, &rec); err != nil {
		return rec, fmt.Errorf("could not get recipe %d %s: %v", rec.ID, rec.Name, err)
	}

	if err := tx.Commit(); err != nil {
		return rec, fmt.Errorf("could not commit transaction: %v", err)
	}

	return rec, nil
}

func (s *SQL) queryRecipes(tx *sql.Tx) ([]recipe.Recipe, error) {
	query := `SELECT id, name FROM recipes`

	s.log.Tracef(query)
	r, err := tx.QueryContext(s.ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not query: %v", err)
	}
	defer r.Close()

	var recs []recipe.Recipe
	for r.Next() {
		var rec recipe.Recipe
		if err := r.Scan(&rec.ID, &rec.Name); err != nil {
			s.log.Warnf("could not scan: %v", err)
			continue
		}
		recs = append(recs, rec)
	}

	if err := r.Err(); err != nil {
		return nil, fmt.Errorf("could not get recipes: %v", err)
	}

	return recs, nil
}

func (s *SQL) queryIngredients(tx *sql.Tx, rec *recipe.Recipe) error {
	query := `
	SELECT
		ingredient_id, amount
	FROM
		recipe_ingredients
	WHERE
		recipe_id = ?
	`

	s.log.Tracef(query)
	ingr, err := tx.QueryContext(s.ctx, query, rec.ID)
	if err != nil {
		return fmt.Errorf("could not query ingredients: %v", err)
	}
	defer ingr.Close()

	rec.Ingredients = make([]recipe.Ingredient, 0)
	for ingr.Next() {
		var i recipe.Ingredient
		if err := ingr.Scan(&i.ProductID, &i.Amount); err != nil {
			return fmt.Errorf("could not scan ingredients: %v", err)
		}
		rec.Ingredients = append(rec.Ingredients, i)
	}

	if err := ingr.Err(); err != nil {
		return fmt.Errorf("could not get ingredients: %v", err)
	}

	return nil
}

func (s *SQL) SetRecipe(r recipe.Recipe) (recipe.ID, error) {
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	verb := "REPLACE"
	if r.ID == 0 {
		// Generate a new IDs until we find one that doesn't exist
		// We never expect to have anywhere near 2^32 (4.3 billion) recipes
		// Collisions are extremely unlikely, but taken care of with the loop
		verb = "INSERT"
		r.ID = recipe.NewRandomID()
	}

	for {
		//nolint:gosec // This is safe because both halves of the query are hardcoded
		queryRecipe := verb + ` INTO recipes (id, name) VALUES (?, ?)`
		s.log.Tracef(queryRecipe)

		_, err := tx.ExecContext(s.ctx, queryRecipe, r.ID, r.Name)
		if err == nil {
			// Success
			break
		}

		if errorIs(err, errKeyExists) {
			// Key conflict: generate a new ID
			r.ID = recipe.NewRandomID()
			continue
		}

		// Some other error
		return 0, fmt.Errorf("could not insert recipe: %v", err)
	}

	if err := s.deleteRecipeIngredients(tx, r.ID); err != nil {
		return 0, fmt.Errorf("could not delete old ingredients: %v", err)
	}

	err = bulkInsert(s, tx,
		"recipe_ingredients(recipe_id, ingredient_id, amount)",
		r.Ingredients, func(i recipe.Ingredient) []interface{} {
			return []interface{}{r.ID, i.ProductID, i.Amount}
		})
	if err != nil {
		return 0, fmt.Errorf("could not insert ingredients: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("could not commit transaction: %v", err)
	}

	return r.ID, nil
}

func (s *SQL) DeleteRecipe(id recipe.ID) error {
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	if err := s.deleteRecipe(tx, id); err != nil {
		return fmt.Errorf("could not delete recipe: %v", err)
	}

	if err := s.deleteRecipeIngredients(tx, id); err != nil {
		return fmt.Errorf("could not delete ingredients: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}

func (s *SQL) deleteRecipe(tx *sql.Tx, id recipe.ID) error {
	query := `DELETE FROM recipes WHERE id = ?`
	s.log.Tracef(query)

	_, err := tx.ExecContext(s.ctx, query, id)
	if err != nil {
		return fmt.Errorf("could not delete product: %v", err)
	}

	return nil
}

func (s *SQL) deleteRecipeIngredients(tx *sql.Tx, id recipe.ID) error {
	query := `DELETE FROM recipe_ingredients WHERE recipe_id = ?`
	s.log.Tracef(query)

	_, err := tx.ExecContext(s.ctx, query, id)
	if err != nil {
		return fmt.Errorf("could not delete ingredients: %v", err)
	}

	return nil
}
