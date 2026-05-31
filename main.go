package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
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

	routes := http.NewServeMux()
	routes.HandleFunc("/api/v1/message", messageHandler)
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
