package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Server is our server powering the API.
// It houses a connection to our database,
// a default context, the logger, and the handler
// for the server.
type Server struct {
	DB     *gorm.DB
	ctx    context.Context
	logger *slog.Logger
	http.Handler
}

// NewServer creates a new server with the given database.
// Along with the database.
//
// Use this when creating new servers because it will
// perform all the necessary server setup that you may
// not do otherwise.
func NewServer(db *gorm.DB) *Server {
	srv := &Server{
		DB:     db,
		ctx:    context.Background(),
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	srv.registerRoutes()

	return srv
}

// registerRoutes registers all the routes for the server.
func (s *Server) registerRoutes() {
	router := mux.NewRouter()

	router.PathPrefix("/user").Handler(s.UserRoutes())
	router.PathPrefix("/blog").Handler(s.BlogRoutes())
	router.PathPrefix("/tags").Handler(s.TagRoutes())
	s.Handler = router
}

// Run starts the server on the given host.
// It will panic if the server fails to run.
// Host is the full host including port.
func (s *Server) Run(host string) {
	fmt.Printf("\nServer running on %s\n", host)
	err := http.ListenAndServe(host, s)
	if err != nil {
		// Panic because running failed
		panic(err)
	}
}
