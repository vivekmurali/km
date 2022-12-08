package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgtype"
)

type note struct {
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

	err = s.insertDB(n.Tags, n.Content, n.Protected)
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
	Tags    []string
	Content string
	Created pgtype.Timestamptz
}

func (s *app) getNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := s.getNotesFromDB()
	if err != nil {
		log.Println("Error getting notes from DB", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
