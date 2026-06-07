//go:build !no_ui

package main

import (
	"embed"
	"io/fs"
	"net/http"
)

// uiFS is a file system that falls back to "index.html" if opening a file fails. When used with
// [http.FileServerFS], this allows for serving an SPA with a fallback for full page navigation.
type uiFS struct {
	files fs.FS
}

func (ui *uiFS) Open(name string) (fs.File, error) {
	file, err := ui.files.Open(name)
	if err == nil {
		return file, nil
	}

	return ui.files.Open("index.html")
}

//go:embed frontend/dist
var uiFiles embed.FS

func registerUIRoute(mux *http.ServeMux) {
	bareUIFiles, err := fs.Sub(uiFiles, "frontend/dist")
	if err != nil {
		panic(err)
	}

	mux.Handle("/", http.FileServerFS(&uiFS{bareUIFiles}))
}
