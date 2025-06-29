package extract

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

const (
	searchRepoPath = "search/repositories"
	perPage        = 5
	sinceCreated   = "created:>2025-01-01"
)

type Owner struct {
	Login string `json:"login"`
}

type RepoSearchResponse struct {
	TotalCount int64  `json:"total_count"`
	Incomplete bool   `json:"incomplete_results"`
	Items      []Repo `json:"items"`
}

type Repo struct {
	GitHubId        int64    `json:"id,omitempty"`
	Name            string   `json:"name"`
	HtmlUrl         string   `json:"html_url"`
	Description     string   `json:"description"`
	Language        string   `json:"language"`
	StargazersCount int64    `json:"stargazers_count"`
	Owner           Owner    `json:"owner"`
	OwnerUsername   string   `json:"-"` // This will be populated from Owner.Login
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
	baseURL := viper.GetString("GITHUB_API_URL")
	return &GithubClient{
		URL:    fmt.Sprintf("%s/%s", baseURL, searchRepoPath),
		APIKey: viper.GetString("GITHUB_PAT"),
		Client: client,
	}
}
