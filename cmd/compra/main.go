package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"example.com/compra/pkg/formatter"
	"example.com/compra/pkg/product"
	"example.com/compra/providers/bonpreu"
	"example.com/compra/providers/mercadona"
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

	var registry product.Registry
	registry.Register("Bonpreu", bonpreu.Get)
	registry.Register("Mercadona", mercadona.Get)

	products, err := run(in, registry)
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

func run(r io.Reader, reg product.Registry) ([]*product.Product, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scanner := bufio.NewScanner(r)
	var products []*product.Product

	var wg sync.WaitGroup
	for scanner.Scan() {
		row := scanner.Text()
		var p product.Product
		products = append(products, &p)

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := p.Parse(row); err != nil {
				log.Warningf("could not parse row %q: %v", row, err)
				return
			}

			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			err := p.Get(ctx, &reg)
			if err != nil {
				log.Warningf("Get: %v", err)
				return
			}
		}()
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("could not scan file: %w", err)
	}

	wg.Wait()
	log.Debug("All products have been processed")

	return products, nil
}

func formatOutput(w io.Writer, products []*product.Product, f formatter.Formatter) error {
	if err := f.PrintHead(w); err != nil {
		return fmt.Errorf("could not write header to output: %w", err)
	}

	for _, p := range products {
		if err := f.Println(w, p); err != nil {
			return fmt.Errorf("could not write results to output: %w", err)
		}
	}

	if err := f.PrintTail(w); err != nil {
		return fmt.Errorf("could not write footer to output: %w", err)
	}

	return nil
}
