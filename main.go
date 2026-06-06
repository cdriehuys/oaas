package main

import (
	"fmt"
	"net/http"

	"github.com/cdriehuys/oaas/api"
)

func main() {
	apiBaseURL := "/api/v1/"
	api := api.NewAPI(apiBaseURL)

	routes := http.NewServeMux()
	routes.Handle(apiBaseURL, api.Routes())

	registerUIRoute(routes)

	server := &http.Server{
		Addr:    ":8000",
		Handler: routes,
	}

	fmt.Printf("Starting server on %q\n", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
