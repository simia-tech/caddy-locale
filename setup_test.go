package locale

import (
	"testing"

	"github.com/caddyserver/caddy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/caddy-locale/method"
)

func TestLocaleParsing(t *testing.T) {
	testFn := func(input string, expectLocales []string, expectMethods []method.Method, expectPathScope, expectCookieName string) func(*testing.T) {
		return func(t *testing.T) {
			localeHandler, err := parseLocale(caddy.NewTestController("http", input))
			require.NoError(t, err)

			assert.Equal(t, expectLocales, localeHandler.AvailableLocales)
			assert.Equal(t, len(expectMethods), len(localeHandler.Methods))
			assert.Equal(t, expectPathScope, localeHandler.PathScope)
			assert.Equal(t, expectCookieName, localeHandler.Configuration.CookieName)
		}
	}

	t.Run("OneLiner", testFn(`locale en de`,
		[]string{"en", "de"}, []method.Method{method.Names["header"]}, "/", "locale"))
	t.Run("PathScope", testFn(`locale en {
			available de
			path /
		}`,
		[]string{"en", "de"}, []method.Method{method.Names["header"]}, "/", "locale"))
	t.Run("DetectMethods", testFn(`locale en de {
			detect cookie header
			cookie language
			path /test
		}`,
		[]string{"en", "de"}, []method.Method{method.Names["cookie"], method.Names["header"]}, "/test", "language"))
}
