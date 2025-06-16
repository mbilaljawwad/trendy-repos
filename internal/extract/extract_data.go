package extract

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func (c *GithubClient) FetchTrendingReposFor2025(
	sort Sort,
	order Order,
) {
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
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	fmt.Println(string(body))

}

func ensureGithubHeader(req *http.Request, apiKey string) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
}
