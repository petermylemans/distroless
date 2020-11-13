package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/distroless/tools/packages"
)

func main() {
	updatePackageBundle("amd64", "9", "stretch")
	updatePackageBundle("arm64", "9", "stretch")
	updatePackageBundle("ppc64le", "9", "stretch")
	updatePackageBundle("s390x", "9", "stretch")
	updatePackageBundle("amd64", "10", "buster")
	updatePackageBundle("arm64", "10", "buster")
	updatePackageBundle("ppc64le", "10", "buster")
	updatePackageBundle("s390x", "10", "buster")
}

func updatePackageBundle(arch string, version string, distro string) {
	debArch := arch
	if arch == "ppc64le" {
		debArch = "ppc64el"
	}
	content, err := ioutil.ReadFile("../package_bundle_" + arch + "_debian" + version + ".bzl")
	if err != nil {
		log.Fatal(err)
	}

	bundle, err := packages.ParseBundle("package_bundle_"+arch+"_debian"+version+".bzl", content)
	if err != nil {
		log.Fatal(err)
	}

	err = bundle.UpdateFromPackageIndex("debian", "http://deb.debian.org/debian", distro, "main", debArch)
	if err != nil {
		log.Fatal(err)
	}
	err = bundle.UpdateFromPackageIndex("debian", "http://deb.debian.org/debian", distro+"-updates", "main", debArch)
	if err != nil {
		log.Fatal(err)
	}
	if distro == "stretch" {
		err = bundle.UpdateFromPackageIndex("debian", "http://deb.debian.org/debian", distro+"-backports", "main", debArch)
		if err != nil {
			log.Fatal(err)
		}
	}

	if distro == "buster" || arch == "amd64" || arch == "arm64" {
		err = bundle.UpdateFromPackageIndex("debian-security", "http://deb.debian.org/debian-security", distro+"/updates", "main", debArch)
		if err != nil {
			log.Fatal(err)
		}
	}

	newFile, err := os.Create("../package_bundle_" + arch + "_debian" + version + ".bzl")
	if err != nil {
		log.Fatal(err)
	}

	err = bundle.Write(newFile)
	if err != nil {
		log.Fatal(err)
	}
}
