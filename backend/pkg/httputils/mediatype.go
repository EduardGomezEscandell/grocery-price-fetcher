package httputils

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	MediaTypeJSON = NewMediaType("application", "json")
	MediaTypeText = NewMediaType("text", "plain")
	MediaTypeHTML = NewMediaType("text", "html")
)

// MediaType represents a media type as defined in RFC 6838,
// with a lenient parser that ignores parameters.
type MediaType struct {
	Type    string
	Subtype string
}

func (m MediaType) String() string {
	return fmt.Sprintf("%s/%s", m.Type, m.Subtype)
}

func NewMediaType(t, s string) MediaType {
	return MediaType{Type: t, Subtype: s}
}

func ParseMediaTypes(s string) ([]MediaType, error) {
	var out []MediaType
	var errs error

	var empty bool
	for _, mt := range strings.Split(s, ",") {
		if len(mt) == 0 {
			empty = true
			continue
		}

		if m, err := ParseMediaType(mt); err != nil {
			errs = errors.Join(errs, err)
		} else {
			out = append(out, m)
		}
	}

	if empty {
		out = append(out, NewMediaType("*", "*"))
	}

	return out, errs
}

func ParseMediaType(s string) (MediaType, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		// Empty string is treated as a wildcard
		return NewMediaType("*", "*"), nil
	}

	data := strings.Split(s, ";")
	if len(data) < 1 {
		return MediaType{}, fmt.Errorf("invalid media type %q", s)
	}

	parts := strings.Split(data[0], "/")
	if len(parts) != 2 {
		return MediaType{}, fmt.Errorf("invalid media type %q", s)
	}

	// We ignore the parameters for now
	return NewMediaType(parts[0], parts[1]), nil
}

func (m MediaType) Match(other ...MediaType) bool {
	match := func(a, b string) bool {
		return a == "*" || b == "*" || strings.EqualFold(a, b)
	}

	for _, o := range other {
		if !match(m.Type, o.Type) {
			continue
		}

		if !match(m.Subtype, o.Subtype) {
			continue
		}

		return true
	}

	return false
}

func ValidateAccepts(r *http.Request, want MediaType) error {
	got, err := ParseMediaTypes(r.Header.Get("Accept"))
	if len(got) == 0 && err != nil {
		// Be lenient and only fail if we could not extract any media types
		return Errorf(http.StatusNotAcceptable, "Invalid Accept header: %v", err)
	}

	if want.Match(got...) {
		return nil
	}

	return Errorf(http.StatusNotAcceptable, "Incompatible Accept header: %v. Only %s is accepted", got, want)
}

func ValidateContentType(r *http.Request, want MediaType) error {
	got, err := ParseMediaTypes(r.Header.Get("Content-Type"))
	if len(got) == 0 && err != nil {
		// Be lenient and only fail if we could not extract any media types
		return Errorf(http.StatusNotAcceptable, "Invalid Content-Type header: %v", err)
	}

	if want.Match(got...) {
		return nil
	}

	return Errorf(http.StatusNotAcceptable, "Incompatible Content-Type header: %v. Only %s is accepted", got, want)
}
