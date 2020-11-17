package packages

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/crypto/openpgp"
	"pault.ag/go/debian/control"
	"pault.ag/go/debian/version"
)

type ReleaseIndex struct {
	control.Paragraph

	RepositoryURL   string
	DistributionURL string

	Description string
	Origin      string
	Label       string
	Version     version.Version
	Suite       string
	Codename    string
	Components  []string                 `delim:" "`
	MD5Sum      []control.MD5FileHash    `delim:"\n" strip:"\n\r\t "`
	SHA256      []control.SHA256FileHash `delim:"\n" strip:"\n\r\t "`
}

func DownloadRelease(repository string, dist string, keyring *openpgp.EntityList) (*ReleaseIndex, error) {
	index := ReleaseIndex{}
	index.RepositoryURL = repository
	index.DistributionURL = repository + "/dists/" + dist

	resp, err := http.Get(repository + "/dists/" + dist + "/InRelease")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		decoder, _ := control.NewDecoder(resp.Body, keyring)
		if err != nil {
			return nil, err
		}

		return &index, decoder.Decode(&index)
	} else if resp.StatusCode == 404 {
		resp, err = http.Get(repository + "/dists/" + dist + "/Release")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		decoder, _ := control.NewDecoder(resp.Body, nil)
		if err != nil {
			return nil, err
		}

		return &index, decoder.Decode(&index)
	}

	return nil, fmt.Errorf("Got an unexpected status code %d for url %s", resp.StatusCode, resp.Request.URL)
}

func (index *ReleaseIndex) GetBinaryIndex(component string, arch string) ([]control.BinaryIndex, error) {
	targetFilename := fmt.Sprintf("%s/binary-%s/Packages.gz", component, arch)
	var filehash *control.FileHash
	for _, file := range index.SHA256 {
		if file.Filename == targetFilename {
			filehash = &file.FileHash
			break
		}
	}

	if filehash == nil {
		return nil, fmt.Errorf("Could not find [%s] in index %s", targetFilename, index.DistributionURL)
	}

	log.Print("Fetching " + index.DistributionURL + "/" + targetFilename)
	resp, err := http.Get(index.DistributionURL + "/" + targetFilename)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	verifier, err := filehash.Verifier()
	if err != nil {
		return nil, err
	}

	reader := io.TeeReader(resp.Body, verifier)
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	binaryIndex, err := control.ParseBinaryIndex(bufio.NewReader(gzipReader))
	if err != nil {
		return nil, err
	}

	err = verifier.Close()
	if err != nil {
		return nil, err
	}

	return binaryIndex, nil
}
