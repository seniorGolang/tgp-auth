package handler

import (
	"fmt"
	"log/slog"
	"time"

	"tgp/core"
	"tgp/core/data"
	"tgp/core/i18n"
	"tgp/plugins/login/internal/config"
	"tgp/plugins/login/internal/handler/polling"
	"tgp/plugins/login/internal/oauth"
)

// HandleLogin обрабатывает команду login.
func HandleLogin(_ string, _ data.Storage) (response data.Storage, err error) {

	slog.Debug(i18n.Msg("login command started"))

	response = data.NewStorage()

	var deviceResp *oauth.DeviceCodeResponse
	if deviceResp, err = oauth.RequestDeviceCode(); err != nil {
		err = fmt.Errorf("%s: %w", i18n.Msg("failed to request device code"), err)
		return
	}

	authURL := oauth.BuildDeviceAuthURL(deviceResp.VerificationURI, deviceResp.UserCode)

	_ = response.Set(config.ResponseKeyAuthURL, authURL)

	slog.Info(fmt.Sprintf(i18n.Msg("authorization code: %s"), deviceResp.UserCode))

	// Интерактивный выбор: открыть браузер или показать адрес
	openBrowserOption := i18n.Msg("open browser")
	showURLOption := i18n.Msg("show url")
	options := []string{
		openBrowserOption,
		showURLOption,
	}
	var selected []string
	var selectErr error
	if selected, selectErr = core.InteractiveSelect(i18n.Msg("open browser question"), options, false, nil); selectErr != nil {
		slog.Warn(i18n.Msg("interactive select failed"), slog.String("error", selectErr.Error()))
		// При ошибке интерактивного выбора показываем адрес
		slog.Info(fmt.Sprintf(i18n.Msg("please open url: %s"), authURL))
	} else if len(selected) > 0 {
		selectedOption := selected[0]
		if selectedOption == openBrowserOption {
			// Пользователь выбрал открыть браузер
			var browserErr error
			if browserErr = oauth.OpenBrowser(authURL); browserErr == nil {
				slog.Info(i18n.Msg("browser opened successfully"))
				_ = response.Set(config.ResponseKeyBrowserOpened, true)
			} else {
				slog.Error(i18n.Msg("failed to open browser"),
					slog.String("url", authURL),
					slog.String("error", browserErr.Error()),
				)
				_ = response.Set(config.ResponseKeyBrowserOpened, false)
				_ = response.Set(config.ResponseKeyBrowserError, browserErr.Error())
				slog.Info(fmt.Sprintf(i18n.Msg("please open url: %s"), authURL))
			}
		} else {
			// Пользователь выбрал показать адрес
			slog.Info(fmt.Sprintf(i18n.Msg("please open url: %s"), authURL))
			_ = response.Set(config.ResponseKeyBrowserOpened, false)
		}
	}

	slog.Debug(i18n.Msg("waiting for authorization"))

	var pollInterval time.Duration
	if deviceResp.Interval > 0 {
		pollInterval = time.Duration(deviceResp.Interval) * time.Second
	} else {
		pollInterval = time.Duration(config.DefaultPollInterval) * time.Second
	}
	expiresAt := time.Now().Add(time.Duration(deviceResp.ExpiresIn) * time.Second)

	slog.Debug(i18n.Msg("starting polling task"), slog.String("interval", pollInterval.String()))
	handler := polling.CreateHandler(deviceResp.DeviceCode, expiresAt)
	if _, err = core.StartTask(pollInterval, handler); err != nil {
		slog.Error(i18n.Msg("failed to start polling task"), slog.Any("error", err))
		err = fmt.Errorf("%s: %w", i18n.Msg("failed to start polling task"), err)
		return
	}
	slog.Debug(i18n.Msg("polling task started successfully"))

	return
}
