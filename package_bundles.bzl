load("//:package_bundle_amd64_debian9.bzl", package_bundle_amd64_debian9 = "PACKAGES")
load("//:package_bundle_amd64_debian10.bzl", package_bundle_amd64_debian10 = "PACKAGES")
load("//:package_bundle_arm64_debian9.bzl", package_bundle_arm64_debian9 = "PACKAGES")
load("//:package_bundle_arm64_debian10.bzl", package_bundle_arm64_debian10 = "PACKAGES")
load("//:package_bundle_ppc64le_debian9.bzl", package_bundle_ppc64le_debian9 = "PACKAGES")
load("//:package_bundle_ppc64le_debian10.bzl", package_bundle_ppc64le_debian10 = "PACKAGES")
load("//:package_bundle_s390x_debian9.bzl", package_bundle_s390x_debian9 = "PACKAGES")
load("//:package_bundle_s390x_debian10.bzl", package_bundle_s390x_debian10 = "PACKAGES")
load(":package_repositories.bzl", "DISTRO_REPOSITORIES")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_file")

BASE_ARCHITECTURES = ["amd64", "arm64"]
ARCHITECTURES = BASE_ARCHITECTURES + ["s390x", "ppc64le"]

DISTRO_SUFFIXES = ("_debian9", "_debian10")

DISTRO_PACKAGE_INFO = {
    "amd64": {
        "_debian9": package_bundle_amd64_debian9,
        "_debian10": package_bundle_amd64_debian10,
    },
    "arm64": {
        "_debian9": package_bundle_arm64_debian9,
        "_debian10": package_bundle_arm64_debian10,
    },
    "s390x": {
        "_debian9": package_bundle_s390x_debian9,
        "_debian10": package_bundle_s390x_debian10,
    },
    "ppc64le": {
        "_debian9": package_bundle_ppc64le_debian9,
        "_debian10": package_bundle_ppc64le_debian10,
    },
}

DISTRO_PACKAGES = {
    arch: {
        suffix: {pkg: "@" + pkg.replace("+", "-") + "_" + arch + suffix + "//file" for pkg in DISTRO_PACKAGE_INFO[arch][suffix]}
        for suffix in DISTRO_SUFFIXES
    }
    for arch in ARCHITECTURES
}

def package_http_files():
    for arch in ARCHITECTURES:
        for suffix in DISTRO_SUFFIXES:
            for pkgName in DISTRO_PACKAGE_INFO[arch][suffix]:
                pkg = DISTRO_PACKAGE_INFO[arch][suffix][pkgName]
                if pkg["repository"] != "":
                    http_file(
                        name = pkgName.replace("+", "-") + "_" + arch + suffix,
                        downloaded_file_path = "file.deb",
                        sha256 = pkg["sha256"],
                        urls = [DISTRO_REPOSITORIES[pkg["repository"]] + "/" + pkg["filename"]],
                    )
