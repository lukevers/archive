package main

import (
	"encoding/base64"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"html/template"
	"net/http"
)

var (
	T  *template.Template
	MB *Mailbox
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", root)
	r.Get("/email/:id", single)

	var err error
	T, err = template.ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}

	MB = NewMailbox(&Mailbox{
		Dirs: []string{"./emails"},
	})

	http.ListenAndServe(":4444", r)
}

func root(w http.ResponseWriter, r *http.Request) {
	var err error
	T, err = template.ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}

	T.ExecuteTemplate(w, "root", MB)
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

	//	log.Println(email.Message)

	if email == nil {
		//
	} else {
		// Convert if it's base64 encoded
		dc, err := base64.StdEncoding.DecodeString(email.Message.Text)
		if err == nil {
			email.Message.Text = string(dc)
		}

		T.ExecuteTemplate(w, "single", email)
	}
}
