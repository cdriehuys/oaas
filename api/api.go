package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/cdriehuys/oaas/api/internal/repositories"
	"github.com/cdriehuys/oaas/api/internal/server"
)

type API struct {
	baseURL string

	server *server.Server
}

func NewAPI(ctx context.Context, baseURL string, database string) (*API, error) {
	repo, err := repositories.NewSQLiteTodoRepo(ctx, database)
	if err != nil {
		return nil, fmt.Errorf("creating todo repository: %v", err)
	}

	api := &API{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		server:  server.NewServer(repo),
	}

	return api, nil
}

func handleRequestError(w http.ResponseWriter, r *http.Request, err error) {
	errType := "/probs/bad-request"
	title := "Bad Request"
	details := err.Error()

	rep := server.ProblemDetails{Type: &errType, Title: &title, Detail: &details}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(http.StatusBadRequest)

	json.NewEncoder(w).Encode(rep)
}

func handleResponseError(w http.ResponseWriter, r *http.Request, err error) {
	log.Println("Unexpected server error:", err)

	errType := "/probs/server-error"
	title := http.StatusText(http.StatusInternalServerError)

	rep := server.ProblemDetails{Type: &errType, Title: &title}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(http.StatusInternalServerError)

	json.NewEncoder(w).Encode(rep)
}

func (a *API) Routes() http.Handler {
	strictHandler := server.NewStrictHandlerWithOptions(a.server, nil, server.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  handleRequestError,
		ResponseErrorHandlerFunc: handleResponseError,
	})

	mux := http.NewServeMux()

	return server.HandlerFromMuxWithBaseURL(strictHandler, mux, a.baseURL)
}
