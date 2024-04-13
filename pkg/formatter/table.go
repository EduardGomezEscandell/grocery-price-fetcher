package formatter

import (
	"io"
	"slices"

	"github.com/sirupsen/logrus"
)

type Table struct {
	fields []string
}

func (f *Table) PrintHead(w io.Writer, fields ...string) error {
	f.fields = fields

	for _, k := range fields {
		_, err := localePrinter().Fprintf(w, "%-20s", k)
		if err != nil {
			return err
		}
	}

	if _, err := localePrinter().Fprintln(w); err != nil {
		return err
	}

	return nil
}

func (f *Table) PrintRow(w io.Writer, data map[string]interface{}) error {
	count := 0
	for _, k := range f.fields {
		v, ok := data[k]
		if !ok {
			localePrinter().Fprintf(w, "%-20s", "")
			continue // Default value: empty
		}

		err := format(v, true).IfString("%-20s").IfFloat32("%20.2f").IfEuro("%10.2f €").Fprint(w)
		if err != nil {
			return err
		}
		count++
	}

	if _, err := localePrinter().Fprintln(w); err != nil {
		return err
	}

	if count != len(data) {
		f.warnUnknownField(data)
	}

	return nil
}

func (f *Table) PrintTail(_ io.Writer) error {
	return nil
}

func (f *Table) warnUnknownField(data map[string]interface{}) {
	if !logrus.IsLevelEnabled(logrus.WarnLevel) {
		return
	}

	for k := range data {
		if slices.Contains(f.fields, k) {
			continue
		}
		logrus.Warningf("Formatter: field %q used in was not registered in PrintHead", k)
	}
}
