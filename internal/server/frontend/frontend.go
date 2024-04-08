package frontend

import (
	"bytes"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"text/template"
	"time"
)

const Path = "/app/"

//go:embed all:dist
var content embed.FS

//go:embed dist/app/index.html
var spaIndexHtml []byte

// Root returns the content root subdirectory
func Root() fs.FS {
	sub, err := fs.Sub(content, "dist")
	if err != nil {
		panic(err)
	}
	return sub
}

// responseRecorder records http.StatusNotFound responses without
// writing the response to the client so the SPA can be served instead.
type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (rw *responseRecorder) Write(p []byte) (int, error) {
	if rw.status == http.StatusNotFound {
		return len(p), nil
	}
	return rw.ResponseWriter.Write(p)
}

func (rw *responseRecorder) WriteHeader(status int) {
	rw.status = status
	if status != http.StatusNotFound {
		rw.ResponseWriter.WriteHeader(status)
	}
}

// SPAFileServer returns a http.Handler that calls h and falls back to serving
// the spa app container when h writes header http.StatusNotFound. Useful when
// the user's browser refreshes the page and the url is handled by a spa like
// react router.
func SPAFileServer(h http.Handler) http.Handler {
	spaIndexTemplate, err := template.New("index.html").Parse(string(spaIndexHtml))
	if err != nil {
		panic(err)
	}
	config := newAppConfig()
	buf := new(bytes.Buffer)
	err = spaIndexTemplate.Execute(buf, config)
	if err != nil {
		panic(err)
	}
	spaIndex := buf.Bytes()
	modTime := time.Time{}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == Path || r.URL.Path == Path+"index.html" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			http.ServeContent(w, r, "", modTime, bytes.NewReader(spaIndex))
			return
		}

		// Serve static assets
		discard404ResponseRecorder := &responseRecorder{ResponseWriter: w}
		h.ServeHTTP(discard404ResponseRecorder, r)
		if discard404ResponseRecorder.status == http.StatusNotFound {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			http.ServeContent(w, r, "", modTime, bytes.NewReader(spaIndex))
		}
	})
}

// appConfig represents app configuration originating from the runtime
// deployment, pass to the front end app via template values in the single page
// app index.html.
type appConfig struct {
	OIDCIssuer string
}

// checkWarnings issues warnings and guidance about config env vars necessary for the front end to function.
func (v appConfig) checkWarnings() {
	if v.OIDCIssuer == "" {
		slog.Warn("OIDCIssuer is empty, set OIDC_ISSUER env var")
	}
}

func newAppConfig() appConfig {
	ac := appConfig{
		OIDCIssuer: os.Getenv("OIDC_ISSUER"),
	}
	ac.checkWarnings()
	return ac
}
