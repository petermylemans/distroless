package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/distroless/tools/packages"
)

func main() {
	updateStretch("amd64", "9", "stretch")
	updateStretch("arm64", "9", "stretch")
	updateStretch("ppc64le", "9", "stretch")
	updateStretch("s390x", "9", "stretch")
	updateStretch("amd64", "10", "buster")
	updateStretch("arm64", "10", "buster")
	updateStretch("ppc64le", "10", "buster")
	updateStretch("s390x", "10", "buster")
}

func updateStretch(arch string, version string, distro string) {
	debArch := arch
	if arch == "ppc64le" {
		debArch = "ppc64el"
	}
	content, err := ioutil.ReadFile("../package_bundle_" + arch + "_debian" + version + ".bzl")
	if err != nil {
		log.Fatal(err)
	}

	row, err := packages.ParseBundle("package_bundle_"+arch+"_debian"+version+".bzl", content)

	packagesFile, err := os.Open("../Packages")
	if err != nil {
		log.Fatal(err)
	}
	defer packagesFile.Close()

	err = row.UpdateFromPackageIndex("debian", "http://deb.debian.org/debian", distro, "main", debArch)
	if err != nil {
		log.Fatal(err)
	}
	err = row.UpdateFromPackageIndex("debian", "http://deb.debian.org/debian", distro+"-updates", "main", debArch)
	if err != nil {
		log.Fatal(err)
	}
	if distro == "stretch" {
		err = row.UpdateFromPackageIndex("debian", "http://deb.debian.org/debian", distro+"-backports", "main", debArch)
		if err != nil {
			log.Fatal(err)
		}
	}

	if arch == "amd64" || arch == "arm64" {
		err = row.UpdateFromPackageIndex("debian-security", "http://deb.debian.org/debian-security", distro+"/updates", "main", debArch)
		if err != nil {
			log.Fatal(err)
		}
	}

	newFile, err := os.Create("../package_bundle_" + arch + "_debian" + version + ".bzl")
	if err != nil {
		log.Fatal(err)
	}

	err = row.Write(newFile)
	if err != nil {
		log.Fatal(err)
	}
}
