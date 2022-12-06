package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "km",
		Usage: "Knowledge Management",
		Action: func(*cli.Context) error {
			fmt.Println("Knowledge Management")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
