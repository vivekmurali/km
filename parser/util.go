package parser

import (
	"log"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type KM struct {
	Front Front
	MD    string
}

type Front struct {
	Title string
	Date  time.Time
	Tags  []string
}

func handle(front, md interface{}) interface{} {
	mdSlice := toIfaceSlice(md)
	frontSlice := toIfaceSlice(front)

	f := frontMatter(frontSlice)

	var markdown strings.Builder

	for _, v := range mdSlice {

		text, ok := v.(string)
		if !ok {
			lfSlice := toIfaceSlice(v)
			lfChar := lfSlice[1].([]uint8)
			if string(lfChar) == "\n" {
				markdown.WriteString(string(lfChar))
			}
		}

		markdown.WriteString(text)

	}

	log.Printf("%+v", f)

	kmData := KM{MD: markdown.String(), Front: f}
	return kmData
}

func frontMatter(f []interface{}) Front {

	frontStrings := toIfaceSlice(f[2])

	var frontBuilder strings.Builder

	for _, v := range frontStrings {
		text := v.(string)
		frontBuilder.WriteString(text)
	}

	var frontObject Front

	_, err := toml.Decode(frontBuilder.String(), &frontObject)
	if err != nil {
		log.Fatal(err)
	}

	return frontObject
}

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}
