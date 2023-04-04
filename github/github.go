package github

import "star-trend/config"

var Gh Github

type Github struct {
	token   string // If the calling interface is restricted, a token can be used token: https://docs.github.com/rest
	perPage int    // The number of results per page (max 100).
}

func init() {
	Gh = Github{
		token:   config.Cfg.GithubToken,
		perPage: config.Cfg.GithubPerPage,
	}
}

func (gh *Github) HasToken() bool {
	return gh.GetToken() != ""
}

func (gh *Github) GetToken() string {
	return gh.token
}

func (gh *Github) GetPerPage() int {
	return gh.perPage
}
