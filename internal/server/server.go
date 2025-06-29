package server

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/mbilaljawwad/trendy-repos/internal/datastore"
	"github.com/mbilaljawwad/trendy-repos/internal/handlers"
	"github.com/mbilaljawwad/trendy-repos/internal/http_middleware"
)

const (
	PORT = ":8080"
)

type AppServer struct {
	Server     *http.Server
	DB         *sqlx.DB
	ApiHandler *handlers.APIHandler
}

func NewAppServer(ctx context.Context) *AppServer {

	db := datastore.NewDataStore(ctx)

	srv := &AppServer{
		DB: db,
		Server: &http.Server{
			Addr: PORT,
		},
		ApiHandler: handlers.NewAPIHandler(db),
	}

	router := srv.initializeRouter()
	srv.Server.Handler = router

	return srv
}

func (app *AppServer) Start() {
	if err := app.Server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func (app *AppServer) initializeRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(http_middleware.CorsMiddleware)

	r.Get("/start", app.ApiHandler.InitiateProcess)
	return r
}
