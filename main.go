package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/cdriehuys/oaas/api"
)

//go:embed frontend/dist
var uiFS embed.FS

type Message struct {
	Message string `json:"message"`
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	message := Message{"Hello, World!"}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(message); err != nil {
		log.Println("Failed to encode message:", err)
	}
}

func main() {
	bareUIFiles, _ := fs.Sub(uiFS, "frontend/dist")

	apiBaseURL := "/api/v1/"
	api := api.NewAPI(apiBaseURL)

	routes := http.NewServeMux()
	routes.Handle(apiBaseURL, api.Routes())
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
