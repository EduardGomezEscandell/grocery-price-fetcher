package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/database"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/formatter"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/menu"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/provider"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/providers/bonpreu"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/providers/mercadona"
	log "github.com/sirupsen/logrus"
)

func main() {
	s := parseInput()
	if s.verbose {
		log.SetLevel(log.DebugLevel)
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

	provider.Register("Bonpreu", bonpreu.New)
	provider.Register("Mercadona", mercadona.New)

	db, err := unmarshalFile[database.DB](s.dbPath)
	if err != nil {
		log.Fatalf("Error: could not parse database: %v", err)
	}

	if err := db.Validate(); err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Debug("Database loaded successfully")
	log.Debugf("Products: %d", len(db.Products))
	log.Debugf("Recipes: %d", len(db.Recipes))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db.UpdatePrices(ctx)

	pantry, err := func() ([]menu.ProductData, error) {
		if s.pantryPath == "" {
			return nil, nil
		}
		return unmarshalFile[[]menu.ProductData](s.pantryPath)
	}()
	if err != nil {
		log.Fatalf("Error: could not parse pantry: %v", err)
	}

	menu, err := unmarshalFile[menu.Menu](s.menuPath)
	if err != nil {
		log.Fatalf("Error: could not parse menu: %v", err)
	}

	log.Debug("Menu loaded successfully")
	log.Debugf("Days: %d", len(menu.Days))

	products, err := menu.Compute(&db, pantry)
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

	menuPath   string
	pantryPath string
	dbPath     string
	outputPath string
}

func parseInput() settings {
	var sett settings

	flag.StringVar(&sett.format, "fmt", "table", "output format (json, csv, tsv, ini)")
	flag.StringVar(&sett.dbPath, "db", "", "database file path")
	flag.StringVar(&sett.pantryPath, "p", "", "pantry file path")
	flag.StringVar(&sett.menuPath, "i", "", "input file path (default: STDIN)")
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
		if skipEmpty && p.Amount == 0 {
			continue
		}

		if err := f.PrintRow(w, map[string]any{
			"Product":   p.Name,
			"Amount":    p.Amount,
			"Unit Cost": formatter.Euro(p.UnitCost),
			"Cost":      formatter.Euro(p.UnitCost * p.Amount),
		}); err != nil {
			return fmt.Errorf("could not write results to output: %w", err)
		}

		total += p.UnitCost * p.Amount
	}

	if err := f.PrintRow(w, map[string]any{
		"Product":   "Total",
		"Amount":    "",
		"Unit Cost": "",
		"Cost":      formatter.Euro(total),
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
