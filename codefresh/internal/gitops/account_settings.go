package gitops

import (
	"fmt"
)

const (
	// Git providers enum from https://github.com/codefresh-io/argo-platform/blob/90f86de326422ca3bd1f64ca5dd26aeedf985e3e/libs/ql/schema/entities/common/integration.graphql#L200
	GitProviderGitHub          string = "GITHUB"
	GitProviderGerrit          string = "GERRIT"
	GitProviderGitlab          string = "GITLAB"
	GitProviderBitbucket       string = "BITBUCKET"
	GitProviderBitbucketServer string = "BITBUCKET_SERVER"
)

func GetSupportedGitProvidersList() []string {
	return []string{GitProviderGitHub, GitProviderGerrit, GitProviderGitlab, GitProviderBitbucket, GitProviderBitbucketServer}
}

// Matching implementation for https://github.com/codefresh-io/argo-platform/blob/3c6af5b5cbb29aef58ef6617e71159e882987f5c/libs/git/src/helpers.ts#L37.
// Must be updated accordingly
func GetDefaultAPIUrlForProvider(gitProvider string) (*string, error) {

	defaultApiUrlProvider := map[string]string{
		GitProviderGitHub:    "https://api.github.com",
		GitProviderGitlab:    "https://gitlab.com/api/v4",
		GitProviderBitbucket: "https://api.bitbucket.org/2.0",
		GitProviderGerrit:    "https://gerrit-review.googlesource.com/a",
	}

	if val, ok := defaultApiUrlProvider[gitProvider]; ok {
		return &val, nil
	}

	return nil, fmt.Errorf("no default API URL for provider %s can be found. For self hosted git providers URL must be provided explicitly", gitProvider)
}
