package locale

import (
	"testing"

	"github.com/caddyserver/caddy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/caddy-locale/method"
)

func TestLocaleParsing(t *testing.T) {
	tests := []struct {
		name                  string
		input                 string
		expectedLocales       []string
		expectedMethods       []method.Method
		expectedPathScope     string
		expectedConfiguration *method.Configuration
	}{
		{"OneLiner", `locale en de`, []string{"en", "de"}, []method.Method{method.Names["header"]}, "/", &method.Configuration{}},
		{"PathScope", `locale en {
        available de
        path /
      }`, []string{"en", "de"}, []method.Method{method.Names["header"]}, "/", &method.Configuration{}},
		{"DetectMethods", `locale en de {
        detect cookie header
        cookie language
        path /test
      }`, []string{"en", "de"}, []method.Method{method.Names["cookie"], method.Names["header"]}, "/test",
			&method.Configuration{CookieName: "language"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			localeHandler, err := parseLocale(caddy.NewTestController("http", test.input))
			require.NoError(t, err)

			assert.Equal(t, test.expectedLocales, localeHandler.AvailableLocales)
			assert.Equal(t, len(test.expectedMethods), len(localeHandler.Methods))
		})
	}
}
