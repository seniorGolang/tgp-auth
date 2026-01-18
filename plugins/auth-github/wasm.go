package main

import (
	"tgp/core"
)

func init() {

	core.InitPlugin(&AuthGithubPlugin{})
}

func main() {

	// Инициализация не требуется для wasip1
}
