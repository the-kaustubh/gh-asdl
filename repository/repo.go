package repository

import "fmt"

var (
	HOST = "api.github.com"
)

type Repo struct {
	Username string
	Name     string
}

func (r *Repo) GetBaseRepoUrl() string {
	return fmt.Sprintf("%s/repos/%s/%s", HOST, r.Username, r.Name)
}

func (r *Repo) GetFullRepoURL() string {
	return fmt.Sprintf("https://%s/releases", r.GetBaseRepoUrl())
}

func (r *Repo) GetAssetUrlWithToken(asset_id, token string) string {
	return fmt.Sprintf("https://%s:@%s/releases/assets/%s",
		token,
		r.GetBaseRepoUrl(),
		asset_id)
}
