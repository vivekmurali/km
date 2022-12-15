package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/vivekmurali/km/cmd"
)

func main() {
	app := &cli.App{
		Name:  "km",
		Usage: "Knowledge Management",
		Action: func(*cli.Context) error {
			fmt.Println("Knowledge Management")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "new",
				Aliases: []string{"n"},
				Usage:   "Create new KM",
				Action:  cmd.New,
			},
			{
				Name:    "commit",
				Aliases: []string{"c"},
				Usage:   "Commit all files to server",
				Action:  cmd.Commit,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
