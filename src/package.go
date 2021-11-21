package rh

import (
	"encoding/json"
	"html/template"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bvisness/releasehistory/src/npm"
)

type Package struct {
	Name     string    `json:"name"`
	Versions []Version `json:"versions"`
}

type Version struct {
	VersionString string `json:"version"`
	Time          int64  `json:"time"`

	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

func PackageHandler(rw http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/p/")
	npmPackage, err := npm.GetPackageInfo(name)
	if err != nil {
		// TODO: Handle bad packages
		internalError(rw, err)
		return
	}

	p := Package{
		Name:     npmPackage.Name,
		Versions: make([]Version, 0, len(npmPackage.Versions)),
	}

	for _, npmVersion := range npmPackage.Versions {
		major, minor, patch, other, ok := ParseNPMSemver(npmVersion.Version)
		if !ok {
			continue
		}

		if other != "" {
			// this is prerelease or something; ignore it
			continue
		}

		t, err := time.Parse(time.RFC3339Nano, npmPackage.Time[npmVersion.Version])
		if err != nil {
			internalError(rw, err)
		}

		p.Versions = append(p.Versions, Version{
			VersionString: npmVersion.Version,
			Time:          t.Unix(),

			Major: major,
			Minor: minor,
			Patch: patch,
		})
	}
	sort.Slice(p.Versions, func(i, j int) bool {
		return p.Versions[i].Time < p.Versions[j].Time
	})

	//
	// Precompute data for template purposes
	//

	type TemplateData struct {
		P     Package
		PJSON string
	}
	td := TemplateData{
		P: p,
	}

	pjsonBytes, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	td.PJSON = string(pjsonBytes)

	//
	// Template and serve
	//

	// TODO: Don't load and parse on every request
	tmpl, err := template.New("package.html").ParseFiles(
		"src/frontend/base.html",
		"src/frontend/package.html",
	)
	if err != nil {
		internalError(rw, err)
		return
	}

	err = tmpl.Execute(rw, td)
	if err != nil {
		internalError(rw, err)
		return
	}
}

var RESemver = regexp.MustCompile(`^[v=]?(\d+)\.(\d+)\.(\d+)(.*)`)

func ParseNPMSemver(version string) (major, minor, patch int, other string, ok bool) {
	groups := RESemver.FindStringSubmatch(version)
	if groups == nil {
		return 0, 0, 0, "", false
	}

	major, _ = strconv.Atoi(groups[1])
	minor, _ = strconv.Atoi(groups[2])
	patch, _ = strconv.Atoi(groups[3])
	other = groups[4]
	ok = true
	return
}

func containsInt(s []int, v int) bool {
	for _, value := range s {
		if value == v {
			return true
		}
	}
	return false
}
