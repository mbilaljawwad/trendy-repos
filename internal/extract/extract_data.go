package extract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// sanitizeResponseBody cleans invalid UTF-8 sequences from HTTP response
func sanitizeResponseBody(body []byte) []byte {
	// Remove common problematic byte sequences before JSON parsing
	problemPatterns := [][]byte{
		{0xe7, 0xbd, 0x2e}, // The specific sequence causing the error
		{0xff, 0xfe},       // BOM markers
		{0xfe, 0xff},       // BOM markers
		{0xef, 0xbb, 0xbf}, // UTF-8 BOM
		{0x00},             // Null bytes
	}

	cleanBody := body
	for _, pattern := range problemPatterns {
		cleanBody = bytes.ReplaceAll(cleanBody, pattern, []byte{})
	}

	// Convert to valid UTF-8
	return []byte(strings.ToValidUTF8(string(cleanBody), ""))
}

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

	// Sanitize the response body before JSON unmarshaling
	cleanBody := sanitizeResponseBody(body)

	var repoSearchResp RepoSearchResponse
	err = json.Unmarshal(cleanBody, &repoSearchResp)
	if err != nil {
		log.Fatalf("Error parsing response body: %v", err)
		return nil
	}

	var repos []Repo
	for _, repo := range repoSearchResp.Items {
		r := Repo{
			GitHubId:        repo.GitHubId,
			Name:            repo.Name,
			HtmlUrl:         repo.HtmlUrl,
			Description:     repo.Description,
			Language:        repo.Language,
			StargazersCount: repo.StargazersCount,
			OwnerUsername:   repo.Owner.Login,
			Topics:          repo.Topics,
		}
		repos = append(repos, r)
	}
	return repos
}

func ensureGithubHeader(req *http.Request, apiKey string) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
}
