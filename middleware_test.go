package locale

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caddyserver/caddy/caddyhttp/httpserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/caddy-locale/method"
)

func TestMiddleware(t *testing.T) {
	testFn := func(availableLocales []string, methods []method.Method, pathScope, path, acceptLanguageHeader, expectLocale string) func(*testing.T) {
		return func(t *testing.T) {
			locale := Middleware{
				Next:             httpserver.HandlerFunc(contentHandler),
				AvailableLocales: availableLocales,
				Methods:          methods,
				PathScope:        pathScope,
			}

			request, err := http.NewRequest("GET", "/", nil)
			require.NoError(t, err)
			request.Header.Set("Accept-Language", acceptLanguageHeader)

			responseRecorder := httpserver.NewResponseRecorder(httptest.NewRecorder())

			_, err = locale.ServeHTTP(responseRecorder, request)
			require.NoError(t, err)

			assert.Equal(t, expectLocale, request.Header.Get("Detected-Locale"))
		}
	}

	t.Run("Single", testFn([]string{"en"}, []method.Method{method.Names["header"]}, "/", "/", "", "en"))
	t.Run("Multiple", testFn([]string{"en", "en-GB"}, []method.Method{method.Names["header"]}, "/", "/", "en-GB,en", "en-GB"))
	t.Run("Weights", testFn([]string{"en", "de"}, []method.Method{method.Names["header"]}, "/", "/", "de,en;q=0.8,en-GB;q=0.6", "de"))
	t.Run("Path", testFn([]string{"en"}, []method.Method{method.Names["header"]}, "/test", "/other/path", "", ""))
	t.Run("CaseSensitivity", testFn([]string{"en", "en-gb"}, []method.Method{method.Names["header"]}, "/", "/", "en-GB,en", "en-gb"))
}

func contentHandler(_ http.ResponseWriter, _ *http.Request) (int, error) {
	return http.StatusOK, nil
}
