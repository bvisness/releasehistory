package npm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bvisness/releasehistory/src/oops"
)

var client = &http.Client{
	Timeout: 5 * time.Second,
}

type Package struct {
	Name     string             `json:"name"`
	Versions map[string]Version `json:"versions"`
	Time     map[string]string  `json:"time"`
}

type Version struct {
	Version string `json:"version"`
}

func GetPackageInfo(name string) (Package, error) {
	res, err := client.Get(fmt.Sprintf("https://registry.npmjs.org/%s", name))
	if err != nil {
		// TODO: handle timeouts maybe?
		return Package{}, oops.New(err, "failed to fetch NPM package")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Package{}, oops.New(err, "failed to read response body")
	}

	var p Package
	err = json.Unmarshal(body, &p)
	if err != nil {
		return Package{}, oops.New(err, "NPM response was not valid JSON")
	}

	return p, nil
}
