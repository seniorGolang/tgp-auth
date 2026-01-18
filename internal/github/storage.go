package github

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"tgp/core/i18n"
)

// LoadTokens загружает токены из файла github.yml.
func LoadTokens() (accessToken string, refreshToken string, err error) {

	tokenFile := filepath.Join(ConfigDir, TokenFileName)

	var data []byte
	if data, err = os.ReadFile(tokenFile); err != nil {
		if os.IsNotExist(err) {
			return "", "", nil
		}
		return "", "", fmt.Errorf("%s: %w", i18n.Msg("failed to read token file"), err)
	}

	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return "", "", fmt.Errorf("%s: %w", i18n.Msg("failed to parse token file"), err)
	}

	return cfg.AccessToken, cfg.RefreshToken, nil
}

// SaveToken сохраняет токены в файл github.yml.
func SaveToken(accessToken string, refreshToken string) (err error) {

	if err = os.MkdirAll(ConfigDir, DirMode); err != nil {
		return fmt.Errorf("%s: %w", i18n.Msg("failed to create config directory"), err)
	}

	cfg := Config{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	var data []byte
	if data, err = yaml.Marshal(&cfg); err != nil {
		return fmt.Errorf("%s: %w", i18n.Msg("failed to marshal token config"), err)
	}

	tokenFile := filepath.Join(ConfigDir, TokenFileName)
	if err = os.WriteFile(tokenFile, data, ConfigFileMode); err != nil {
		return fmt.Errorf("%s: %w", i18n.Msg("failed to write token file"), err)
	}

	return nil
}

// RemoveToken удаляет файл github.yml.
func RemoveToken() (err error) {

	tokenFile := filepath.Join(ConfigDir, TokenFileName)

	if err = os.Remove(tokenFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("%s: %w", i18n.Msg("failed to remove token file"), err)
	}

	return nil
}
