package formatter

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/sirupsen/logrus"
)

type CharSV struct {
	Separator string
	fields    []string
}

func (f *CharSV) PrintHead(w io.Writer, fields ...string) error {
	f.fields = fields

	for i, k := range fields {
		if i > 0 {
			if _, err := localePrinter().Fprint(w, f.Separator); err != nil {
				return err
			}
		}

		if strings.Contains(k, f.Separator) {
			k = fmt.Sprintf("%q", k)
		}

		_, err := localePrinter().Fprint(w, k)
		if err != nil {
			return err
		}
	}

	if _, err := localePrinter().Fprintln(w); err != nil {
		return err
	}

	return nil
}

func (f *CharSV) PrintRow(w io.Writer, data map[string]interface{}) error {
	count := 0
	for i, k := range f.fields {
		if i > 0 {
			if _, err := localePrinter().Fprintf(w, f.Separator); err != nil {
				return err
			}
		}

		v, ok := data[k]
		if !ok {
			continue // Default value: empty
		}

		out := format(v, true).IfFloat32("%.2f").IfEuro("%.2f").IfString("%s").OrElse("%v").String()
		if strings.Contains(out, f.Separator) {
			out = fmt.Sprintf("%q", out)
		}

		_, err := fmt.Fprint(w, out)
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

func (f *CharSV) PrintTail(_ io.Writer) error {
	return nil
}

func (f *CharSV) warnUnknownField(data map[string]interface{}) {
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
