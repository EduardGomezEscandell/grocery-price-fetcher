package formatter

import (
	"io"
	"slices"
	"strings"

	"golang.org/x/exp/maps"
)

type JSON struct {
	isFirstRow bool
}

func (f *JSON) PrintHead(w io.Writer, _ ...string) error {
	_, err := fmter.Fprintf(w, "[")
	f.isFirstRow = true
	return err
}

func (f *JSON) PrintRow(w io.Writer, data map[string]interface{}) error {
	if f.isFirstRow {
		if _, err := fmter.Fprintf(w, "\n\t{"); err != nil {
			return err
		}
		f.isFirstRow = false
	} else {
		// Print comma only if it's not the first row
		if _, err := fmter.Fprintf(w, ",\n\t{"); err != nil {
			return err
		}
	}

	// Print the key-value pairs
	first := true
	keys := maps.Keys(data)
	slices.Sort(keys)
	for _, k := range keys {
		if !first {
			// Print comma if it's not the first pair
			if _, err := fmter.Fprintf(w, ", "); err != nil {
				return err
			}
		}
		first = false

		out := format(data[k]).IfFloat32("%.2f").IfEuro("%.2f").IfString("%q").OrElse("%q")

		_, err := fmter.Fprintf(w, "%q: %s", strings.ToLower(k), out)
		if err != nil {
			return err
		}
	}

	if _, err := fmter.Fprint(w, "}"); err != nil {
		return err
	}

	return nil
}

func (f *JSON) PrintTail(w io.Writer) error {
	_, err := fmter.Fprintln(w, "\n]")
	return err
}
