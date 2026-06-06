package api

import (
	"net/http"
	"strings"

	"github.com/cdriehuys/oaas/api/internal/server"
)

type API struct {
	baseURL string

	server *server.Server
}

func NewAPI(baseURL string) *API {
	return &API{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		server:  &server.Server{}}
}

func (a *API) Routes() http.Handler {
	strictHandler := server.NewStrictHandler(a.server, nil)
	mux := http.NewServeMux()

	return server.HandlerFromMuxWithBaseURL(strictHandler, mux, a.baseURL)
}
