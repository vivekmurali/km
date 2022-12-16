package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/urfave/cli/v2"
	"github.com/vivekmurali/km/parser"
)

type note struct {
	Title     string
	Tags      []string
	Content   string
	Protected bool
}

func Commit(ctx *cli.Context) error {
	fmt.Println("Committing all notes")
	err := createIfNotExists("notes/archive")
	if err != nil {
		return err
	}

	files, err := os.ReadDir("notes")
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, file := range files {
		if !file.IsDir() {
			wg.Add(1)
			err = serverCommit(file, &wg)
			if err != nil {
				return err
			}
		}
	}
	wg.Wait()
	return nil
}

func serverCommit(f os.DirEntry, wg *sync.WaitGroup) error {
	defer wg.Done()

	path := fmt.Sprintf("notes/%s", f.Name())
	archivePath := fmt.Sprintf("notes/archive/%s", f.Name())
	text, err := parser.ParseFile(path)
	if err != nil {
		return err
	}
	km := text.(parser.KM)

	var n note

	n.Title = km.Front.Title
	n.Tags = km.Front.Tags
	n.Content = km.MD

	body, err := json.Marshal(n)
	if err != nil {
		return err
	}

	env := os.Getenv("APP_ENV")
	var url string
	if env == "prod" {
		url = "https://notes.vivekmurali.in/notes"
	} else {
		url = "http://localhost:3000/notes"
	}
	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return errors.New("Error posting to server")
	}

	err = moveFile(path, archivePath)
	if err != nil {
		return err
	}

	fmt.Println("Commited: ", f.Name())
	return nil
}

func moveFile(prev, next string) error {
	err := os.Rename(prev, next)
	if err != nil {
		return err
	}
	return nil
}
