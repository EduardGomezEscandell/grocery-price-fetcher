package formatter

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"golang.org/x/exp/maps"
)

type JSON struct {
	isFirstRow bool
}

func (f *JSON) PrintHead(w io.Writer, _ ...string) error {
	_, err := fmt.Fprintf(w, "[")
	f.isFirstRow = true
	return err
}

func (f *JSON) PrintRow(w io.Writer, data map[string]interface{}) error {
	if f.isFirstRow {
		if _, err := fmt.Fprintf(w, "\n\t{"); err != nil {
			return err
		}
		f.isFirstRow = false
	} else {
		// Print comma only if it's not the first row
		if _, err := fmt.Fprintf(w, ",\n\t{"); err != nil {
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
			if _, err := fmt.Fprintf(w, ", "); err != nil {
				return err
			}
		}
		first = false

		key := strings.ReplaceAll(strings.ToLower(k), " ", "_")
		value := format(data[k], false).IfFloat32("%.2f").IfEuro("%.2f").IfString("%q").OrElse("%q").String()
		_, err := fmt.Fprintf(w, "%q: %s", key, value)
		if err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(w, "}"); err != nil {
		return err
	}

	return nil
}

func (f *JSON) PrintTail(w io.Writer) error {
	_, err := fmt.Fprintln(w, "\n]")
	return err
}
