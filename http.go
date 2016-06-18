package main

import (
	"fmt"
	"strconv"
	"log"
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/base64"
	"mime/quotedprintable"

	"github.com/satori/go.uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/KingCrunch/mocksmtp/store"
	"github.com/KingCrunch/mocksmtp/mail"
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
		index,err := strconv.Atoi(params.ByName("part"))
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
		log.Print(m.Parts[index].Header)
		w.Header().Add("Content-Type", m.Parts[index].Header.Get("Content-Type"))
		//w.Header().Add("Content-Transfer-Encoding", m.Parts[index].Header.Get("Content-Transfer-Encoding"))
		switch m.Parts[index].Header.Get("Content-Transfer-Encoding") {
		case "base64":
			x, _ := base64.RawStdEncoding.DecodeString(string(m.Parts[index].Data))
			fmt.Fprint(w, string(x))
		case "quoted-printable":
			x := quotedprintable.NewReader(bytes.NewReader(m.Parts[index].Data))
			x2, _ := ioutil.ReadAll(x)
			fmt.Fprint(w, string(x2))
		default:
			fmt.Fprint(w, string(m.Parts[index].Data))
		}
	})
	router.GET("/mail/single/:id", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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

		w.Header().Add("Content-Type", m.Header.Get("Content-Type"))
		switch m.Header.Get("Content-Transfer-Encoding") {
		case "base64":
			x, _ := base64.RawStdEncoding.DecodeString(string(m.Data))
			fmt.Fprint(w, string(x))
		case "quoted-printable":
			x := quotedprintable.NewReader(bytes.NewReader(m.Data))
			x2, _ := ioutil.ReadAll(x)
			fmt.Fprint(w, string(x2))
		default:
			fmt.Fprint(w, string(m.Data))
		}
	})

	router.DELETE("/mail", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		s.Purge()

		w.WriteHeader(http.StatusOK)
	})

	err := http.ListenAndServe(bind, router)
	check(err)
}
