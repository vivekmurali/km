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
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	fmt.Println("Committing all notes")
	err = createIfNotExists(home + "/notes/archive")
	if err != nil {
		return err
	}

	files, err := os.ReadDir(home + "/notes")
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, file := range files {
		if !file.IsDir() {
			wg.Add(1)
			// TODO: goroutine
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

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/notes/%s", home, f.Name())
	archivePath := fmt.Sprintf("%s/notes/archive/%s", home, f.Name())
	text, err := parser.ParseFile(path)
	if err != nil {
		return err
	}
	km := text.(parser.KM)

	var n note

	n.Title = km.Front.Title
	n.Tags = km.Front.Tags
	n.Content = km.MD
	n.Protected = km.Front.Protected

	body, err := json.Marshal(n)
	if err != nil {
		return err
	}

	env := os.Getenv("APP_ENV")
	var url string
	if env == "dev" {
		url = "http://localhost:3000/notes"
	} else {
		url = "https://notes.vivekmurali.in/notes"
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
