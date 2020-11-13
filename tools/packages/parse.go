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
			if attrValue.Key.(*build.Ident).Name == "version" {
				pkg.Version, _ = version.Parse(attrValue.Value.(*build.StringExpr).Value)
			} else if attrValue.Key.(*build.Ident).Name == "repository" {
				pkg.Repository = attrValue.Value.(*build.StringExpr).Value
			} else if attrValue.Key.(*build.Ident).Name == "files" {
				pkg.Files = parseFiles(attrValue.Value.(*build.DictExpr))
			}
		}
	}

	return &pkg
}

func parseFiles(expr *build.DictExpr) Files {
	files := make(Files)

	for _, keyValueExpr := range expr.List {
		arch := keyValueExpr.Key.(*build.Ident).Name
		files[arch] = parseFile(keyValueExpr.Value.(*build.DictExpr))
	}

	return files
}

func parseFile(expr *build.DictExpr) *File {
	file := File{}

	for _, attrValue := range expr.List {
		if attrValue.Key.(*build.Ident).Name == "filename" {
			file.Filename = attrValue.Value.(*build.StringExpr).Value
		} else if attrValue.Key.(*build.Ident).Name == "sha256" {
			file.Sha256 = attrValue.Value.(*build.StringExpr).Value
		}
	}

	return &file
}
