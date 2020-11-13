package packages

import (
	"github.com/bazelbuild/buildtools/build"
	"io"
	"sort"
)

func (bundle *PackageBundle) Write(writer io.Writer) error {
	file := build.File{
		Type: build.TypeBzl,
		Stmt: []build.Expr{
			&build.StringExpr{Value: "\nPackage definition manifest extracted from debian release file.\n\nGenerated file, do not update manually.\n", TripleQuote: true},
		},
	}

	dict := build.DictExpr{ForceMultiLine: true}
	for _, packageInformation := range *bundle {
		dict.List = append(dict.List, bundle.convertPackageInformation(packageInformation))
	}
	sort.Slice(dict.List, func(i, j int) bool {
		return dict.List[i].Key.(*build.StringExpr).Value < dict.List[j].Key.(*build.StringExpr).Value
	})

	assign := build.AssignExpr{
		LHS: &build.Ident{Name: "PACKAGES"},
		Op:  "=",
		RHS: &dict,
	}

	file.Stmt = append(file.Stmt, &assign)

	_, err := writer.Write(build.Format(&file))
	return err
}

func (bundle *PackageBundle) convertPackageInformation(info PackageInformation) *build.KeyValueExpr {
	version := build.KeyValueExpr{
		Key:   &build.Ident{Name: "version"},
		Value: &build.StringExpr{Value: info.Version.String()},
	}
	repository := build.KeyValueExpr{
		Key:   &build.Ident{Name: "repository"},
		Value: &build.StringExpr{Value: info.Repository},
	}
	fileDict := build.DictExpr{ForceMultiLine: true}
	if info.Files != nil {
		for arch, file := range info.Files {
			filename := build.KeyValueExpr{
				Key:   &build.Ident{Name: "filename"},
				Value: &build.StringExpr{Value: file.Filename},
			}
			sha256 := build.KeyValueExpr{
				Key:   &build.Ident{Name: "sha256"},
				Value: &build.StringExpr{Value: file.Sha256},
			}
			fileDict.List = append(fileDict.List, &build.KeyValueExpr{
				Key:   &build.Ident{Name: arch},
				Value: &build.DictExpr{List: []*build.KeyValueExpr{&filename, &sha256}, ForceMultiLine: true},
			})
		}
	}
	files := build.KeyValueExpr{
		Key:   &build.Ident{Name: "files"},
		Value: &fileDict,
	}

	return &build.KeyValueExpr{
		Key: &build.StringExpr{Value: info.Name},
		Value: &build.DictExpr{
			List: []*build.KeyValueExpr{
				&version,
				&repository,
				&files,
			},
			ForceMultiLine: true,
		},
	}
}
