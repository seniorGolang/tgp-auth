package api

import (
	"encoding/json"
	"fmt"
	"io"

	"tgp/core/http"
	"tgp/core/i18n"
	"tgp/internal/github"
	"tgp/plugins/auth-github/internal/config"
)

// FetchUserInfo получает информацию о пользователе из GitHub API.
func FetchUserInfo(token string) (account *Account, err error) {

	var req *http.Request
	if req, err = http.NewRequest(http.MethodGet, config.GitHubAPIUserURL, nil); err != nil {
		return nil, fmt.Errorf("%s: %w", i18n.Msg("failed to create request"), err)
	}

	req.Header.Set(config.HeaderAuthorization, config.BearerPrefix+token)
	req.Header.Set(github.HeaderAccept, config.GitHubAPIAccept)

	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
		return nil, fmt.Errorf("%s: %w", i18n.Msg("failed to execute request"), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("%s (status: %d)", config.TokenInvalidError, resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		var body []byte
		body, _ = io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%s: %d, body: %s", i18n.Msg("unexpected status code"), resp.StatusCode, string(body))
	}

	var accountData Account
	if err = json.NewDecoder(resp.Body).Decode(&accountData); err != nil {
		return nil, fmt.Errorf("%s: %w", i18n.Msg("failed to decode response"), err)
	}

	account = &accountData
	return
}
