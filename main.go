package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"github.com/vivekmurali/km/cmd"
)

func main() {

	godotenv.Load()

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
			{
				Name:   "auth",
				Usage:  "login",
				Action: cmd.Auth,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
