package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/cmd/needs/formatter"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/product"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/provider"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/recipe"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/providers/bonpreu"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/providers/mercadona"
	log "github.com/sirupsen/logrus"
)

func main() {
	s := parseInput()
	if s.verbose {
		log.SetLevel(log.DebugLevel)
	}

	outFmt, err := formatter.Get(s.format)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	var in *os.File
	if s.inputPath == "" {
		in = os.Stdin
		log.Info("Reading from STDIN. Write your products and press Ctrl+D to finish.")
	} else {
		in, err = os.Open(s.inputPath)
		if err != nil {
			log.Fatalf("could not open file: %v", err)
		}
		defer in.Close()
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

	products, err := run(in)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if err := formatOutput(out, products, outFmt); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

type settings struct {
	verbose    bool
	format     string
	inputPath  string
	outputPath string
}

func parseInput() settings {
	var sett settings

	flag.BoolVar(&sett.verbose, "v", false, "verbose mode")
	flag.StringVar(&sett.format, "fmt", "table", "output format (json, csv, tsv, ini)")
	flag.StringVar(&sett.inputPath, "i", "", "input file path (default: STDIN)")
	flag.StringVar(&sett.outputPath, "o", "", "output file path (default: STDOUT)")

	flag.Parse()
	return sett
}

type ProductCount struct {
	product product.Product
	count   float32
}

type Database struct {
	Products []product.Product
	Recipes  []recipe.Recipe
}

func run(r io.Reader) ([]ProductCount, error) {
	var db Database

	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &db); err != nil {
		return nil, err
	}

	counts := make(map[string]float32)
	for _, r := range db.Products {
		counts[r.Name] = 0
	}

	log.Debug("Database loaded successfully")
	log.Debugf("Products: %d", len(db.Products))
	log.Debugf("Recipes: %d", len(db.Recipes))

	// Calculate the amount of each product needed
	for _, r := range db.Recipes {
		for _, i := range r.Ingredients {
			_, ok := counts[i.Name]
			if !ok {
				log.Warningf("Recipe %q contains product %q which is not registered", r.Name, i.Name)
				continue
			}
			counts[i.Name] += i.Amount
		}
	}

	// Create the output
	products := make([]ProductCount, 0, len(counts))
	for _, p := range db.Products {
		products = append(products, ProductCount{
			product: p,
			count:   counts[p.Name],
		})
	}

	return products, nil
}

func formatOutput(w io.Writer, products []ProductCount, f formatter.Formatter) error {
	if err := f.PrintHead(w); err != nil {
		return fmt.Errorf("could not write header to output: %w", err)
	}

	for _, p := range products {
		if err := f.Println(w, p.product.Name, p.count); err != nil {
			return fmt.Errorf("could not write results to output: %w", err)
		}
	}

	if err := f.PrintTail(w); err != nil {
		return fmt.Errorf("could not write footer to output: %w", err)
	}

	return nil
}
