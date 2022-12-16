package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type app struct {
	db *pgxpool.Pool
}

func (s *app) dbInit() {
	var err error
	s.db, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	var greeting string
	err = s.db.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)
}

func (s *app) closeDB() {
	s.db.Close()
}

func (s *app) insertDB(tags []string, title, content string, protected bool) error {

	decodedContent, err := url.QueryUnescape(content)
	if err != nil {
		return err
	}

	tag, err := s.db.Exec(context.Background(), "insert into notes(TITLE, TAGS, CONTENT, PROTECTED) values($1, $2, $3, $4)", &title, &tags, &decodedContent, &protected)

	if err != nil {
		return err
	}

	if !tag.Insert() {
		return errors.New("Not insert")
	}

	return nil
}

func (s *app) getNotesFromDB() ([]notes, error) {
	var n []notes

	rows, err := s.db.Query(context.Background(), "select id, title, created, tags from notes where protected=false order by created desc limit 30")
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

		singlenote.Created = created.Format(time.RFC822)
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
	err := s.db.QueryRow(context.Background(), "select title, created, tags, content from notes where id=$1", id).Scan(&note.Title, &created, &note.Tags, &note.Content)
	if err != nil {
		return note, err
	}

	note.Created = created.Format(time.RFC822)

	return note, nil
}
