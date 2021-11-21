package rh

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func Main() {
	serveFrontendFile("/tachyons.css", "text/css")
	serveFrontendFile("/style.css", "text/css")
	serveFrontendFile("/templates.js", "application/css")

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("404 Not Found"))
			return
		}

		// TODO: Don't load and parse on every request
		tmpl, err := template.New("index.html").ParseFiles(
			"src/frontend/base.html",
			"src/frontend/index.html",
		)
		if err != nil {
			internalError(rw, err)
			return
		}

		tmpl.Execute(rw, map[string]interface{}{
			"Now": time.Now(),
		})
	})

	http.HandleFunc("/p/", PackageHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func internalError(rw http.ResponseWriter, err error) {
	log.Print(err)
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Write([]byte("Internal Server Error"))
}

func serveFrontendFile(path string, contentType string) {
	http.HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", contentType)

		f, err := os.Open(filepath.Join("src/frontend", path))
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(rw, f)
		if err != nil {
			panic(err)
		}
	})
}
