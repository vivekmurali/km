package parser

import (
	"testing"
)

func TestParse(t *testing.T) {

	// t.Run("TEXT", func(t *testing.T) {
	// 	text, err := ParseFile("test.km")
	// 	if err != nil {
	// 		t.Fatal("Unexpected err: ", err)
	// 	}
	// 	if text != "\n# Test\n" {
	// 		t.Fatal("Expected # Test but got: ", text)
	// 	}
	// })

	t.Run("HEADING1", func(t *testing.T) {

		text, err := Parse("HEADING1", []byte("---\n hi\n---\n hello\n\n"))
		if err != nil {
			t.Fatal("Unexpected err: ", err)
		}
		if text != " hello\n\n" {
			t.Fatal("Expected Hello but got: ", text)
		}

	})
}
