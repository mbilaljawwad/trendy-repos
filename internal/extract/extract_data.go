package extract

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (c *GithubClient) FetchTrendingReposFor2025(
	sort Sort,
	order Order,
) []Repo {
	fmt.Println("Fetching trending repos...")
	searchUrl := fmt.Sprintf("%s?q=%s&sort=%s&order=%s&per_page=%d", c.URL, sinceCreated, sort, order, perPage)
	req, err := http.NewRequest("GET", searchUrl, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	ensureGithubHeader(req, c.APIKey)
	resp, err := c.Client.Do(req)
	if err != nil {
		log.Fatalf("Error fetching trending repos: %v", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
		return nil
	}

	var repoSearchResp RepoSearchResponse
	err = json.Unmarshal(body, &repoSearchResp)
	if err != nil {
		log.Fatalf("Error parsing response body: %v", err)
		return nil
	}

	var repos []Repo
	for _, repo := range repoSearchResp.Items {
		repos = append(repos, Repo{
			GitHubId:        repo.GitHubId,
			Name:            repo.Name,
			HtmlUrl:         repo.HtmlUrl,
			Description:     repo.Description,
			Language:        repo.Language,
			StargazersCount: repo.StargazersCount,
			OwnerUsername:   repo.OwnerUsername,
			Topics:          repo.Topics,
		})
	}
	return repos
}

func ensureGithubHeader(req *http.Request, apiKey string) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
}
