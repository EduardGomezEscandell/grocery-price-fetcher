package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/recipe"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/utils"
)

var recipeTables = []tableDef{
	{
		name: "recipes",
		columns: []string{
			"id INT UNSIGNED NOT NULL AUTO_INCREMENT",
			"name VARCHAR(255) NOT NULL",
			"user VARCHAR(255) NOT NULL",
			"FOREIGN KEY (user) REFERENCES users(id) ON DELETE CASCADE",
			"UNIQUE KEY (name, user)",
			"PRIMARY KEY (id)",
		},
	},
	{
		name: "recipe_ingredients",
		columns: []string{
			"recipe INT UNSIGNED REFERENCES recipes(id) ON DELETE CASCADE",
			"product INT UNSIGNED REFERENCES products(id) ON DELETE CASCADE",
			"amount FLOAT NOT NULL",
			"PRIMARY KEY (recipe, product)",
		},
	},
}

func (s *SQL) Recipes(user string) ([]recipe.Recipe, error) {
	tx, err := s.db.BeginTx(s.ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	recs, err := s.queryRecipes(tx, user)
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

func (s *SQL) LookupRecipe(asUser string, id recipe.ID) (recipe.Recipe, error) {
	tx, err := s.db.BeginTx(s.ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return recipe.Recipe{}, fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	rec := recipe.Recipe{
		ID: id,
	}

	row := tx.QueryRowContext(s.ctx, "SELECT user, name FROM recipes WHERE id = ? AND user = ?", id, asUser)
	if err := row.Scan(&rec.User, &rec.Name); errorIs(err, errKeyNotFound) {
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

func (s *SQL) queryRecipes(tx *sql.Tx, user string) ([]recipe.Recipe, error) {
	query := `SELECT id, user, name FROM recipes WHERE user = ?`

	s.log.Tracef(query)
	r, err := tx.QueryContext(s.ctx, query, user)
	if err != nil {
		return nil, fmt.Errorf("could not query: %v", err)
	}
	defer r.Close()

	var recs []recipe.Recipe
	for r.Next() {
		var rec recipe.Recipe
		if err := r.Scan(&rec.ID, &rec.User, &rec.Name); err != nil {
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
		product, amount
	FROM
		recipe_ingredients
	WHERE
		recipe = ?
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
	if r.User == "" {
		return 0, errors.New("user cannot be empty")
	} else if r.Name == "" {
		return 0, errors.New("name cannot be empty")
	}

	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	if r.ID == 0 {
		if r.ID, err = s.insertNewRecipe(tx, r); err != nil {
			return 0, fmt.Errorf("could not insert recipe: %v", err)
		}
	} else {
		if err := s.insertRecipeWithID(tx, r); err != nil {
			return 0, fmt.Errorf("could not insert recipe: %v", err)
		}
	}

	if err := s.deleteRecipeIngredients(tx, r.ID); err != nil {
		return 0, fmt.Errorf("could not delete old ingredients: %v", err)
	}

	err = bulkInsert(s, tx,
		"recipe_ingredients(recipe, product, amount)",
		r.Ingredients, func(i recipe.Ingredient) []any {
			return []any{uint(r.ID), uint(i.ProductID), i.Amount}
		})
	if err != nil {
		return 0, fmt.Errorf("could not insert ingredients: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("could not commit transaction: %v", err)
	}

	return r.ID, nil
}

func (s *SQL) insertRecipeWithID(tx *sql.Tx, r recipe.Recipe) error {
	query := `
	UPDATE
		recipes
	SET
		name = ?
	WHERE
		id = ?
		AND user = ?
	`
	s.log.Tracef(query)

	result, err := tx.ExecContext(s.ctx, query, r.Name, r.ID, r.User)
	if err != nil {
		return fmt.Errorf("could not update recipe: %v", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected: %v", err)
	}

	s.log.Tracef("Rows affected: %d", affected)

	switch affected {
	case 0:
	case 1:
		return nil
	default:
		return fmt.Errorf("unexpected number of rows affected: %d", affected)
	}

	if !s.allowInsertNewID {
		return nil
	}

	// No rows affected, try to insert a new one if allowed
	query = `INSERT INTO recipes (id, user, name) VALUES (?, ?, ?)`
	s.log.Tracef(query)

	_, err = tx.ExecContext(s.ctx, query, r.ID, r.User, r.Name)
	if err != nil && !errorIs(err, errKeyExists) {
		return fmt.Errorf("could not insert recipe: %v", err)
	}

	return nil
}

func (s *SQL) insertNewRecipe(tx *sql.Tx, r recipe.Recipe) (recipe.ID, error) {
	queryRecipe := `INSERT INTO recipes (user, name) VALUES (?, ?)`
	s.log.Tracef(queryRecipe)

	result, err := tx.ExecContext(s.ctx, queryRecipe, r.User, r.Name)
	if err != nil {
		return 0, fmt.Errorf("could not insert recipe: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("could not get last insert ID: %v", err)
	}

	s.log.Tracef("Last insert ID: %d", id)

	return recipe.ID(id), nil
}

func (s *SQL) DeleteRecipe(asUser string, id recipe.ID) error {
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() //nolint:errcheck // The error is irrelevant

	if err := s.deleteRecipe(tx, asUser, id); err != nil {
		return fmt.Errorf("could not delete recipe: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}

func (s *SQL) deleteRecipe(tx *sql.Tx, asUser string, id recipe.ID) error {
	query := `DELETE FROM recipes WHERE id = ?`
	argv := []any{id}

	if asUser != "" {
		query += ` AND user = ?`
		argv = append(argv, asUser)
	}

	s.log.Tracef(query)

	_, err := tx.ExecContext(s.ctx, query, argv...)
	if err != nil {
		return fmt.Errorf("could not delete product: %v", err)
	}

	return nil
}

func (s *SQL) deleteRecipeIngredients(tx *sql.Tx, id recipe.ID) error {
	query := `DELETE FROM recipe_ingredients WHERE recipe = ?`
	s.log.Tracef(query)

	_, err := tx.ExecContext(s.ctx, query, id)
	if err != nil {
		return fmt.Errorf("could not delete ingredients: %v", err)
	}

	return nil
}
