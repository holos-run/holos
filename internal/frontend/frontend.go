package frontend

import (
	"bytes"
	"embed"
	"io/fs"
	"net/http"
	"text/template"
	"time"
)

// Path needs to match <base href="/ui/"> in index.html and baseHref in angular.json
const Path = "/ui/"

// Output must be the relative path to where the frontend build too places the
// output index.html file.  Tip: use the holos server frontend ls command to
// list the embedded file system.
// Refer to: https://angular.io/guide/workspace-config#output-path-configuration
// This should be the base output path with the browser field set to "ui" in angular.json
const OutputPath = "holos/dist/holos"

//go:embed all:holos/dist/holos/ui
var Dist embed.FS

//go:embed holos/dist/holos/ui/index.html
var spaIndexHtml []byte

// Root returns the content root subdirectory
func Root() fs.FS {
	sub, err := fs.Sub(Dist, OutputPath)
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
func NewSPAFileServer(issuer string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		spaIndexTemplate, err := template.New("index.html").Parse(string(spaIndexHtml))
		if err != nil {
			panic(err)
		}
		config := appConfig{OIDCIssuer: issuer}
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
}

// appConfig represents app configuration originating from the runtime
// deployment, pass to the front end app via template values in the single page
// app index.html.
type appConfig struct {
	OIDCIssuer string
}
