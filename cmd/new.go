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
	err := createIfNotExists("notes")
	if err != nil {
		return err
	}

	// Generate random number and convert to string
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(1000)
	numString := strconv.Itoa(num)

	fileName := fmt.Sprintf("notes/%s.km", numString)

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	date := time.Now().Format("2006-01-02")
	// Write frontmatter to file
	s := fmt.Sprintf("---\ntitle: \ndate: %s \ntags: \n---\n", date)
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
