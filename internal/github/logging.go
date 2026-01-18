package github

import (
	"log/slog"

	"tgp/core/http"
	"tgp/core/i18n"
)

type loggingRoundTripper struct {
	rt http.RoundTripper
}

// RoundTrip выполняет запрос и логирует только ошибки.
func (l *loggingRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	resp, err = l.rt.RoundTrip(req)
	if err != nil {
		slog.Error(i18n.Msg("HTTP request failed"),
			slog.String("method", req.Method),
			slog.String("url", req.URL.String()),
			slog.Any("error", err),
		)
	}

	return
}

// NewLoggingClient создает HTTP клиент с логированием запросов и ответов.
func NewLoggingClient(baseClient *http.Client) (loggingClient *http.Client) {

	var baseTransport http.RoundTripper
	if baseClient.Transport == nil {
		baseTransport = http.DefaultTransport
	} else {
		baseTransport = baseClient.Transport
	}

	return &http.Client{
		Transport: &loggingRoundTripper{
			rt: baseTransport,
		},
		Timeout:       baseClient.Timeout,
		Jar:           baseClient.Jar,
		CheckRedirect: baseClient.CheckRedirect,
	}
}
