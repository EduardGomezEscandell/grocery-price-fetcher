package httputils_test

import (
	"testing"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
	"github.com/stretchr/testify/require"
)

func TestMediatypes(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		Accept  string
		Receive string
		Reject  bool
	}{
		// Matching types
		"Accept matching types":                               {Accept: "application/json", Receive: "application/json"},
		"Accept matching types with a parameter":              {Accept: "application/json; charset=utf-8", Receive: "application/json"},
		"Accept matching types in list":                       {Accept: "application/json", Receive: "application/json,application/css"},
		"Accept matching second type in list":                 {Accept: "application/css", Receive: "application/json,application/css"},
		"Accept matching second type in list with whitespace": {Accept: "application/css", Receive: "application/json,              application/css    , text/html"},
		"Accept matching types in list with a parameter":      {Accept: "application/json; charset=utf-8, text/html", Receive: "application/json"},

		// Full wildcards
		"Accept types with a wildcard":                          {Accept: "*/*", Receive: "application/json"},
		"Accept types with an empty string":                     {Accept: "", Receive: "application/json"},
		"Accept types with a wildcard with a parameter":         {Accept: "*/*; q=0.8", Receive: "application/json"},
		"Accept types with a wildcard in list":                  {Accept: "*/*", Receive: "application/json, application/css"},
		"Accept types with a wildcard in list with a parameter": {Accept: "*/*; q=0.8", Receive: "application/json, application/css"},

		// Subtype wildcards
		"Accept subtype wildcard subtypes":                  {Accept: "application/*", Receive: "application/json"},
		"Accept subtype wildcard subtypes with a parameter": {Accept: "application/*; q=0.8", Receive: "application/json"},
		"Accept subtype wildcard subtypes in list":          {Accept: "application/*", Receive: "application/json, application/css"},

		"Reject mismatching types":                          {Accept: "application/json", Receive: "text/html", Reject: true},
		"Reject mismatching types in list":                  {Accept: "application/json", Receive: "text/html, application/css", Reject: true},
		"Reject mismatching types in list with a parameter": {Accept: "application/json; charset=utf-8, text/html", Receive: "text/html", Reject: true},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			accepts, err := httputils.ParseMediaType(tc.Accept)
			require.NoError(t, err)

			receives, err := httputils.ParseMediaTypes(tc.Receive)
			require.NoError(t, err)

			t.Logf("Accepts: %v", accepts)
			t.Logf("Receives: %v", receives)

			ok := accepts.Match(receives...)
			if tc.Reject {
				require.False(t, ok, "Verification should fail")
				return
			}
			require.True(t, ok, "Verification should pass")
		})
	}
}
