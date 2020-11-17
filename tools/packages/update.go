package packages

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"pault.ag/go/debian/control"
	"pault.ag/go/debian/version"
)

// https://stackoverflow.com/a/33853856/5441396
func downloadFile(filepath string, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Download from %s failed with statuscode %d", url, resp.StatusCode)
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (bundle PackageBundle) UpdateFromPackageIndex(repository string, index []control.BinaryIndex) {
	for _, pkgIndex := range index {
		for idx := range bundle {
			info := &bundle[idx]
			if pkgIndex.Package == info.Name {
				if version.Compare(info.Version, pkgIndex.Version) <= 0 {
					info.Version = pkgIndex.Version
					info.Repository = repository
				}
				if info.Version == pkgIndex.Version && info.Repository == repository {
					info.Filename = pkgIndex.Filename
					info.Sha256 = pkgIndex.SHA256
				}
			}
		}
	}
}
