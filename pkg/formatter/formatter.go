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

var fmter = func() *message.Printer {
	locale := strings.Split(os.Getenv("LC_NUMERIC"), ".")[0]

	tag, err := language.Parse(locale)
	if err != nil {
		logrus.Warningf("Locale: defaulting to english because locale %q was not found: %v", locale, err)
		tag = language.English
	} else {
		logrus.Debugf("Using locale %s", tag)
	}

	return message.NewPrinter(tag)
}()

type Formatter interface {
	PrintHead(w io.Writer, columns ...string) error
	PrintRow(w io.Writer, row map[string]interface{}) error
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

type finalValue struct {
	string
}

type valueFmt struct {
	value any
}

func format(v any) *valueFmt {
	return &valueFmt{v}
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
		f.value = finalValue{fmter.Sprintf(fmt, e.float32)}
	}

	return f
}

func (f *valueFmt) String() string {
	if v, ok := f.value.(finalValue); ok {
		return v.string
	}

	return fmter.Sprint(f.value)
}

func (f *valueFmt) OrElse(fmt string) string {
	if _, ok := f.value.(finalValue); !ok {
		f.value = finalValue{fmter.Sprintf(fmt, f.value)}
	}

	return f.String()
}

func fmtfmt[T any](f *valueFmt, fmt string) {
	if e, ok := f.value.(T); ok {
		f.value = finalValue{fmter.Sprintf(fmt, e)}
	}
}
