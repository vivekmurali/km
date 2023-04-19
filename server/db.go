package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
)

type app struct {
	db    *pgxpool.Pool
	store *sessions.CookieStore
}

var Index bleve.Index

func (s *app) dbInit() error {
	var err error
	s.db, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	var greeting string
	err = s.db.QueryRow(context.Background(), "select 'Database Init!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	index, err := bleve.Open("index.bleve")
	if err == bleve.ErrorIndexPathDoesNotExist {
		log.Println("Creating new index")
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New("index.bleve", mapping)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	Index = index
	fmt.Println(greeting)
	return nil
}

func (s *app) closeDB() {
	s.db.Close()
}

func (s *app) editInDB(id int64, tags []string, title, content string, protected bool) error {

	tag, err := s.db.Exec(context.Background(), "update notes set tags = $1, title = $2, content = $3, protected = $4 where id = $5", tags, title, content, protected, id)

	if err != nil {
		return err
	}

	if !tag.Update() {
		return errors.New("Not update")
	}
	return nil
}

func (s *app) insertDB(tags []string, title, content string, protected bool) (int, error) {

	decodedContent, err := url.QueryUnescape(content)
	if err != nil {
		return 0, err
	}

	var id int
	err = s.db.QueryRow(context.Background(), "insert into notes(TITLE, TAGS, CONTENT, PROTECTED) values($1, $2, $3, $4) RETURNING id", &title, &tags, &decodedContent, &protected).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *app) getNotesFromDB(page int, protected bool) ([]notes, error) {
	var n []notes

	offset := page * 50
	rows, err := s.db.Query(context.Background(), "select id, title, created, tags from notes where protected=$1 order by created desc limit 50 offset $2", protected, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var singlenote notes
		var created time.Time
		err = rows.Scan(&singlenote.ID, &singlenote.Title, &created, &singlenote.Tags)
		if err != nil {
			return nil, err
		}

		loc, _ := time.LoadLocation("America/New_York")

		singlenote.Created = created.In(loc).Format(time.RFC822)
		n = append(n, singlenote)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *app) getSingleNote(id int64) (notes, error) {
	var note notes
	var created time.Time
	err := s.db.QueryRow(context.Background(), "select title, created, tags, content, protected from notes where id=$1", id).Scan(&note.Title, &created, &note.Tags, &note.Content, &note.Protected)
	if err != nil {
		return note, err
	}

	loc, _ := time.LoadLocation("America/New_York")
	note.Created = created.In(loc).Format(time.RFC822)

	return note, nil
}

func (s *app) deleteNoteFromDB(id int64) error {

	tag, err := s.db.Exec(context.Background(), "delete from notes where id = $1", id)
	if err != nil {
		return err
	}

	if !tag.Delete() {
		return errors.New("Not delete")
	}

	return nil
}

func (s *app) searchDB(term string) ([]notes, error) {

	var n []notes

	rows, err := s.db.Query(context.Background(), "select * from search_notes($1)", term)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var singlenote notes
		var created time.Time
		var rank any
		// id integer, title text, tags text[], created timestamp with time zone, content text, rank real
		err = rows.Scan(&singlenote.ID, &singlenote.Title, &singlenote.Tags, &created, &singlenote.Content, &rank)
		if err != nil {
			return nil, err
		}

		loc, _ := time.LoadLocation("America/New_York")

		singlenote.Created = created.In(loc).Format(time.RFC822)
		n = append(n, singlenote)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return n, nil

}

func (s *app) tags() ([]string, error) {
	var tags []string

	rows, err := s.db.Query(context.Background(), "select distinct trim(unnest(tags)) as tags from notes where tags is not null and protected=false")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (s *app) tagNotes(tag string) ([]notes, error) {
	var n []notes

	rows, err := s.db.Query(context.Background(), "select id, title, created, tags from notes where $1=ANY(tags)", tag)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var singlenote notes
		var created time.Time
		err = rows.Scan(&singlenote.ID, &singlenote.Title, &created, &singlenote.Tags)
		if err != nil {
			return nil, err
		}

		loc, _ := time.LoadLocation("America/New_York")

		singlenote.Created = created.In(loc).Format(time.RFC822)
		n = append(n, singlenote)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return n, nil
}
