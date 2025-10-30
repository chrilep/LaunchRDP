package web

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed index.html
var indexHTML []byte

//go:embed all:static
var staticFiles embed.FS

// UI is a struct that holds the embedded UI files.
type UI struct {
	staticFS http.FileSystem
}

// NewUI creates a new UI instance.
func NewUI() (*UI, error) {
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return nil, err
	}
	return &UI{staticFS: http.FS(staticFS)}, nil
}

// ServeHTTP serves the UI files.
func (u *UI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/static/") {
		http.StripPrefix("/static/", http.FileServer(u.staticFS)).ServeHTTP(w, r)
		return
	}

	http.NotFound(w, r)
}