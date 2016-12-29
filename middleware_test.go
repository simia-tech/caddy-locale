package locale

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/caddy-locale/method"
)

func TestMiddleware(t *testing.T) {
	tests := []struct {
		name                 string
		availableLocales     []string
		methods              []method.Method
		pathScope            string
		path                 string
		acceptLanguageHeader string
		expectedLocaleHeader string
	}{
		{"Single", []string{"en"}, []method.Method{method.Names["header"]}, "/", "/", "", "en"},
		{"Multiple", []string{"en", "en-GB"}, []method.Method{method.Names["header"]}, "/", "/", "en-GB,en", "en-GB"},
		{"Weights", []string{"en", "de"}, []method.Method{method.Names["header"]}, "/", "/", "de,en;q=0.8,en-GB;q=0.6", "de"},
		{"Path", []string{"en"}, []method.Method{method.Names["header"]}, "/test", "/other/path", "", ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			locale := Middleware{
				Next:             httpserver.HandlerFunc(contentHandler),
				AvailableLocales: test.availableLocales,
				Methods:          test.methods,
				PathScope:        test.pathScope,
			}

			request, err := http.NewRequest("GET", "/", nil)
			require.NoError(t, err)
			request.Header.Set("Accept-Language", test.acceptLanguageHeader)

			responseRecorder := httpserver.NewResponseRecorder(httptest.NewRecorder())

			_, err = locale.ServeHTTP(responseRecorder, request)
			require.NoError(t, err)

			assert.Equal(t, test.expectedLocaleHeader, request.Header.Get("Detected-Locale"))
		})
	}
}

func contentHandler(_ http.ResponseWriter, _ *http.Request) (int, error) {
	return http.StatusOK, nil
}
