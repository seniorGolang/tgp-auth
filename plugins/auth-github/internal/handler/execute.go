package handler

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	"tgp/core/data"
	"tgp/core/i18n"
	"tgp/internal/github"
	"tgp/plugins/auth-github/internal/api"
	"tgp/plugins/auth-github/internal/config"
	"tgp/plugins/auth-github/internal/storage"
)

// HandleExecute обрабатывает выполнение плагина auth-github.
func HandleExecute(request data.Storage, path ...string) (response data.Storage, err error) {

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s: %v", i18n.Msg("panic in auth plugin"), r)
			slog.Error(i18n.Msg("panic recovered in auth plugin"),
				slog.Any("panic", r),
			)
			if response == nil {
				response = data.NewStorage()
			}
		}
	}()

	response = copyRequestToResponse(request)

	// Пропускаем проверку токена для команды login
	if len(path) > 0 && path[0] == config.CommandLogin {
		return response, nil
	}

	instanceID, err := ensureInstanceID()
	if err != nil {
		slog.Error(i18n.Msg("failed to ensure instance ID"), slog.Any("error", err))
	}

	if err = response.Set(config.StorageKeyInstanceID, instanceID); err != nil {
		slog.Error(i18n.Msg("failed to set instanceID in response"), slog.Any("error", err))
	}

	account, err := authenticateAndGetAccount()
	if err != nil {
		return response, err
	}

	if err = response.Set(config.StorageKeyAccount, account); err != nil {
		slog.Error(i18n.Msg("failed to set account in response"), slog.Any("error", err))
	}

	return
}

// copyRequestToResponse копирует данные из request в response.
func copyRequestToResponse(request data.Storage) (response data.Storage) {

	response = data.NewStorage()
	if request == nil {
		return response
	}

	storageMap, ok := request.(*data.MapStorage)
	if !ok {
		return response
	}

	for k, raw := range *storageMap {
		if err := response.Set(k, raw); err != nil {
			slog.Warn(i18n.Msg("failed to copy request data to response"), slog.String("key", k), slog.Any("error", err))
		}
	}

	return response
}

// ensureInstanceID обеспечивает наличие instance ID, загружая его или создавая новый.
func ensureInstanceID() (instanceID string, err error) {

	slog.Debug(i18n.Msg("auth directory ready"), slog.String("path", config.AuthDir))

	if instanceID, err = storage.LoadInstanceID(); err != nil {
		slog.Warn(i18n.Msg("failed to load instance ID"), slog.Any("error", err))
	}

	if instanceID == "" {
		instanceID = uuid.New().String()
		if err = storage.SaveInstanceID(instanceID); err != nil {
			slog.Error(i18n.Msg("failed to save instance ID"), slog.Any("error", err))
			return "", err
		}
		slog.Info(i18n.Msg("generated new instance ID"), slog.String("instanceID", instanceID))
	}

	return instanceID, nil
}

// authenticateAndGetAccount выполняет аутентификацию и получает данные аккаунта.
func authenticateAndGetAccount() (account *api.Account, err error) {

	accessToken, refreshToken, err := github.LoadTokens()
	if err != nil {
		slog.Warn(i18n.Msg("failed to load tokens"), slog.Any("error", err))
	}

	if accessToken == "" {
		return nil, fmt.Errorf("%s", i18n.Msg("login required via tg login"))
	}

	account, err = api.FetchUserInfo(accessToken)
	if err == nil {
		return account, nil
	}

	errStr := err.Error()
	if !strings.Contains(errStr, config.TokenInvalidError) {
		return nil, fmt.Errorf("%s: %w", i18n.Msg("failed to execute request"), err)
	}

	// Токен невалиден, пытаемся обновить
	if refreshToken == "" {
		var removeErr error
		if removeErr = github.RemoveToken(); removeErr != nil {
			slog.Debug(i18n.Msg("failed to remove invalid token file"), slog.Any("error", removeErr))
		}
		return nil, fmt.Errorf("%s: %s", i18n.Msg("login required via tg login"), errStr)
	}

	return refreshTokenAndGetAccount(refreshToken)
}

// refreshTokenAndGetAccount обновляет токен и получает данные аккаунта.
func refreshTokenAndGetAccount(refreshToken string) (account *api.Account, err error) {

	slog.Debug(i18n.Msg("attempting to refresh token"))

	newTokens, err := github.RefreshToken(refreshToken)
	if err != nil {
		var removeErr error
		if removeErr = github.RemoveToken(); removeErr != nil {
			slog.Debug(i18n.Msg("failed to remove invalid token file"), slog.Any("error", removeErr))
		}
		return nil, fmt.Errorf("%s: %w", i18n.Msg("login required via tg login"), err)
	}

	if err = github.SaveToken(newTokens.AccessToken, newTokens.RefreshToken); err != nil {
		slog.Error(i18n.Msg("failed to save refreshed tokens"), slog.Any("error", err))
	}

	account, err = api.FetchUserInfo(newTokens.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", i18n.Msg("failed to fetch account data after refresh"), err)
	}

	return account, nil
}
