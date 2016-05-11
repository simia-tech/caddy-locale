package locale

import (
	"net/http"
	"strings"

	"github.com/mholt/caddy/middleware"

	"github.com/simia-tech/caddy-locale/method"
)

// Middleware is a middleware to detect the user's locale.
type Middleware struct {
	Next             middleware.Handler
	AvailableLocales []string
	Methods          []method.Method
	PathScope        string
	Configuration    *method.Configuration
}

// ServeHTTP implements the middleware.Handler interface.
func (l *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if !middleware.Path(r.URL.Path).Matches(l.PathScope) {
		return l.Next.ServeHTTP(w, r)
	}

	candidates := []string{}
	for _, method := range l.Methods {
		candidates = append(candidates, method(r, l.Configuration)...)
	}

	locale := l.firstValid(candidates)
	if locale == "" {
		locale = l.defaultLocale()
	}
	r.Header.Set("Detected-Locale", locale)

	if rr, ok := w.(*middleware.ResponseRecorder); ok && rr.Replacer != nil {
		rr.Replacer.Set("locale", locale)
	}

	return l.Next.ServeHTTP(w, r)
}

func (l *Middleware) defaultLocale() string {
	return l.AvailableLocales[0]
}

func (l *Middleware) firstValid(candidates []string) string {
	for _, candidate := range candidates {
		candidate = strings.TrimSpace(candidate)
		if l.isValid(candidate) {
			return candidate
		}
	}
	return ""
}

func (l *Middleware) isValid(locale string) bool {
	for _, validLocale := range l.AvailableLocales {
		if locale == validLocale {
			return true
		}
	}
	return false
}
