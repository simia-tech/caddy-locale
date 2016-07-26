package locale

import (
	"testing"

	"github.com/mholt/caddy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/caddy-locale/method"
)

func TestLocaleParsing(t *testing.T) {
	tests := []struct {
		input                 string
		expectedLocales       []string
		expectedMethods       []method.Method
		expectedPathScope     string
		expectedConfiguration *method.Configuration
	}{
		{`locale en de`, []string{"en", "de"}, []method.Method{method.Names["header"]}, "/", &method.Configuration{}},
		{`locale en {
        available de
        path /
      }`, []string{"en", "de"}, []method.Method{method.Names["header"]}, "/", &method.Configuration{}},
		{`locale en de {
        detect cookie header
        cookie language
        path /test
      }`, []string{"en", "de"}, []method.Method{method.Names["cookie"], method.Names["header"]}, "/test",
			&method.Configuration{CookieName: "language"}},
	}

	for index, test := range tests {
		localeHandler, err := parseLocale(caddy.NewTestController("http", test.input))
		require.NoError(t, err, "test #%d", index)

		assert.Equal(t, test.expectedLocales, localeHandler.AvailableLocales, "test #%d", index)
		assert.Equal(t, len(test.expectedMethods), len(localeHandler.Methods), "test #%d", index)
	}
}
