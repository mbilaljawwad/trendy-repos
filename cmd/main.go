package main

import (
	"fmt"

	"github.com/mbilaljawwad/trendy-repos/internal/config"
	"github.com/mbilaljawwad/trendy-repos/internal/extract"
)

func main() {
	fmt.Println("Starting the application...")
	config.InitConfig()

	client := extract.NewGithubClient()
	client.FetchTrendingReposFor2025(extract.SortStars, extract.OrderDesc)
}
