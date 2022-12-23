package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	r := chi.NewRouter()

	s := &app{}
	s.store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	s.dbInit()
	defer s.closeDB()

	r.Use(middleware.Logger)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "server/static"))
	fmt.Println(filepath.Join(workDir, "static"))
	FileServer(r, "/static", filesDir)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(s.UnloggedInRedirector)

		r.Get("/new", func(w http.ResponseWriter, r *http.Request) {
			tmpl, err := template.ParseFiles("server/templates/new.html")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
			}
			tmpl.Execute(w, nil)
		})

		r.Route("/edit/{id}", func(r chi.Router) {
			r.Get("/", s.showEditNote)
			r.Post("/", s.editNote)
		})

		r.Delete("/delete/{id}", s.deleteNote)
		r.Post("/notes", s.postNote)
	})

	//Unprotected routes
	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/notes", http.StatusTemporaryRedirect)
		})

		r.Group(func(r chi.Router) {

			r.Use(s.LoggedInRedirector)
			r.Route("/login", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					tmpl, err := template.ParseFiles("server/templates/login.html")
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						w.Write([]byte(err.Error()))
					}
					tmpl.Execute(w, nil)
				})
				r.Post("/", s.login)
			})

		})
		r.Get("/notes", s.getNotes)
		r.Get("/notes/{id}", s.singleNote)
	})

	http.ListenAndServe(os.Getenv("PORT"), r)
}
