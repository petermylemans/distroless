package packages

import (
	"github.com/bazelbuild/buildtools/build"
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

	return &pkg
}
