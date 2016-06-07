package main

import (
	"html/template"
	"log"
	"time"
)

var (
	templateMap = template.FuncMap{
		"datetime": func(t time.Time) string {
			return t.Format(time.RFC3339)
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
