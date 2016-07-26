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
		availableLocales     []string
		methods              []method.Method
		pathScope            string
		path                 string
		acceptLanguageHeader string
		expectedLocaleHeader string
	}{
		{[]string{"en"}, []method.Method{method.Names["header"]}, "/", "/", "", "en"},
		{[]string{"en", "en-GB"}, []method.Method{method.Names["header"]}, "/", "/", "en-GB,en", "en-GB"},
		{[]string{"en", "de"}, []method.Method{method.Names["header"]}, "/", "/", "de,en;q=0.8,en-GB;q=0.6", "de"},
		{[]string{"en"}, []method.Method{method.Names["header"]}, "/test", "/other/path", "", ""},
	}

	for index, test := range tests {
		locale := Middleware{
			Next:             httpserver.HandlerFunc(contentHandler),
			AvailableLocales: test.availableLocales,
			Methods:          test.methods,
			PathScope:        test.pathScope,
		}

		request, err := http.NewRequest("GET", "/", nil)
		require.NoError(t, err, "test #%d", index)
		request.Header.Set("Accept-Language", test.acceptLanguageHeader)

		responseRecorder := httpserver.NewResponseRecorder(httptest.NewRecorder())

		_, err = locale.ServeHTTP(responseRecorder, request)
		require.NoError(t, err, "test #%d", index)

		assert.Equal(t, test.expectedLocaleHeader, request.Header.Get("Detected-Locale"), "test #%d", index)
	}
}

func contentHandler(_ http.ResponseWriter, _ *http.Request) (int, error) {
	return http.StatusOK, nil
}
