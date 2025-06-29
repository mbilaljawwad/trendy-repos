package transform

import (
	"strings"

	"github.com/mbilaljawwad/trendy-repos/internal/extract"
	"github.com/mbilaljawwad/trendy-repos/internal/pkg/utils"
)

type Repo = extract.Repo

type NormalizedRepo struct {
	GitHubId        int64  `db:"github_id"`
	Name            string `db:"name"`
	HtmlUrl         string `db:"url"`
	Description     string `db:"description"`
	Language        string `db:"language"`
	StargazersCount int64  `db:"stars_count"`
	OwnerUsername   string `db:"owner_username"`
	Topics          string `db:"topics"`
}

const (
	TRUNCATED_DESCRIPTION_LENGTH = 100
	MAX_STRING_SIZE              = 1000
)

/**
* 1. Convert language to lowercase
* 2. Convert owner username to lowercase
* 3. If descriptions empty, set it to "No description"
* 4. Truncate description to 100 characters
* 5. Join topics with comma separated string and lowercase it
* 6. Sanitize all strings to ensure valid UTF-8 encoding
 */
func NormalizeData(repos []Repo) []NormalizedRepo {
	var normalizedRepos []NormalizedRepo

	for _, repo := range repos {

		var nRepo NormalizedRepo

		nRepo.GitHubId = repo.GitHubId
		nRepo.Name, _ = utils.ConvertToUTF8(repo.Name)
		nRepo.HtmlUrl, _ = utils.ConvertToUTF8(repo.HtmlUrl)
		nRepo.StargazersCount = repo.StargazersCount
		// 1. Convert language to lowercase
		nRepo.Language, _ = utils.ConvertToUTF8(strings.ToLower(repo.Language))

		// 2. Convert owner username to lowercase
		nRepo.OwnerUsername, _ = utils.ConvertToUTF8(strings.ToLower(repo.OwnerUsername))

		// 3. If descriptions empty, set it to "No description"
		if repo.Description == "" {
			nRepo.Description = "No description"
		} else {
			// 4. Truncate description to 100 characters
			description, _ := utils.ConvertToUTF8(repo.Description)
			if len(description) > TRUNCATED_DESCRIPTION_LENGTH {
				nRepo.Description = description[:TRUNCATED_DESCRIPTION_LENGTH] + "..."
			} else {
				nRepo.Description = description
			}
		}

		// 5. Join topics with comma separated string and lowercase it
		topicsStr := strings.Join(repo.Topics, ",")
		nRepo.Topics, _ = utils.ConvertToUTF8(strings.ToLower(topicsStr))

		normalizedRepos = append(normalizedRepos, nRepo)
	}
	return normalizedRepos
}
