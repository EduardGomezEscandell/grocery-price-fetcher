package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/blank"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/bonpreu"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/providers/mercadona"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Settings struct {
	Verbosity int
	Overwrite bool
	Input     database.Settings
	Output    database.Settings
}

const (
	ExitSuccess = iota
	ExitBadInput
	ExitError
)

func main() {
	var settings string
	var verbosity int

	flag.StringVar(&settings, "settings", "", "Path to the settings file")
	flag.IntVar(&verbosity, "v", 0, "Enable verbose logging (defaults to the settings file value, otherwise 4)")
	flag.Parse()

	if settings == "" {
		fmt.Fprintln(os.Stderr, "No settings file provided")
		os.Exit(ExitBadInput)
	}

	out, err := os.ReadFile(settings)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read settings file: %v\n", err)
		os.Exit(ExitBadInput)
	}

	s := Settings{
		Input:     database.Settings{}.Defaults(),
		Output:    database.Settings{}.Defaults(),
		Verbosity: int(logrus.InfoLevel),
	}

	if err := yaml.Unmarshal(out, &s); err != nil {
		fmt.Fprintf(os.Stderr, "Could not unmarshal settings: %v\n", err)
		os.Exit(ExitBadInput)
	}

	// Override verbosity with command line flag
	if verbosity != 0 {
		s.Verbosity = verbosity
	}

	os.Exit(run(s))
}

func run(s Settings) int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := logger.New()
	log.SetLevel(s.Verbosity)

	if s.Verbosity > int(logrus.InfoLevel) {
		if out, err := yaml.Marshal(s); err == nil {
			log.Debugf("Running with options:\n%s", string(out))
		}
	}

	providers.Register(blank.Provider{})
	providers.Register(bonpreu.New(log))
	providers.Register(mercadona.New(log))

	input, err := database.New(ctx, log, s.Input)
	if err != nil {
		fmt.Printf("Could not create input database: %v\n", err)
		return ExitError
	}
	defer input.Close()

	output, err := database.New(ctx, log, s.Output)
	if err != nil {
		fmt.Printf("Could not create output database: %v\n", err)
		return ExitError
	}
	defer output.Close()

	if !s.Overwrite {
		if ok, err := verifyEmpty(log, output); err != nil {
			fmt.Fprintf(os.Stderr, "Failed output database empty check: %v\n", err)
			return ExitError
		} else if !ok {
			return ExitSuccess
		}
	}

	for _, copyFunc := range []func(logger.Logger, database.DB, database.DB) error{
		copyProducts,
		copyRecipes,
		copyMenus,
		copyPantries,
		copyShoppingLists,
	} {
		if err := copyFunc(log, input, output); err != nil {
			fmt.Fprintf(os.Stderr, "Could not copy data: %v\n", err)
			return ExitError
		}
	}

	log.Debug("Data copied successfully")
	return ExitSuccess
}

func copyProducts(log logger.Logger, input, output database.DB) error {
	products, err := input.Products()
	if err != nil {
		return fmt.Errorf("could not get products: %w", err)
	}

	for _, p := range products {
		log.Tracef("Copying %T %s", p, p.Name)

		if err := output.SetProduct(p); err != nil {
			return fmt.Errorf("could not set product: %w", err)
		}
	}

	return nil
}

func copyRecipes(log logger.Logger, input, output database.DB) error {
	recipes, err := input.Recipes()
	if err != nil {
		return fmt.Errorf("could not get recipes: %w", err)
	}

	for _, r := range recipes {
		log.Tracef("Copying %T %s", r, r.Name)

		if err := output.SetRecipe(r); err != nil {
			return fmt.Errorf("could not set recipe: %w", err)
		}
	}

	return nil
}

func copyMenus(log logger.Logger, input, output database.DB) error {
	menus, err := input.Menus()
	if err != nil {
		return fmt.Errorf("could not get menus: %w", err)
	}

	for _, m := range menus {
		log.Tracef("Copying %T %s", m, m.Name)

		if err := output.SetMenu(m); err != nil {
			return fmt.Errorf("could not set menu: %w", err)
		}
	}

	return nil
}

func copyPantries(log logger.Logger, input, output database.DB) error {
	pantries, err := input.Pantries()
	if err != nil {
		return fmt.Errorf("could not get pantries: %w", err)
	}

	for _, p := range pantries {
		log.Tracef("Copying %T %s", p, p.Name)

		if err := output.SetPantry(p); err != nil {
			return fmt.Errorf("could not set pantry: %w", err)
		}
	}

	return nil
}

func copyShoppingLists(log logger.Logger, input, output database.DB) error {
	shoppingLists, err := input.ShoppingLists()
	if err != nil {
		return fmt.Errorf("could not get shopping lists: %w", err)
	}

	for _, s := range shoppingLists {
		log.Tracef("Copying %T %s", s, s.Name)

		if err := output.SetShoppingList(s); err != nil {
			return fmt.Errorf("could not set shopping list: %w", err)
		}
	}

	return nil
}

func verifyEmpty(log logger.Logger, db database.DB) (bool, error) {
	if products, err := db.Products(); err != nil {
		return false, fmt.Errorf("could not read output DB products: %v", err)
	} else if len(products) != 0 {
		log.Info("Output DB products not empty, exiting")
		return false, nil
	}

	if recipes, err := db.Recipes(); err != nil {
		return false, fmt.Errorf("could not read output DB recipes: %v", err)
	} else if len(recipes) != 0 {
		log.Info("Output DB recipes not empty, exiting")
		return false, nil
	}

	if menus, err := db.Menus(); err != nil {
		return false, fmt.Errorf("could not read output DB menus: %v", err)
	} else if len(menus) != 0 {
		log.Info("Output DB menus not empty, exiting")
		return false, nil
	}

	if pantries, err := db.Pantries(); err != nil {
		return false, fmt.Errorf("could not read output DB pantries: %v", err)
	} else if len(pantries) != 0 {
		log.Info("Output DB pantries not empty, exiting")
		return false, nil
	}

	if shoppingLists, err := db.ShoppingLists(); err != nil {
		return false, fmt.Errorf("could not read output DB shopping lists: %v", err)
	} else if len(shoppingLists) != 0 {
		log.Info("Output DB shopping lists not empty, exiting")
		return false, nil
	}

	return true, nil
}
