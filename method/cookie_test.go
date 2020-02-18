package method

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCookieParsing(t *testing.T) {
	cookie := Names["cookie"]
	configuration := &Configuration{CookieName: "locale"}

	testFn := func(name, value string, expectLocales []string) (string, func(*testing.T)) {
		return name, func(t *testing.T) {
			request, _ := http.NewRequest("GET", "/", nil)
			if name != "" {
				request.Header.Set("Cookie", (&http.Cookie{Name: name, Value: value}).String())
			}

			locales := cookie(request, configuration)
			assert.Equal(t, expectLocales, locales)
		}
	}

	t.Run(testFn("", "", []string{}))
	t.Run(testFn("locale", "en", []string{"en"}))
}
