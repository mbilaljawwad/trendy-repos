package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/mbilaljawwad/trendy-repos/internal/extract"
	"github.com/mbilaljawwad/trendy-repos/internal/transform"
)

type APIHandler struct {
	database *sqlx.DB
}

func NewAPIHandler(db *sqlx.DB) *APIHandler {
	return &APIHandler{
		database: db,
	}
}

func (apiHandler *APIHandler) InitiateProcess(w http.ResponseWriter, r *http.Request) {
	client := extract.NewGithubClient()
	repos := client.FetchTrendingReposFor2025(extract.SortStars, extract.OrderDesc)
	cleanedRepos := transform.NormalizeData(repos)
	_, err := loadData(apiHandler.database, cleanedRepos)
	if err != nil {
		http.Error(w, "Failed to load data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Extracted and normalized data successfully",
	})
}

func loadData(db *sqlx.DB, repos []transform.NormalizedRepo) (bool, error) {
	log.Println("Loading data into database")
	query := `INSERT INTO repositories (
		github_id,
		name,
		description,
		url,
		stars_count,
		language,
		topics
	) VALUES (
		:github_id,
		:name,
		:description,
		:url,
		:stars_count,
		:language,
		:topics
	)`

	_, err := db.NamedExec(query, repos)
	if err != nil {
		log.Printf("Error loading data into database: %v", err)
		return false, err
	}

	return true, nil
}
