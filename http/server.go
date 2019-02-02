package http

import (
	"net/http"

	"github.com/holman_dw/ukjent/store"
)

type server struct {
	store store.Store
}

// Run builds a new HTTP router
func Run(s store.Store) error {
	server := server{s}
	mux := server.registerRoutes()
	if err := http.ListenAndServe(":8080", mux); err != nil {
		return err
	}
	return nil
}
