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

	log "github.com/sirupsen/logrus"

	"example.com/compra/pkg/formatter"
	"example.com/compra/pkg/product"
	"example.com/compra/providers/bonpreu"
	"example.com/compra/providers/mercadona"
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

	in, err := os.Open(s.inputPath)
	if err != nil {
		log.Fatalf("could not open file: %v", err)
	}
	defer in.Close()

	out, err := os.Open(s.outputPath)
	if err != nil {
		log.Fatalf("could not open file: %v", err)
	}
	defer out.Close()

	var registry product.Registry
	registry.Register("Bonpreu", bonpreu.Get)
	registry.Register("Mercadona", mercadona.Get)

	products, err := run(in, registry)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if err := printOutput(out, products, outFmt); err != nil {
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
	flag.StringVar(&sett.format, "fmt", "", "output format (json, csv, tsv, ini)")
	flag.StringVar(&sett.inputPath, "i", os.Stdin.Name(), "input file path")
	flag.StringVar(&sett.outputPath, "o", os.Stdout.Name(), "output file path")

	flag.Parse()
	return sett
}

func run(r io.Reader, reg product.Registry) ([]*product.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

			log.Debug("Parsing row: ", row)

			if err := p.Parse(reg, row); err != nil {
				log.Warningf("could not parse row %q: %v", row, err)
				return
			}

			err := p.Get(ctx, &reg)
			if err != nil {
				log.Warningf("Get: %v", err)
				return
			}
		}()
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("could not scan file: %v", err)
	}

	wg.Wait()

	return products, nil
}

func printOutput(w io.Writer, products []*product.Product, f formatter.Formatter) error {
	if err := f.PrintHead(w); err != nil {
		return fmt.Errorf("could not write to output: %v", err)
	}

	for _, p := range products {
		if err := f.Println(w, p); err != nil {
			return fmt.Errorf("could not write to output: %v", err)
		}
	}

	if err := f.PrintTail(w); err != nil {
		return fmt.Errorf("could not write to output: %v", err)
	}

	return nil
}
