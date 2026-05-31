package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed frontend/dist
var uiFS embed.FS

func main() {
	bareUIFiles, _ := fs.Sub(uiFS, "frontend/dist")

	routes := http.NewServeMux()
	routes.Handle("/", http.FileServerFS(bareUIFiles))

	server := &http.Server{
		Addr:    ":8000",
		Handler: routes,
	}

	fmt.Printf("Starting server on %q\n", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
