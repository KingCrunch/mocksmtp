package main

import (
	"fmt"
	"net/http"
	"strconv"
	"log"

	"github.com/KingCrunch/mocksmtp/store"
	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
)

func RunHttpServer(bind string, s store.Store) {
	router := httprouter.New()

	// Serve static assets via the "static" directory
	router.ServeFiles("/static/*filepath", assetFS())

	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		l, err := s.List()
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}


		err = templates.ExecuteTemplate(w, "templates/index.html", l)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	router.GET("/mail/meta/:id", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id, err := uuid.FromString(params.ByName("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var m *store.Item
		m, err = s.Get(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err = templates.ExecuteTemplate(w, "templates/meta.html", m)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	router.GET("/mail/multi/:id/:part", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id, err := uuid.FromString(params.ByName("id"))
		index, err := strconv.Atoi(params.ByName("part"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var m *store.Item
		m, err = s.Get(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		ctype, _, e2 := m.Message.Parts[index].Header.ContentType()
		if e2 != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", ctype)
		fmt.Fprint(w, string(m.Message.Parts[index].Body))
	})
	router.GET("/mail/single/:id", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id, err := uuid.FromString(params.ByName("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var m *store.Item
		m, err = s.Get(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		ctype, _, e2 := m.Message.Header.ContentType()
		if e2 != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", ctype)
		fmt.Fprint(w, string(m.Message.Body))
	})

	router.DELETE("/mail", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		s.Purge()

		w.WriteHeader(http.StatusOK)
	})

	err := http.ListenAndServe(bind, router)
	if err != nil {
		log.Fatal(err)
	}
}
