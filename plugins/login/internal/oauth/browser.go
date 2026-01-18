package oauth

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"tgp/core/exec"
	"tgp/core/i18n"
	"tgp/plugins/login/internal/config"
)

// DetectHostOS определяет реальную ОС хоста через команду uname или переменную окружения OSTYPE.
func DetectHostOS() (hostOS string) {

	cmd := exec.Command(config.CommandUname, "-s")
	if err := cmd.Start(); err == nil {
		var stdoutPipe io.ReadCloser
		if stdoutPipe, err = cmd.StdoutPipe(); err == nil {
			var stdoutBytes []byte
			var readErr error
			if stdoutBytes, readErr = io.ReadAll(stdoutPipe); readErr == nil {
				stdoutPipe.Close()

				var waitErr error
				if waitErr = cmd.Wait(); waitErr == nil {
					unameOutput := strings.TrimSpace(string(stdoutBytes))

					unameLower := strings.ToLower(unameOutput)
					if strings.Contains(unameLower, config.OSDarwin) {
						return config.OSDarwin
					}
					if strings.Contains(unameLower, config.OSLinux) {
						return config.OSLinux
					}
					return unameOutput
				}
			}
		}
	}

	var ostype string
	if ostype = os.Getenv(config.EnvVarOSType); ostype != "" {
		ostypeLower := strings.ToLower(ostype)
		if strings.Contains(ostypeLower, config.OSDarwin) {
			return config.OSDarwin
		}
		if strings.Contains(ostypeLower, config.OSLinux) {
			return config.OSLinux
		}
	}

	return ""
}

// OpenBrowser открывает URL в браузере пользователя.
func OpenBrowser(url string) (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s: %v", i18n.Msg("panic in openBrowser"), r)
			slog.Error(i18n.Msg("panic recovered in openBrowser"),
				slog.String("url", url),
				slog.Any("panic", r),
			)
		}
	}()

	hostOS := DetectHostOS()

	var cmd *exec.Cmd
	switch hostOS {
	case config.OSDarwin:
		cmd = exec.Command(config.CommandOpen, url)
	case config.OSLinux:
		cmd = exec.Command(config.CommandXdgOpen, url)
	default:
		commands := []struct {
			name string
			args []string
		}{
			{config.CommandOpen, []string{url}},
			{config.CommandXdgOpen, []string{url}},
			{config.CommandCmd, []string{"/c", "start", url}},
		}

		var lastErr error
		for _, cmdInfo := range commands {
			tryCmd := exec.Command(cmdInfo.name, cmdInfo.args...)

			if startErr := tryCmd.Start(); startErr != nil {
				lastErr = startErr
				continue
			}

			if waitErr := tryCmd.Wait(); waitErr != nil {
				lastErr = waitErr
				continue
			}

			return nil
		}

		slog.Error(i18n.Msg("failed to open browser"),
			slog.String("url", url),
			slog.Any("error", lastErr),
		)
		if lastErr != nil {
			return fmt.Errorf("%s: %w", i18n.Msg("failed to open browser with any command"), lastErr)
		}
		return fmt.Errorf("%s", i18n.Msg("failed to open browser: no commands available"))
	}

	if err = cmd.Start(); err != nil {
		slog.Error(i18n.Msg("failed to open browser"),
			slog.String("url", url),
			slog.Any("error", err),
		)
		return fmt.Errorf("%s: %w", i18n.Msg("failed to start browser command"), err)
	}

	if err = cmd.Wait(); err != nil {
		slog.Error(i18n.Msg("failed to open browser"),
			slog.String("url", url),
			slog.Any("error", err),
		)
		return fmt.Errorf("%s: %w", i18n.Msg("browser command failed"), err)
	}

	return nil
}
