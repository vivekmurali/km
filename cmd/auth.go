package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/urfave/cli/v2"
)

func Auth(ctx *cli.Context) error {

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	err = createIfNotExists(home + "/.config/km")
	if err != nil {
		return err
	}

	otp := ctx.Args().Get(0)
	if otp == "" {
		return errors.New("Invalid OTP")
	}

	v := url.Values{}
	v.Set("otp", otp)
	res, err := http.PostForm("https://notes.vivekmurali.in/login", v)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return errors.New("Error posting to server")
	}

	file, err := os.Create(home + "/.config/km/host")
	if err != nil {
		return err
	}

	var str string

	for _, v := range res.Cookies() {
		if v.Name == "session" {
			str = v.Value
		}
	}
	fmt.Fprint(file, str)

	return nil

}
