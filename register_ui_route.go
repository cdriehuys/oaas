//go:build !no_ui

package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed frontend/dist
var uiFS embed.FS

func registerUIRoute(mux *http.ServeMux) {
	bareUIFiles, err := fs.Sub(uiFS, "frontend/dist")
	if err != nil {
		panic(err)
	}

	mux.Handle("/", http.FileServerFS(bareUIFiles))
}
