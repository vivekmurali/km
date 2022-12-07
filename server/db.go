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

func (s *app) insertDB(name, content string, protected bool) error {

	tag, err := s.db.Exec(context.Background(), "insert into notes(NAME, CONTENT, PROTECTED) values($1, $2, $3)", &name, &content, &protected)

	if err != nil {
		return err
	}

	if !tag.Insert() {
		return errors.New("Not insert")
	}

	return nil
}
