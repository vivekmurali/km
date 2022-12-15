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
	// 	kmData := text.(km)
	// 	if kmData.md != "\n# Test\n" {
	// 		t.Fatal("Expected # Test but got: ", text)
	// 	}
	// })

	t.Run("HEADING1", func(t *testing.T) {

		text, err := Parse("HEADING1", []byte("---\nTags=['Hi']\n---\n hello\n\n"))
		if err != nil {
			t.Fatal("Unexpected err: ", err)
		}
		kmData := text.(KM)
		if kmData.md != " hello\n\n" {
			t.Fatal("Expected hello but got: ", text)
		}
		if kmData.front.Tags[0] != "Hi" {
			t.Fatal("Not the right tags")
		}

	})
}
