package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gomarkdown/markdown"
)

type note struct {
	Title     string
	Tags      []string
	Content   string
	Protected bool
}

func (s *app) postNote(w http.ResponseWriter, r *http.Request) {

	var n note

	err := json.NewDecoder(r.Body).Decode(&n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.insertDB(n.Tags, n.Title, n.Content, n.Protected)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	n.Content = ""

	returnJson, err := json.Marshal(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(returnJson)
}

type notes struct {
	Title   string
	Tags    []string
	ID      int64
	Created string
	Content string
	HTML    template.HTML
}

func (s *app) singleNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	note, err := s.getSingleNote(intID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	note.Content = string(markdown.NormalizeNewlines([]byte(note.Content)))
	htmlData := markdown.ToHTML([]byte(note.Content), nil, nil)
	note.HTML = template.HTML((htmlData))

	tmpl, err := template.ParseFiles("server/templates/note.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not open template"))
	}
	err = tmpl.Execute(w, note)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not open template"))
	}
}

func (s *app) getNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := s.getNotesFromDB()
	if err != nil {
		log.Println("Error getting notes from DB", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("server/templates/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not open template"))
	}

	err = tmpl.Execute(w, notes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not open template"))
	}
}
