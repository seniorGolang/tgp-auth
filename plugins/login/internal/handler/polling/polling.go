package polling

import (
	"log/slog"
	"time"

	"tgp/core/i18n"
	"tgp/internal/github"
	"tgp/plugins/login/internal/oauth"
)

// CreateHandler создает обработчик задачи для polling токена авторизации.
func CreateHandler(deviceCode string, expiresAt time.Time) (handler func() (next bool)) {

	return func() (next bool) {

		defer func() {
			if r := recover(); r != nil {
				slog.Error(i18n.Msg("panic in polling handler"), slog.Any("error", r))
			}
		}()

		var tokens *github.TokenPair
		var shouldStop bool
		var pollErr error
		if tokens, shouldStop, pollErr = oauth.PollOnce(deviceCode, expiresAt); pollErr != nil {
			slog.Error(i18n.Msg("polling error"), slog.Any("error", pollErr))
			return false
		}

		if shouldStop {
			if tokens != nil && tokens.AccessToken != "" {
				slog.Info(i18n.Msg("tokens received, saving to storage"))
				var saveErr error
				if saveErr = github.SaveToken(tokens.AccessToken, tokens.RefreshToken); saveErr != nil {
					slog.Error(i18n.Msg("failed to save tokens"), slog.Any("error", saveErr))
					return false
				}
				slog.Info(i18n.Msg("login completed successfully"))
				return false
			}
			slog.Warn(i18n.Msg("polling stopped without tokens"))
			return false
		}

		return true
	}
}
