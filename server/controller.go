package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/pquerna/otp/totp"
)

type notes struct {
	Title     string
	Tags      []string
	ID        int64
	Created   string
	Content   string
	Protected bool
	HTML      template.HTML
}

func (s *app) postNote(w http.ResponseWriter, r *http.Request) {

	var n notes

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

func (s *app) singleNote(w http.ResponseWriter, r *http.Request) {

	session, _ := s.store.Get(r, "session")
	authenticated, ok := session.Values["authenticated"].(bool)

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

	if note.Protected && (!authenticated || !ok) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Not authorized")
		return
	}

	//All links open in new tab
	htmlFlags := html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// Difference between \n and \r\n
	note.Content = string(markdown.NormalizeNewlines([]byte(note.Content)))

	// Add footnotes support
	extensions := parser.CommonExtensions | parser.Footnotes
	parser := parser.NewWithExtensions(extensions)

	htmlData := markdown.ToHTML([]byte(note.Content), parser, renderer)

	// HTML will not render if you do not do this (security)
	note.HTML = template.HTML((htmlData))
	note.ID = intID

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

func (s *app) getProtectedNotes(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	var intPage int
	var err error

	if page == "" {
		intPage = 0
	} else {
		intPage, err = strconv.Atoi(page)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}
	notes, err := s.getNotesFromDB(intPage, true)
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

func (s *app) getNotes(w http.ResponseWriter, r *http.Request) {

	page := r.URL.Query().Get("page")
	var intPage int
	var err error

	if page == "" {
		intPage = 0
	} else {
		intPage, err = strconv.Atoi(page)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}

	notes, err := s.getNotesFromDB(intPage, false)
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

func (s *app) login(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session")

	secret := os.Getenv("TOTP_SECRET")

	otp := r.FormValue("otp")

	if !totp.Validate(otp, secret) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unable to login"))
		return
	}

	session.Values["authenticated"] = true
	session.Save(r, w)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged in"))
}

func (s *app) editNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	intID, err := strconv.ParseInt(id, 10, 64)
	var n notes

	err = json.NewDecoder(r.Body).Decode(&n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.editInDB(intID, n.Tags, n.Title, n.Content, n.Protected)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Edited note"))
}

func (s *app) showEditNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	intID, err := strconv.ParseInt(id, 10, 64)

	tmpl, err := template.ParseFiles("server/templates/edit.html")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	note, err := s.getSingleNote(intID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	note.ID = intID
	tmpl.Execute(w, note)
}

func (s *app) deleteNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	intID, err := strconv.ParseInt(id, 10, 64)

	err = s.deleteNoteFromDB(intID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted"))
}

func (s *app) search(w http.ResponseWriter, r *http.Request) {

	term := r.URL.Query().Get("q")

	if term == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Empty query"))
		return
	}

	sp := strings.Split(term, " ")
	term = strings.Join(sp, " or ")

	notes, err := s.searchDB(term)
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

func (s *app) getTags(w http.ResponseWriter, r *http.Request) {

	tags, err := s.tags()
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("server/templates/tags.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not open template"))
	}

	err = tmpl.Execute(w, tags)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not open template"))
	}

}

func (s *app) getTagNotes(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")

	notes, err := s.tagNotes(tag)
	if err != nil {
		w.Write([]byte(err.Error()))
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
