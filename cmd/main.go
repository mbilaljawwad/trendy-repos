package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mbilaljawwad/trendy-repos/internal/config"
	"github.com/mbilaljawwad/trendy-repos/internal/extract"
	"github.com/mbilaljawwad/trendy-repos/internal/http_middleware"
	"github.com/mbilaljawwad/trendy-repos/internal/transform"
)

const (
	PORT = ":8080"
)

func main() {
	fmt.Println("Starting the application...")
	config.InitConfig()

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(http_middleware.CorsMiddleware)

	router.Get("/extract", func(w http.ResponseWriter, req *http.Request) {
		client := extract.NewGithubClient()
		repos := client.FetchTrendingReposFor2025(extract.SortStars, extract.OrderDesc)
		normalizedRepos := transform.NormalizeData(repos)

		fmt.Println("normalizedRepos: ", normalizedRepos)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Extracted and normalized data successfully",
		})
	})

	srv := http.Server{
		Addr:    PORT,
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
