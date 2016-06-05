package main

import (
	"net/http"
	"fmt"
	"html"
	"path"
	"html/template"
	"github.com/satori/go.uuid"
)

func RunHttpServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.New("index.html").ParseFiles("templates/index.html", "templates/style.css", "templates/script.js")
		check(err)
		err = t.Execute(w, mailBucket)
		check(err)
	})
	http.HandleFunc("/mail", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, html.EscapeString("Click on a mail"))
	})
	http.HandleFunc("/mail/meta/", func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.FromString(path.Base(r.URL.Path))
		check(err)
		t, err := template.New("meta.html").ParseFiles("templates/meta.html")
		check(err)

		err = t.Execute(w, mailBucket[id])
		check(err)
	})
	http.HandleFunc("/mail/", func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.FromString(path.Base(r.URL.Path))
		check(err)

		fmt.Fprint(w, html.EscapeString(string(mailBucket[id].Data)))

		// TODO Parse multipart-request
	})

	err := http.ListenAndServe(":8080", nil)
	check(err)
}
