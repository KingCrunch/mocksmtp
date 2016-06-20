package main

import (
	"html/template"
	"log"
	"time"
	"github.com/veqryn/go-email/email"
)

var (
	templateMap = template.FuncMap{
		"datetime": func(t time.Time) string {
			return t.Format(time.RFC3339)
		},
		"contenttype": func(header email.Header) (string, error) {
			h, _, e := header.ContentType()
			return h, e
		},
		"contentdisposition": func(header email.Header) (string, error) {
			h, _, e := header.ContentDisposition()
			return h, e
		},
	}

	templates = template.New("").Funcs(templateMap)
)

func init() {
	for _, path := range AssetNames() {
		bytes, err := Asset(path)
		if err != nil {
			log.Panicf("Unable to parse: path=%s, err=%s", path, err)
		}
		templates.New(path).Parse(string(bytes))
	}
}
