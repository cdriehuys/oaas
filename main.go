package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cdriehuys/oaas/api"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	apiBaseURL := "/api/v1/"
	api, err := api.NewAPI(ctx, apiBaseURL, "db.sqlite")
	if err != nil {
		panic(err)
	}

	routes := http.NewServeMux()
	routes.Handle(apiBaseURL, api.Routes())

	registerUIRoute(routes)

	server := &http.Server{
		Addr:    ":8000",
		Handler: routes,
	}

	fmt.Printf("Starting server on %q\n", server.Addr)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Println("Server errored:", err)
			}
		}
	}()

	<-ctx.Done()
	cancel()

	shutdownCtx, stopTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopTimeout()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Println("Error shutting down server:", err)
	}

	stopTimeout()

	log.Println("Shutdown successful.")
}
