package packages

import (
	"bufio"
	"io"
	"pault.ag/go/debian/control"
	"pault.ag/go/debian/version"
)

func (bundle PackageBundle) UpdateFromPackageIndex(indexFile io.Reader, arch string, repository string) error {
	index, err := control.ParseBinaryIndex(bufio.NewReader(indexFile))
	if err != nil {
		return err
	}

	for _, pkgIndex := range index {
		for idx := range bundle {
			info := &bundle[idx]
			if pkgIndex.Package == info.Name {
				if version.Compare(info.Version, pkgIndex.Version) < 0 {
					info.Files = nil
					info.Version = pkgIndex.Version
					info.Repository = repository
				}
				if info.Version == pkgIndex.Version {
					if info.Files == nil {
						files := make(Files)
						info.Files = files
					}
					info.Files[arch] = &File{
						Filename: pkgIndex.Filename,
						Sha256:   pkgIndex.SHA256,
					}
				}
			}
		}
	}

	return err
}
