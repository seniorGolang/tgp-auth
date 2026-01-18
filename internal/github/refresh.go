package github

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"strings"

	"tgp/core/http"
	"tgp/core/i18n"
)

// RefreshToken обновляет access token используя refresh token.
func RefreshToken(refreshToken string) (tokens *TokenPair, err error) {

	if refreshToken == "" {
		return nil, fmt.Errorf("%s", i18n.Msg("refresh token is empty"))
	}

	client := NewLoggingClient(http.DefaultClient)

	reqBody := url.Values{}
	reqBody.Set(FormFieldClientID, ClientID)
	reqBody.Set(FormFieldRefreshToken, refreshToken)
	reqBody.Set(FormFieldGrantType, GrantTypeRefreshToken)

	var refreshReq *http.Request
	if refreshReq, err = http.NewRequest(http.MethodPost, GitHubOAuthTokenURL, strings.NewReader(reqBody.Encode())); err != nil {
		return nil, fmt.Errorf("%s: %w", i18n.Msg("failed to create refresh request"), err)
	}

	refreshReq.Header.Set(HeaderContentType, FormURLEncodedContentType)
	refreshReq.Header.Set(HeaderAccept, "application/json")

	var resp *http.Response
	if resp, err = client.Do(refreshReq); err != nil {
		return nil, fmt.Errorf("%s: %w", i18n.Msg("failed to execute refresh request"), err)
	}
	defer resp.Body.Close()

	var bodyBytes []byte
	if bodyBytes, err = io.ReadAll(resp.Body); err != nil {
		return nil, fmt.Errorf("%s: %w", i18n.Msg("failed to read response body"), err)
	}

	if resp.StatusCode != http.StatusOK {
		slog.Warn(i18n.Msg("refresh token request failed"), slog.Int("status", resp.StatusCode), slog.String("body", string(bodyBytes)))
		return nil, fmt.Errorf("%s: %d, body: %s", i18n.Msg("unexpected status code"), resp.StatusCode, string(bodyBytes))
	}

	var tokenResp RefreshTokenResponse
	if err = json.Unmarshal(bodyBytes, &tokenResp); err != nil {
		return nil, fmt.Errorf("%s: %w", i18n.Msg("failed to decode response"), err)
	}

	if tokenResp.Error != "" {
		return nil, fmt.Errorf("%s: %s - %s", i18n.Msg("refresh token error"), tokenResp.Error, tokenResp.ErrorDescription)
	}

	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("%s", i18n.Msg("no access token in refresh response"))
	}

	slog.Info(i18n.Msg("token refreshed successfully"))
	return &TokenPair{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
	}, nil
}
