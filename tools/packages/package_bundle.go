package packages

import "pault.ag/go/debian/version"

type File struct {
	Filename string
	Sha256   string
}

type Files map[string]*File

type PackageInformation struct {
	Name       string
	Version    version.Version
	Repository string
	Files      Files
}

type PackageBundle []PackageInformation
