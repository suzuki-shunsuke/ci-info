package domain

type Params struct {
	Owner            string
	Repo             string
	SHA              string
	Dir              string
	PRNum            int
	GitHubToken      string
	GitHubAPIURL     string
	GitHubGraphQLURL string
	LogLevel         string
	Prefix           string
}
