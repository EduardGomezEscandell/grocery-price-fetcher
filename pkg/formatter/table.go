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
		_, err := fmter.Fprintf(w, "%-20s", k)
		if err != nil {
			return err
		}
	}

	if _, err := fmter.Fprintln(w); err != nil {
		return err
	}

	return nil
}

func (f *Table) PrintRow(w io.Writer, data map[string]interface{}) error {
	count := 0
	for _, k := range f.fields {
		v, ok := data[k]
		if !ok {
			fmter.Fprintf(w, "%-20s", "")
			continue // Default value: empty
		}

		out := format(v).IfFloat32("%.2f").IfEuro("%.2f €").String()
		_, err := fmter.Fprintf(w, "%-20s", out)
		if err != nil {
			return err
		}
		count++
	}

	if _, err := fmter.Fprintln(w); err != nil {
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
