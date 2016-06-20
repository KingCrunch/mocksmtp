package main

import (
	"fmt"
	"strconv"
	"net/http"

	"github.com/satori/go.uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/KingCrunch/mocksmtp/store"
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

		var m *store.Item
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
		index,err := strconv.Atoi(params.ByName("part"))
		if (err != nil) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var m *store.Item
		m, err = s.Get(id)
		if (err != nil) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		ctype, _, _ := m.Message.Parts[index].Header.ContentType()
		w.Header().Add("Content-Type", ctype)
		fmt.Fprint(w, string(m.Message.Parts[index].Body))
	})
	router.GET("/mail/single/:id", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id, err := uuid.FromString(params.ByName("id"))
		if (err != nil) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var m *store.Item
		m, err = s.Get(id)
		if (err != nil) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		ctype, _, _ := m.Message.Header.ContentType()
		w.Header().Add("Content-Type", ctype)
		fmt.Fprint(w, string(m.Message.Body))
	})

	router.DELETE("/mail", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		s.Purge()

		w.WriteHeader(http.StatusOK)
	})

	err := http.ListenAndServe(bind, router)
	check(err)
}
