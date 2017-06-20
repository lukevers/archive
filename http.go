package main

import (
	"encoding/base64"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"html/template"
	"net/http"
)

var templates *template.Template

func route() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", root)
	r.Get("/email/:id", single)

	var err error
	templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}

	// TODO: not hardcoded
	http.ListenAndServe(":4444", r)
}

func root(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "root", MB)
}

func single(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var email *Email

	for _, e := range MB.Emails {
		if e.Message.MessageId == id {
			email = e
			break
		}
	}

	if email == nil {
		// TODO: 404
	} else {
		// Convert if it's base64 encoded
		// TODO: cleaner
		dc, err := base64.StdEncoding.DecodeString(email.Message.Text)
		if err == nil {
			email.Message.Text = string(dc)
		}

		templates.ExecuteTemplate(w, "single", email)
	}
}
