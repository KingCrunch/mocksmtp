package main

import (
	"net/http"
	"fmt"
	"html"
	"github.com/satori/go.uuid"
	"github.com/julienschmidt/httprouter"
)

func RunHttpServer(bind string) {
	router := httprouter.New()

	// Serve static assets via the "static" directory
	router.ServeFiles("/static/*filepath", assetFS())

	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		err := templates.ExecuteTemplate(w, "templates/index.html", mailBucket)
		check(err)
	})
	router.GET("/mail/meta/:id", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id, err := uuid.FromString(params.ByName("id"))
		check(err)

		err = templates.ExecuteTemplate(w, "templates/meta.html", mailBucket[id])
		check(err)
	})
	router.GET("/mail/multi/:id/:part", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id, err := uuid.FromString(params.ByName("id"))
		check(err)

		fmt.Fprint(w, html.EscapeString(string(mailBucket[id].Data)))

		// TODO Parse multipart-request
	})

	err := http.ListenAndServe(bind, router)
	check(err)
}
