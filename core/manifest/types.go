package manifest

// Manifest представляет упрощенную версию plugin.Info для публикации в каталог.
type Manifest struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Author       string   `json:"author"`
	License      string   `json:"license"`
	Dependencies []string `json:"dependencies,omitempty"`
}
