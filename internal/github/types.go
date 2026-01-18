package github

// TokenPair представляет пару access и refresh токенов.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// Config представляет структуру файла github.yml.
type Config struct {
	AccessToken  string `yaml:"access_token"`
	RefreshToken string `yaml:"refresh_token,omitempty"`
}

// RefreshTokenResponse представляет ответ от GitHub при обновлении токена.
type RefreshTokenResponse struct {
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	TokenType             string `json:"token_type"`
	Scope                 string `json:"scope"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	Error                 string `json:"error"`
	ErrorDescription      string `json:"error_description"`
}
