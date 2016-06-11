package main

import (
	"net/http"
	"fmt"
	"html"
	"github.com/satori/go.uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/KingCrunch/visualsmtp/store"
	"github.com/KingCrunch/visualsmtp/mail"
)

func RunHttpServer(bind string, s store.Store) {
	router := httprouter.New()

	// Serve static assets via the "static" directory
	router.ServeFiles("/static/*filepath", assetFS())

	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		l, err := s.List()
		check(err)
		err = templates.ExecuteTemplate(w, "templates/index.html", l)
		check(err)
	})
	router.GET("/mail/meta/:id", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id, err := uuid.FromString(params.ByName("id"))
		check(err)

		var m mail.Mail
		m, err = s.Get(id)
		check(err)

		data := struct{
			Id uuid.UUID
			Mail mail.Mail
		}{
			id,
			m,
		}
		err = templates.ExecuteTemplate(w, "templates/meta.html", data)
		check(err)
	})
	router.GET("/mail/multi/:id/:part", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id, err := uuid.FromString(params.ByName("id"))
		check(err)

		var m mail.Mail
		m, err = s.Get(id)
		check(err)

		fmt.Fprint(w, html.EscapeString(string(m.Data)))

		// TODO Parse multipart-request
	})

	err := http.ListenAndServe(bind, router)
	check(err)
}
