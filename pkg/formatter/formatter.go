package formatter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Formatter interface {
	PrintHead(w io.Writer, columns ...string) error
	PrintRow(w io.Writer, row map[string]any) error
	PrintTail(w io.Writer) error
}

func New(format string) (Formatter, error) {
	switch format {
	case "json":
		return &JSON{}, nil
	case "tsv":
		return &CharSV{Separator: "\t"}, nil
	case "csv":
		return &CharSV{Separator: ","}, nil
	case "table":
		return &Table{}, nil
	default:
		return nil, fmt.Errorf("unsupported format %q", format)
	}
}

var lp *localizedPrinter

func localePrinter() *localizedPrinter {
	if lp != nil {
		return lp
	}

	locale := strings.Split(os.Getenv("LC_NUMERIC"), ".")[0]

	tag, err := language.Parse(locale)
	if err != nil {
		logrus.Warningf("Locale: defaulting to english because locale %q was not found: %v", locale, err)
		tag = language.English
	} else {
		logrus.Debugf("Using locale %s", tag)
	}

	lp = &localizedPrinter{message.NewPrinter(tag)}
	return lp
}

type printer interface {
	Sprintf(format string, a ...any) string
	Sprint(a ...any) string
}

type valueFmt struct {
	value any
	fmt   printer
}

type stdPrinter struct{}

func (stdPrinter) Sprintf(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

func (stdPrinter) Sprint(a ...any) string {
	return fmt.Sprint(a...)
}

type localizedPrinter struct {
	*message.Printer
}

func (p localizedPrinter) Sprintf(format string, a ...any) string {
	return p.Printer.Sprintf(format, a...)
}

func (p localizedPrinter) Sprint(a ...any) string {
	return p.Printer.Sprint(a...)
}

type finalValue struct {
	string
}

func format(v any, localized bool) *valueFmt {
	if localized {
		return &valueFmt{v, localePrinter()}
	}
	return &valueFmt{v, stdPrinter{}}
}

func (f *valueFmt) Localize() *valueFmt {
	return f
}

func (f *valueFmt) IfString(fmt string) *valueFmt {
	fmtfmt[string](f, fmt)
	return f
}

func (f *valueFmt) IfFloat32(fmt string) *valueFmt {
	fmtfmt[float32](f, fmt)
	return f
}

func (f *valueFmt) IfEuro(fmt string) *valueFmt {
	if e, ok := f.value.(euro); ok {
		f.value = finalValue{f.fmt.Sprintf(fmt, e.float32)}
	}

	return f
}

func (f *valueFmt) String() string {
	if v, ok := f.value.(finalValue); ok {
		return v.string
	}

	return f.fmt.Sprint(f.value)
}

func (f *valueFmt) OrElse(fmt string) *valueFmt {
	if _, ok := f.value.(finalValue); !ok {
		f.value = finalValue{f.fmt.Sprintf(fmt, f.value)}
	}

	return f
}

func fmtfmt[T any](f *valueFmt, fmt string) {
	if e, ok := f.value.(T); ok {
		f.value = finalValue{f.fmt.Sprintf(fmt, e)}
	}
}

func (f *valueFmt) Fprint(w io.Writer) error {
	_, err := fmt.Fprint(w, f.String())
	return err
}

func (f *valueFmt) Fprintf(w io.Writer, format string) error {
	_, err := fmt.Fprintf(w, format, f.String())
	return err
}
