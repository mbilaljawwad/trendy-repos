package transform

import (
	"strings"

	"github.com/mbilaljawwad/trendy-repos/internal/extract"
)

type Repo = extract.Repo

type NormalizedRepo struct {
	GitHubId        int64
	Name            string
	HtmlUrl         string
	Description     string
	Language        string
	StargazersCount int64
	OwnerUsername   string
	Topics          string
}

const (
	TRUNCATED_DESCRIPTION_LENGTH = 100
)

/**
* 1. Convert language to lowercase
* 2. Convert owner username to lowercase
* 3. If descriptions empty, set it to "No description"
* 4. Truncate description to 100 characters
* 5. Join topics with comma separated string and lowercase it
 */
func NormalizeData(repos []Repo) []NormalizedRepo {
	var normalizedRepos []NormalizedRepo

	for _, repo := range repos {
		var nRepo NormalizedRepo
		// 1. Convert language to lowercase
		nRepo.Language = strings.ToLower(repo.Language)

		// 2. Convert owner username to lowercase
		nRepo.OwnerUsername = strings.ToLower(repo.OwnerUsername)

		// 3. If descriptions empty, set it to "No description"
		if repo.Description == "" {
			nRepo.Description = "<No description>"
		} else {
			// 4. Truncate description to 100 characters
			if len(repo.Description) > 100 {
				nRepo.Description = repo.Description[:TRUNCATED_DESCRIPTION_LENGTH] + "..."
			}
		}

		// 5. Join topics with comma separated string and lowercase it
		nRepo.Topics = strings.ToLower(strings.Join(repo.Topics, ","))

		normalizedRepos = append(normalizedRepos, nRepo)
	}
	return normalizedRepos
}
