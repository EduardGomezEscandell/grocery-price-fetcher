package formatter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/product"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var formats = map[string]Formatter{
	"table": {
		head: "Product               Price\n",
		body: "%-20s %5.2f €",
		sep:  "\n",
		tail: "\n",
	}, "json": {
		head: "[",
		body: `{"name":%q,"price":%.2f}`,
		sep:  ",",
		tail: "]\n",
	}, "csv": {
		head: "Product,Price\n",
		body: "%s,%.2f",
		sep:  "\n",
		tail: "\n",
	}, "tsv": {
		head: "Product\tPrice\n",
		body: "%s\t%.2f",
		sep:  "\n",
		tail: "\n",
	}, "ini": {
		head: "[products]\n",
		body: "%s = %.2f €",
		sep:  "\n",
		tail: "\n",
	},
}

func Get(name string) (Formatter, error) {
	f, ok := formats[name]
	if !ok {
		return Formatter{}, fmt.Errorf("unknown format %q", name)
	}

	locale := strings.Split(os.Getenv("LC_NUMERIC"), ".")[0]

	tag, err := language.Parse(locale)
	if err != nil {
		log.Warningf("Locale: defaulting to english because locale %q was not found: %v", locale, err)
		tag = language.English
	} else {
		log.Debugf("Using locale %s", tag)
	}

	f.printer = message.NewPrinter(tag)

	return f, nil
}

type Formatter struct {
	head, body, sep, tail string
	line                  int
	printer               *message.Printer
}

func (f Formatter) PrintHead(w io.Writer) error {
	_, err := f.printer.Fprint(w, f.head)
	return err
}

func (f *Formatter) Println(w io.Writer, p *product.Product) error {
	if f.line != 0 {
		_, err := f.printer.Fprint(w, f.sep)
		if err != nil {
			return err
		}
	}

	if _, err := f.printer.Fprintf(w, f.body, p.Name, p.Price); err != nil {
		return err
	}

	f.line++
	return nil
}

func (f Formatter) PrintTail(w io.Writer) error {
	_, err := f.printer.Fprint(w, f.tail)
	return err
}
