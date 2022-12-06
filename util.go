package parser

type heading struct {
	level   int
	content string
}

func newHeading(level int, content string) (*heading, error) {
	return &heading{level, content}, nil
}
