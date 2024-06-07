package cmd

import "github.com/lucafroeschke/go-package-server/config"

func RunConfigCommand() {
	err := config.CreateConfig()
	if err != nil {
		return
	}
}
