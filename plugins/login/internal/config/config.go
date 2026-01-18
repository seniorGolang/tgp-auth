package config

const (
	// GitHubDeviceCodeURL - URL для запроса device code.
	GitHubDeviceCodeURL = "https://github.com/login/device/code"

	// FormFieldDeviceCode - имя поля формы для device_code.
	FormFieldDeviceCode = "device_code"

	// GrantTypeDeviceCode - тип grant для device code flow.
	GrantTypeDeviceCode = "urn:ietf:params:oauth:grant-type:device_code"

	// OAuthScope - область доступа OAuth.
	OAuthScope = "read:user"

	// ErrorAuthorizationPending - ошибка "authorization_pending" от GitHub.
	ErrorAuthorizationPending = "authorization_pending"

	// ErrorSlowDown - ошибка "slow_down" от GitHub.
	ErrorSlowDown = "slow_down"

	// JSONAccept - значение Accept заголовка для JSON.
	JSONAccept = "application/json"

	// DefaultPollInterval - интервал polling по умолчанию в секундах.
	DefaultPollInterval = 5

	// AllowedHostGitHub - разрешенный хост для GitHub.
	AllowedHostGitHub = "github.com"

	// AllowedHostGitHubAPI - разрешенный хост для GitHub API.
	AllowedHostGitHubAPI = "api.github.com"

	// AllowedHostLocalhost - разрешенный хост для localhost.
	AllowedHostLocalhost = "127.0.0.1"

	// AllowedPathGitHub - разрешенный путь для GitHub токенов.
	AllowedPathGitHub = "@tg/github"

	// AllowedPathModeWrite - режим доступа для записи.
	AllowedPathModeWrite = "w"

	// EnvVarOSType - переменная окружения для определения ОС.
	EnvVarOSType = "OSTYPE"

	// CommandOpen - команда для открытия браузера на macOS.
	CommandOpen = "open"

	// CommandXdgOpen - команда для открытия браузера на Linux.
	CommandXdgOpen = "xdg-open"

	// CommandCmd - команда для открытия браузера на Windows.
	CommandCmd = "cmd"

	// CommandUname - команда для определения ОС.
	CommandUname = "uname"

	// OSDarwin - название ОС Darwin (macOS).
	OSDarwin = "darwin"

	// OSLinux - название ОС Linux.
	OSLinux = "linux"

	// ResponseKeyAuthURL - ключ для authURL в response.
	ResponseKeyAuthURL = "authURL"

	// ResponseKeyBrowserOpened - ключ для browserOpened в response.
	ResponseKeyBrowserOpened = "browserOpened"

	// ResponseKeyBrowserError - ключ для browserError в response.
	ResponseKeyBrowserError = "browserError"

	// ResponseKeyMessage - ключ для message в response.
	ResponseKeyMessage = "message"

	// ResponseKeySuccess - ключ для success в response.
	ResponseKeySuccess = "success"

	// FormFieldScope - имя поля формы для scope.
	FormFieldScope = "scope"

	// URLQueryParamUserCode - параметр запроса user_code для URL.
	URLQueryParamUserCode = "user_code"
)
