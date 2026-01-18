package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"tgp/core/http"
	"tgp/core/i18n"
	"tgp/internal/github"
	"tgp/plugins/login/internal/config"
)

// RequestDeviceCode запрашивает device code у GitHub.
func RequestDeviceCode() (deviceResp *DeviceCodeResponse, err error) {

	reqBody := url.Values{}
	reqBody.Set(github.FormFieldClientID, github.ClientID)
	reqBody.Set(config.FormFieldScope, config.OAuthScope)

	var req *http.Request
	if req, err = http.NewRequest(http.MethodPost, config.GitHubDeviceCodeURL, strings.NewReader(reqBody.Encode())); err != nil {
		return nil, fmt.Errorf("%s: %w", i18n.Msg("failed to create request"), err)
	}

	req.Header.Set(github.HeaderContentType, github.FormURLEncodedContentType)
	req.Header.Set(github.HeaderAccept, config.JSONAccept)

	client := github.NewLoggingClient(http.DefaultClient)

	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		return nil, fmt.Errorf("%s: %w", i18n.Msg("failed to execute request"), err)
	}
	defer resp.Body.Close()

	var bodyBytes []byte
	if bodyBytes, err = io.ReadAll(resp.Body); err != nil {
		return nil, fmt.Errorf("%s: %w", i18n.Msg("failed to read response body"), err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %d, body: %s", i18n.Msg("unexpected status code"), resp.StatusCode, string(bodyBytes))
	}

	var deviceRespData DeviceCodeResponse
	if err = json.Unmarshal(bodyBytes, &deviceRespData); err != nil {
		return nil, fmt.Errorf("%s: %w", i18n.Msg("failed to decode response"), err)
	}

	return &deviceRespData, nil
}

// PollOnce выполняет один запрос polling для получения токена.
func PollOnce(deviceCode string, expiresAt time.Time) (tokens *github.TokenPair, shouldStop bool, err error) {

	if time.Now().After(expiresAt) {
		return nil, true, fmt.Errorf("%s", i18n.Msg("device code expired"))
	}

	client := github.NewLoggingClient(http.DefaultClient)

	reqBody := url.Values{}
	reqBody.Set(github.FormFieldClientID, github.ClientID)
	reqBody.Set(config.FormFieldDeviceCode, deviceCode)
	reqBody.Set(github.FormFieldGrantType, config.GrantTypeDeviceCode)

	var pollReq *http.Request
	if pollReq, err = http.NewRequest(http.MethodPost, github.GitHubOAuthTokenURL, strings.NewReader(reqBody.Encode())); err != nil {
		return nil, true, fmt.Errorf("%s: %w", i18n.Msg("failed to create polling request"), err)
	}

	pollReq.Header.Set(github.HeaderContentType, github.FormURLEncodedContentType)
	pollReq.Header.Set(github.HeaderAccept, config.JSONAccept)

	var resp *http.Response
	if resp, err = client.Do(pollReq); err != nil {
		return nil, true, fmt.Errorf("%s: %w", i18n.Msg("failed to execute polling request"), err)
	}

	var bodyBytes []byte
	if bodyBytes, err = io.ReadAll(resp.Body); err != nil {
		resp.Body.Close()
		return nil, true, fmt.Errorf("%s: %w", i18n.Msg("failed to read response body"), err)
	}
	resp.Body.Close()

	var tokenResp DeviceTokenResponse
	if err = json.Unmarshal(bodyBytes, &tokenResp); err != nil {
		return nil, true, fmt.Errorf("%s: %w", i18n.Msg("failed to decode response"), err)
	}

	if tokenResp.Error != "" {
		if tokenResp.Error == config.ErrorAuthorizationPending || tokenResp.Error == config.ErrorSlowDown {
			return nil, false, nil
		}
		return nil, true, fmt.Errorf("%s: %s - %s", i18n.Msg("token polling error"), tokenResp.Error, tokenResp.ErrorDescription)
	}

	if tokenResp.AccessToken != "" {
		return &github.TokenPair{
			AccessToken:  tokenResp.AccessToken,
			RefreshToken: tokenResp.RefreshToken,
		}, true, nil
	}

	return nil, false, nil
}

// BuildDeviceAuthURL строит URL для авторизации устройства с user code.
func BuildDeviceAuthURL(verificationURI string, userCode string) (authURL string) {

	parsedURL, err := url.Parse(verificationURI)
	if err != nil {
		return verificationURI
	}

	query := parsedURL.Query()
	query.Set(config.URLQueryParamUserCode, userCode)
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String()
}
