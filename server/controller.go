package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"
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
	Content string
	Created time.Time
}

func (s *app) getNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := s.getNotesFromDB()
	if err != nil {
		log.Println("Error getting notes from DB", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// err = json.NewEncoder(w).Encode(notes)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

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
