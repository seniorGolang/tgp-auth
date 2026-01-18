package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"tgp/core/i18n"
	"tgp/internal/github"
	"tgp/plugins/auth-github/internal/config"
)

// LoadInstanceID загружает instance ID из файла instance.id.
func LoadInstanceID() (instanceID string, err error) {

	instanceFile := filepath.Join(config.AuthDir, config.InstanceIDFileName)

	var data []byte
	if data, err = os.ReadFile(instanceFile); err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("%s: %w", i18n.Msg("failed to read instance ID file"), err)
	}

	return string(data), nil
}

// SaveInstanceID сохраняет instance ID в файл instance.id.
func SaveInstanceID(instanceID string) (err error) {

	if err = os.MkdirAll(config.AuthDir, github.DirMode); err != nil {
		return fmt.Errorf("%s: %w", i18n.Msg("failed to create auth directory"), err)
	}

	instanceFile := filepath.Join(config.AuthDir, config.InstanceIDFileName)
	if err = os.WriteFile(instanceFile, []byte(instanceID), github.ConfigFileMode); err != nil {
		return fmt.Errorf("%s: %w", i18n.Msg("failed to write instance ID file"), err)
	}

	return nil
}
