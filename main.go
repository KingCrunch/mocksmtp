package main

import (
	"bitbucket.org/chrj/smtpd"
	"net/http"
	"fmt"
	"path"
	"log"
	"html/template"
	"strconv"
	"html"
)

var bucket []smtpd.Envelope = make([]smtpd.Envelope, 0, 5)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.New("index.html").ParseFiles("templates/index.html")
		check(err)
		err = t.Execute(w, bucket)
		check(err)
	})
	http.HandleFunc("/mail/", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(path.Base(r.URL.Path))
		check(err)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, html.EscapeString(string(bucket[id].Data)))
	})
	go func() {
		err := http.ListenAndServe(":8080", nil)
		check(err)
	}();

	server := &smtpd.Server{
		Handler: func(peer smtpd.Peer, env smtpd.Envelope) error {
			log.Printf("New Mail from %q to %q", env.Sender, env.Recipients)
			bucket = append(bucket, env)

			return nil
		},
	}
	err := server.ListenAndServe("127.0.0.1:10025")
	check(err)
}

func check (err error) {
	if err != nil {
		log.Fatal(err)
	}
}