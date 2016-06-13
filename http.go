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
		if (err != nil) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var m mail.Mail
		m, err = s.Get(id)
		if (err != nil) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err = templates.ExecuteTemplate(w, "templates/meta.html", m)
		check(err)
	})
	router.GET("/mail/multi/:id/:part", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id, err := uuid.FromString(params.ByName("id"))
		if (err != nil) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var m mail.Mail
		m, err = s.Get(id)
		if (err != nil) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		fmt.Fprint(w, html.EscapeString(string(m.Data)))

		// TODO Parse multipart-request
	})

	router.DELETE("/mail", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		s.Purge()

		w.WriteHeader(http.StatusOK)
	})

	err := http.ListenAndServe(bind, router)
	check(err)
}
