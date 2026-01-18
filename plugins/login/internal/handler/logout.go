package handler

import (
	"fmt"
	"log/slog"

	"tgp/core/data"
	"tgp/core/i18n"
	"tgp/internal/github"
	"tgp/plugins/login/internal/config"
)

// HandleLogout обрабатывает команду logout.
func HandleLogout(rootDir string, request data.Storage) (response data.Storage, err error) {

	slog.Info(i18n.Msg("logout command started"))

	response = data.NewStorage()

	if err = github.RemoveToken(); err != nil {
		err = fmt.Errorf("%s: %w", i18n.Msg("failed to remove token"), err)
		return
	}

	slog.Info(i18n.Msg("logout completed, token removed"))

	_ = response.Set(config.ResponseKeySuccess, true)

	return
}
