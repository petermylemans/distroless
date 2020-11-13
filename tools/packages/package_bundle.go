package packages

import "pault.ag/go/debian/version"

type PackageInformation struct {
	Name       string
	Version    version.Version
	Repository string
	Filename   string
	Sha256     string
}

type PackageBundle []PackageInformation
