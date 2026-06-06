package api

import (
	"context"
	"fmt"
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

func (a *API) Routes() http.Handler {
	strictHandler := server.NewStrictHandler(a.server, nil)
	mux := http.NewServeMux()

	return server.HandlerFromMuxWithBaseURL(strictHandler, mux, a.baseURL)
}
