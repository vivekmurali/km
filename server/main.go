package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
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

	migration := os.Getenv("KM_MIGRATION")
	if migration == "TRUE" {
		s.migrateSearch()
		return
	}

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
		r.Get("/protected", s.getProtectedNotes)
	})

	//Unprotected routes
	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/notes", http.StatusTemporaryRedirect)
		})

		r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "server/static/favicon.ico")
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
		r.Get("/tags", s.getTags)
		r.Get("/tags/{tag}", s.getTagNotes)
		r.Get("/search", s.search)
		r.Get("/notes", s.getNotes)
		r.Get("/notes/{id}", s.singleNote)
	})

	http.ListenAndServe(os.Getenv("PORT"), r)
}

func searchInit() (bleve.Index, error) {
	mapping := bleve.NewIndexMapping()
	var index bleve.Index
	var err error

	index, err = bleve.Open("index.bleve")
	if err != nil {
		index, err = bleve.New("index.bleve", mapping)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return index, nil
	}
	return nil, nil
}

func (s *app) migrateSearch() {
	n, err := s.getNotesFromDB(0, false)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range n {
		singleNote, err := s.getSingleNote(v.ID)
		if err != nil {
			log.Fatal(err)
		}
		singleNote.HTML = ""
		Index.Index(fmt.Sprintf("%d", v.ID), singleNote)
	}
}
