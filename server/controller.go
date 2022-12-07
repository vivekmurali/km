package main

import (
	"encoding/json"
	"net/http"
)

type note struct {
	Name      string
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

	err = s.insertDB(n.Name, n.Content, n.Protected)
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
