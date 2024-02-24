package server

import "github.com/lucafroeschke/go-package-server/config"

type PackageResponse struct {
	Config  config.Config
	Package config.Package
}
