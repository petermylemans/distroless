package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/distroless/tools/packages"
)

func main() {
	content, err := ioutil.ReadFile("../package_bundle_debian9.bzl")
	if err != nil {
		log.Fatal(err)
	}

	row, err := packages.ParseBundle("package_bundle_debian9.bzl", content)

	packagesFile, err := os.Open("../Packages")
	if err != nil {
		log.Fatal(err)
	}
	defer packagesFile.Close()

	err = row.UpdateFromPackageIndex("debian", "http://deb.debian.org/debian", "stretch", "main", "amd64")
	if err != nil {
		log.Fatal(err)
	}
	err = row.UpdateFromPackageIndex("debian", "http://deb.debian.org/debian", "stretch-updates", "main", "amd64")
	if err != nil {
		log.Fatal(err)
	}
	err = row.UpdateFromPackageIndex("debian", "http://deb.debian.org/debian", "stretch-backports", "main", "amd64")
	if err != nil {
		log.Fatal(err)
	}
	err = row.UpdateFromPackageIndex("debian-security", "http://deb.debian.org/debian-security", "stretch/updates", "main", "amd64")
	if err != nil {
		log.Fatal(err)
	}

	newFile, err := os.Create("../package_bundle_debian9.bzl")
	if err != nil {
		log.Fatal(err)
	}

	err = row.Write(newFile)
	if err != nil {
		log.Fatal(err)
	}

}
