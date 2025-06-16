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
