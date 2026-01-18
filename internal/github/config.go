package github

// ClientID - OAuth App Client ID для GitHub.
// Устанавливается через ldflags при сборке.
var ClientID string

const (
	// GitHubOAuthTokenURL - URL для получения токена через OAuth.
	//nolint:gosec // Это публичный URL, не credentials
	GitHubOAuthTokenURL = "https://github.com/login/oauth/access_token"

	// ConfigDir - путь к директории конфигурации GitHub.
	ConfigDir = "/tg/github"

	// TokenFileName - имя файла с токенами.
	//nolint:gosec // Это имя файла, не credentials
	TokenFileName = "github.yml"

	// FormFieldClientID - имя поля формы для client_id.
	FormFieldClientID = "client_id"

	// FormFieldGrantType - имя поля формы для grant_type.
	FormFieldGrantType = "grant_type"

	// FormFieldRefreshToken - имя поля формы для refresh_token.
	FormFieldRefreshToken = "refresh_token"

	// GrantTypeRefreshToken - тип grant для обновления токена.
	GrantTypeRefreshToken = "refresh_token"

	// HeaderContentType - заголовок Content-Type.
	HeaderContentType = "Content-Type"

	// HeaderAccept - заголовок Accept.
	HeaderAccept = "Accept"

	// FormURLEncodedContentType - тип контента для form-urlencoded.
	FormURLEncodedContentType = "application/x-www-form-urlencoded"

	// AllowedHostGitHubAPI - разрешенный хост для GitHub API.
	AllowedHostGitHubAPI = "api.github.com"

	// AllowedPathGitHub - разрешенный путь для GitHub токенов.
	AllowedPathGitHub = "@tg/github"

	// AllowedPathModeWrite - режим доступа для записи.
	AllowedPathModeWrite = "w"

	// ConfigFileMode - права доступа для файла конфигурации.
	ConfigFileMode = 0600

	// DirMode - права доступа для директории.
	DirMode = 0755
)
