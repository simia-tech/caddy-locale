package method

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaderParsing(t *testing.T) {
	header := Names["header"]
	request, _ := http.NewRequest("GET", "/", nil)

	testFn := func(value string, expectLocales []string) (string, func(*testing.T)) {
		return value, func(t *testing.T) {
			request.Header.Set("Accept-Language", value)

			locales := header(request, nil)
			assert.Equal(t, expectLocales, locales)
		}
	}

	t.Run(testFn("de,en;q=0.8,en-GB;q=0.6", []string{"de", "en", "en-GB"}))
	t.Run(testFn("de;q=0.2,en;q=0.8,en-GB;q=0.6", []string{"en", "en-GB", "de"}))
	t.Run(testFn("de,,en-GB;q=0.6", []string{"de", "en-GB"}))
	t.Run(testFn("en; q=0.8, de", []string{"de", "en"}))
}
