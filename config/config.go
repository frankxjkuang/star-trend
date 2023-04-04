package config

import (
	"github.com/caarlos0/env/v7"
)

var Cfg Config

type Config struct {
	GithubToken   string `env:"GITHUB_TOKEN"` // https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/about-authentication-to-github
	GithubPerPage int    `env:"GITHUB_PER_PAGE" envDefault:"100"`
	Port          string `env:"PORT" envDefault:"9000"` // 监听端口
}

func init() {
	err := env.Parse(&Cfg)
	if err != nil {
		panic(err)
	}
}
