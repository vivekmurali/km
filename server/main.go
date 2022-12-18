package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	r := chi.NewRouter()

	s := &app{}
	s.dbInit()
	defer s.closeDB()

	r.Use(middleware.Logger)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "server/static"))
	fmt.Println(filepath.Join(workDir, "static"))
	FileServer(r, "/static", filesDir)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/notes", http.StatusTemporaryRedirect)
	})
	r.Get("/new", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("server/templates/new.html")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}
		tmpl.Execute(w, nil)
	})

	r.Route("/notes", func(r chi.Router) {
		r.Get("/", s.getNotes)
		r.Post("/", s.postNote)
		r.Get("/{id}", s.singleNote)
	})
	http.ListenAndServe(os.Getenv("PORT"), r)
}
