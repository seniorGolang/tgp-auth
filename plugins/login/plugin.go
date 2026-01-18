package main

import (
	"fmt"
	"log/slog"

	"tgp/core/data"
	"tgp/core/i18n"
	"tgp/core/plugin"
	"tgp/plugins/login/internal/config"
	"tgp/plugins/login/internal/handler"
)

// LoginPlugin реализует интерфейс Plugin.
type LoginPlugin struct{}

// Execute выполняет основную логику плагина.
func (p *LoginPlugin) Execute(rootDir string, request data.Storage, path ...string) (response data.Storage, err error) {

	response = data.NewStorage()

	if len(path) > 0 {
		switch path[0] {
		case "login":
			if response, err = handler.HandleLogin(rootDir, request); err != nil {
				slog.Error(i18n.Msg("Execute: HandleLogin returned"), slog.Any("error", err), slog.Bool("hasResponse", response != nil))
			}
			return
		case "logout":
			return handler.HandleLogout(rootDir, request)
		default:
			err = fmt.Errorf("%s: %s", i18n.Msg("unknown command"), path[0])
			return
		}
	}

	return
}

// Info возвращает информацию о плагине.
func (p *LoginPlugin) Info() (info plugin.Info, err error) {

	info = plugin.Info{
		Name:        "login",
		Description: i18n.Msg("GitHub authentication commands"),
		Author:      "AlexK <seniorGolang@gmail.com>",
		License:     "MIT",
		Category:    "utility",
		Commands: []plugin.Command{
			{
				Path:        []string{"login"},
				Description: i18n.Msg("Authenticate with GitHub using OAuth App flow"),
			},
			{
				Path:        []string{"logout"},
				Description: i18n.Msg("Logout from GitHub (remove saved token)"),
			},
		},
		AllowedHosts: []string{
			config.AllowedHostGitHub,
			config.AllowedHostGitHubAPI,
			config.AllowedHostLocalhost,
		},
		AllowedShellCMDs: []string{
			config.CommandOpen,
			config.CommandXdgOpen,
			config.CommandCmd,
			config.CommandUname,
		},
		AllowedEnvVars: []string{
			config.EnvVarOSType,
		},
		AllowedPaths: map[string]string{
			config.AllowedPathGitHub: config.AllowedPathModeWrite,
		},
	}

	return
}
