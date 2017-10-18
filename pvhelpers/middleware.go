package pvhelpers

import (
	"github.com/codegangsta/negroni"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type NoStaticLogger struct {
	// Logger inherits from log.Logger used to log messages with the Logger middleware
	*log.Logger
}

// NewLogger returns a new Logger instance
func NewNoStaticLogger() *NoStaticLogger {
	return &NoStaticLogger{log.New(os.Stdout, "[negroni] ", 0)}
}

func (l *NoStaticLogger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	track := true
	if strings.HasPrefix(r.URL.Path, "/static/") {
		track = false
	}

	if track {
		if Cfg.App.LogLevel == "Verbose" {
			l.Printf("Started %s %s", r.Method, r.URL.Path)
		}
	}

	next(rw, r)

	if track {
		if Cfg.App.LogLevel == "Verbose" {
			res := rw.(negroni.ResponseWriter)
			l.Printf("Completed %v %s in %v", res.Status(), http.StatusText(res.Status()), time.Since(start))
		}
	}
}
