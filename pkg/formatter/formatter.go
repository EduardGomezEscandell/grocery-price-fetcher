package formatter

import (
	"fmt"
	"io"

	"example.com/compra/pkg/product"
)

var formats = map[string]Formatter{
	"":     {"Product             Price\n", "%-20s %5.2f €\n", ""},
	"json": {"[\n", `{"name": "%s", "price": %5.2f}`, "]\n"},
	"csv":  {"Product,Price\n", "%s,%5.2f\n", ""},
	"tsv":  {"Product\tPrice\n", "%s\t%5.2f\n", ""},
	"ini":  {"[products]\n", "%s = %5.2f €\n", ""},
}

func Get(name string) (Formatter, error) {
	f, ok := formats[name]
	if !ok {
		return Formatter{}, fmt.Errorf("unknown format %q", name)
	}

	return f, nil
}

type Formatter struct {
	head, body, tail string
}

func (f Formatter) PrintHead(w io.Writer) error {
	_, err := fmt.Fprint(w, f.head)
	return err
}

func (f Formatter) Println(w io.Writer, p *product.Product) error {
	_, err := fmt.Fprintf(w, f.body, p.Name, p.Price)
	return err
}

func (f Formatter) PrintTail(w io.Writer) error {
	_, err := fmt.Fprint(w, f.tail)
	return err
}
