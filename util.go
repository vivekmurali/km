package parser

import (
	"strings"
	"time"
)

type km struct {
	front front
	md    string
}

type front struct {
	date time.Time
	tags []string
}

func handle(front, md interface{}) interface{} {
	mdSlice := toIfaceSlice(md)

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
	return markdown.String()
}

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}
