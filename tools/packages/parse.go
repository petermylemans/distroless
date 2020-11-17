package packages

import (
	"github.com/bazelbuild/buildtools/build"
	"pault.ag/go/debian/version"
)

func ParseBundle(filename string, data []byte) (*PackageBundle, error) {
	f, err := build.ParseBzl(filename, data)
	if err != nil {
		return nil, err
	}

	var bundle PackageBundle
	for _, value := range f.Stmt {
		assignExpr, ok := value.(*build.AssignExpr)
		if ok {
			dict, ok := assignExpr.RHS.(*build.DictExpr)
			if ok {
				for _, keyValue := range dict.List {
					bundle = append(bundle, *parsePackageInformation(keyValue))
				}
			}
		}
	}

	return &bundle, err
}

func parsePackageInformation(keyValue *build.KeyValueExpr) *PackageInformation {
	pkg := PackageInformation{
		Name: keyValue.Key.(*build.StringExpr).Value,
	}

	attrs, ok := keyValue.Value.(*build.DictExpr)
	if ok {
		for _, attrValue := range attrs.List {
			if attrValue.Key.(*build.StringExpr).Value == "version" {
				pkg.Version, _ = version.Parse(attrValue.Value.(*build.StringExpr).Value)
			} else if attrValue.Key.(*build.StringExpr).Value == "repository" {
				pkg.Repository = attrValue.Value.(*build.StringExpr).Value
			} else if attrValue.Key.(*build.StringExpr).Value == "filename" {
				pkg.Filename = attrValue.Value.(*build.StringExpr).Value
			} else if attrValue.Key.(*build.StringExpr).Value == "sha256" {
				pkg.Sha256 = attrValue.Value.(*build.StringExpr).Value
			}
		}
	}

	return &pkg
}

