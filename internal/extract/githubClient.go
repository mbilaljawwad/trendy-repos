package extract

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

const (
	baseURL        = "https://api.github.com"
	searchRepoPath = "search/repositories"
	perPage        = 5
	sinceCreated   = "created:>2025-01-01"
)

// name, html_url, description, language, stargazers_count, owner.login, topics

type RepoSearchResponse struct {
	TotalCount int64  `json:"total_count"`
	Incomplete bool   `json:"incomplete_results"`
	Items      []Repo `json:"items"`
}

type Repo struct {
	GitHubId        int64    `json:"id"`
	Name            string   `json:"name"`
	HtmlUrl         string   `json:"html_url"`
	Description     string   `json:"description"`
	Language        string   `json:"language"`
	StargazersCount int64    `json:"stargazers_count"`
	OwnerUsername   string   `json:"owner.login"`
	Topics          []string `json:"topics"`
}

type GithubClient struct {
	URL    string
	APIKey string
	Client *http.Client
}

type Order string

const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

type Sort string

const (
	SortStars Sort = "stars"
)

func NewGithubClient() *GithubClient {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	return &GithubClient{
		URL:    fmt.Sprintf("%s/%s", baseURL, searchRepoPath),
		APIKey: viper.GetString("GITHUB_PAT"),
		Client: client,
	}
}
