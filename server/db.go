package main

import (
	"context"
	"errors"
	"fmt"
	"os"

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

func (s *app) insertDB(tags []string, content string, protected bool) error {

	tag, err := s.db.Exec(context.Background(), "insert into notes(TAGS, CONTENT, PROTECTED) values($1, $2, $3)", &tags, &content, &protected)

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

	rows, err := s.db.Query(context.Background(), "select tags, content, created from notes where protected=false limit 10")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var singlenote notes
		err = rows.Scan(&singlenote.Tags, &singlenote.Content, &singlenote.Created)
		if err != nil {
			return nil, err
		}
		n = append(n, singlenote)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return n, nil
}
