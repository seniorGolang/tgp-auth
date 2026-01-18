package config

const (
	// GitHubAPIUserURL - URL для получения информации о пользователе GitHub.
	GitHubAPIUserURL = "https://api.github.com/user"

	// AuthDir - путь к директории auth для хранения instance.id.
	AuthDir = "/tg/auth"

	// InstanceIDFileName - имя файла с instance ID.
	InstanceIDFileName = "instance.id"

	// TokenInvalidError - текст ошибки о невалидном токене.
	TokenInvalidError = "token is invalid or revoked"

	// StorageKeyInstanceID - ключ для хранения instance ID в storage.
	StorageKeyInstanceID = "instanceID"

	// StorageKeyAccount - ключ для хранения данных аккаунта в storage.
	StorageKeyAccount = "account"

	// AllowedPathAuth - разрешенный путь для auth данных.
	AllowedPathAuth = "@tg/auth"

	// HeaderAuthorization - заголовок Authorization.
	HeaderAuthorization = "Authorization"

	// BearerPrefix - префикс для Bearer токена.
	BearerPrefix = "Bearer "

	// GitHubAPIAccept - значение Accept заголовка для GitHub API v3.
	GitHubAPIAccept = "application/vnd.github.v3+json"

	// CommandLogin - команда login для пропуска проверки токена.
	CommandLogin = "login"
)
