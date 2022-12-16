package cmd

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
)

func New(ctx *cli.Context) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	err = createIfNotExists(home + "/notes")
	if err != nil {
		return err
	}

	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(1000)
	numString := strconv.Itoa(num)

	fileName := fmt.Sprintf("%s/notes/%s.km", home, numString)

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	s := fmt.Sprintf("---\ntitle=\"\"\nprotected= false\ntags= []\n---\n")
	_, err = f.WriteString(s)

	fmt.Println("Created file: ", fileName)
	return nil
}

func createIfNotExists(path string) error {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}
