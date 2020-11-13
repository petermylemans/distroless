package packages

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
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

func (bundle PackageBundle) UpdateFromPackageIndex(repository string, mirror string, distribution string, component string, arch string) (err error) {
	packageIndex, err := ioutil.TempFile("", "Packages*")
	packageIndexGz, err := ioutil.TempFile("", "Packages*.gz")
	if err != nil {
		return err
	}

	err = downloadFile(packageIndexGz.Name(), fmt.Sprintf("%s/dists/%s/%s/binary-%s/Packages.gz", mirror, distribution, component, arch))
	if err != nil {
		return err
	}

	zipReader, err := gzip.NewReader(packageIndexGz)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	content, err := ioutil.ReadAll(zipReader)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(packageIndex.Name(), os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(content)
	if err != nil {
		return err
	}

	index, err := control.ParseBinaryIndex(bufio.NewReader(packageIndex))
	if err != nil {
		return err
	}

	for _, pkgIndex := range index {
		for idx := range bundle {
			info := &bundle[idx]
			if pkgIndex.Package == info.Name {
				if version.Compare(info.Version, pkgIndex.Version) < 0 {
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

	return err
}
