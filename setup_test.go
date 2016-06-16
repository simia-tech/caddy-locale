package locale

import (
	"reflect"
	"testing"

	"github.com/mholt/caddy"

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
		localeHandler, err := parseLocale(caddy.NewTestController(test.input))
		if err != nil {
			t.Fatalf("test %d: unexpected error: %s", index, err)
		}

		if !reflect.DeepEqual(localeHandler.AvailableLocales, test.expectedLocales) {
			t.Errorf("test %d: expected handler to have available locales %#v, got: %#v", index,
				test.expectedLocales, localeHandler.AvailableLocales)
		}
		if am, em := len(localeHandler.Methods), len(test.expectedMethods); em != am {
			t.Errorf("test %d: expected handler to have %d detect methods, got: %d", index, em, am)
		}
	}
}
