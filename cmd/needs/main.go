package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/formatter"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/logger"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/providers"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/providers/bonpreu"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/providers/mercadona"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/services/menu"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/services/pricing"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func main() {
	s := parseInput()
	log := logger.New()
	if s.verbose {
		log.SetLevel(int(logrus.DebugLevel))
	}

	outFmt, err := formatter.New(s.format)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	var out *os.File
	if s.outputPath == "" {
		out = os.Stdout
	} else {
		out, err = os.Create(s.outputPath)
		if err != nil {
			log.Fatalf("could not open file: %v", err)
		}
		defer out.Close()
	}

	providers.Register("Bonpreu", bonpreu.New)
	providers.Register("Mercadona", mercadona.New)

	db, err := loadDB(context.Background(), log, s.dbPath)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer db.Close()

	log.Debug("Database loaded successfully")
	log.Debugf("Products: %d", len(db.Products()))
	log.Debugf("Recipes: %d", len(db.Recipes()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pricing.OneShot(ctx, log, db)

	pantry, err := func() ([]menu.ProductData, error) {
		if s.pantryPath == "" {
			return nil, nil
		}
		return unmarshalFile[[]menu.ProductData](s.pantryPath)
	}()
	if err != nil {
		log.Fatalf("Error: could not parse pantry: %v", err)
	}

	m := db.Menus()
	switch len(m) {
	case 0:
		log.Fatalf("Error: no menus found in database")
	case 1:
		break
	default:
		log.Fatalf("Error: more than one menu found in database")
	}

	products, err := menu.OneShot(log, db, m[0], pantry)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if err := formatOutput(out, products, outFmt, s.skipEmpty); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

type settings struct {
	verbose   bool
	format    string
	skipEmpty bool

	pantryPath string
	dbPath     string
	outputPath string
}

func parseInput() settings {
	var sett settings

	flag.StringVar(&sett.format, "fmt", "table", "output format (json, csv, tsv, ini)")
	flag.StringVar(&sett.dbPath, "db", "", "database manifest path")
	flag.StringVar(&sett.pantryPath, "p", "", "pantry file path")
	flag.StringVar(&sett.outputPath, "o", "", "output file path (default: STDOUT)")
	flag.BoolVar(&sett.verbose, "v", false, "verbose mode")
	flag.BoolVar(&sett.skipEmpty, "skip-empty", false, "skip empty products")

	flag.Parse()
	return sett
}

func formatOutput(w io.Writer, products []menu.ProductData, f formatter.Formatter, skipEmpty bool) error {
	if err := f.PrintHead(w, "Product", "Amount", "Cost"); err != nil {
		return fmt.Errorf("could not write header to output: %w", err)
	}

	var total float32
	for _, p := range products {
		if skipEmpty && p.Need <= p.Have {
			fmt.Printf("%+v\n", p)
			continue
		}

		need := max(p.Need-p.Have, 0)
		if err := f.PrintRow(w, map[string]any{
			"Product": p.Name,
			"Amount":  need,
			"Cost":    formatter.Euro(p.UnitCost * need),
		}); err != nil {
			return fmt.Errorf("could not write results to output: %w", err)
		}

		total += p.UnitCost * need
	}

	if err := f.PrintRow(w, map[string]any{
		"Product": "Total",
		"Amount":  "",
		"Cost":    formatter.Euro(total),
	}); err != nil {
		return fmt.Errorf("could not write total to output: %w", err)
	}

	if err := f.PrintTail(w); err != nil {
		return fmt.Errorf("could not write footer to output: %w", err)
	}

	return nil
}

func unmarshalFile[T any](path string) (T, error) {
	var t T

	b, err := os.ReadFile(path)
	if err != nil {
		return t, err
	}

	if err := json.Unmarshal(b, &t); err != nil {
		return t, err
	}

	return t, nil
}

func loadDB(ctx context.Context, log logger.Logger, path string) (database.DB, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}

	out, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read database manifest: %v", err)
	}

	var conf database.Settings
	if err := yaml.Unmarshal(out, &conf); err != nil {
		return nil, fmt.Errorf("could not parse database manifest: %v", err)
	}

	db, err := database.New(ctx, log, conf)
	if err != nil {
		return nil, fmt.Errorf("could not load database: %v", err)
	}

	return db, nil
}
