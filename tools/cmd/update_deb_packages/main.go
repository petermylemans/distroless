package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/distroless/tools/packages"
	"golang.org/x/crypto/openpgp"
)

func main() {
	updateStretch()
	updateBuster()
}

func updateStretch() {
	keyring := readKeyRingFile([]string{"archive-key-9.asc", "archive-key-9-security.asc"})

	architectures := []string{"amd64", "arm64", "ppc64el", "s390x"}

	securitySource := downloadRelease("http://deb.debian.org/debian-security", "stretch/updates", keyring)
	for _, arch := range []string{"amd64", "arm64"} {
		updateFromSource("package_bundle_"+arch+"_debian9.bzl", "debian-security", arch, securitySource)
	}

	for _, dist := range []string{"stretch", "stretch-updates", "stretch-backports"} {
		source := downloadRelease("http://deb.debian.org/debian", dist, keyring)
		for _, arch := range architectures {
			updateFromSource("package_bundle_"+arch+"_debian9.bzl", "debian", arch, source)
		}
	}
}

func updateBuster() {
	keyring := readKeyRingFile([]string{"archive-key-10.asc", "archive-key-10-security.asc"})

	architectures := []string{"amd64", "arm64", "ppc64el", "s390x"}

	securitySource := downloadRelease("http://deb.debian.org/debian-security", "buster/updates", keyring)
	for _, arch := range architectures {
		updateFromSource("package_bundle_"+arch+"_debian10.bzl", "debian-security", arch, securitySource)
	}

	for _, dist := range []string{"buster", "buster-updates"} {
		source := downloadRelease("http://deb.debian.org/debian", dist, keyring)
		for _, arch := range architectures {
			updateFromSource("package_bundle_"+arch+"_debian10.bzl", "debian", arch, source)
		}
	}
}

func updateFromSource(bundleName string, repository string, arch string, source *packages.ReleaseIndex) {
	bundleName = strings.ReplaceAll(bundleName, "ppc64el", "ppc64le")
	content, err := ioutil.ReadFile("../" + bundleName)
	logFatal(err)

	bundle, err := packages.ParseBundle(bundleName, content)
	logFatal(err)

	binaryIndex, err := source.GetBinaryIndex("main", arch)
	logFatal(err)

	bundle.UpdateFromPackageIndex(repository, binaryIndex)

	newFile, err := os.Create("../" + bundleName)
	logFatal(err)
	err = bundle.Write(newFile)
	logFatal(err)
}

func readKeyRingFile(filenames []string) *openpgp.EntityList {
	var keyring openpgp.EntityList
	for _, filename := range filenames {
		f, err := os.Open(filename)
		logFatal(err)

		defer f.Close()
		entry, err := openpgp.ReadArmoredKeyRing(f)
		logFatal(err)

		keyring = append(keyring, entry...)
	}

	return &keyring
}

func downloadRelease(mirror string, distribution string, keyring *openpgp.EntityList) *packages.ReleaseIndex {
	release, err := packages.DownloadRelease(mirror, distribution, keyring)
	logFatal(err)

	return release
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
