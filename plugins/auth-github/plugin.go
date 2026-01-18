package main

import (
	"tgp/core/data"
	"tgp/core/i18n"
	"tgp/core/plugin"
	"tgp/internal/github"
	"tgp/plugins/auth-github/internal/config"
	"tgp/plugins/auth-github/internal/handler"
)

// AuthGithubPlugin реализует интерфейс Plugin.
type AuthGithubPlugin struct{}

// Execute выполняет основную логику плагина.
func (p *AuthGithubPlugin) Execute(rootDir string, request data.Storage, path ...string) (response data.Storage, err error) {

	return handler.HandleExecute(request, path...)
}

// Info возвращает информацию о плагине.
func (p *AuthGithubPlugin) Info() (info plugin.Info, err error) {

	info = plugin.Info{
		Name:        "auth-github",
		Description: i18n.Msg("Authentication transformer that enriches execution chain with GitHub account data"),
		Author:      "AlexK <seniorGolang@gmail.com>",
		License:     "MIT",
		Category:    "utility",
		Kind:        "pre",
		Silent:      true,
		Always:      true,
		AllowedHosts: []string{
			github.AllowedHostGitHubAPI,
		},
		AllowedPaths: map[string]string{
			github.AllowedPathGitHub: github.AllowedPathModeWrite,
			config.AllowedPathAuth:   github.AllowedPathModeWrite,
		},
	}

	return
}
