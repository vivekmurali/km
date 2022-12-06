package parser

import (
	"testing"
)

func TestParse(t *testing.T) {

	t.Run("TEXT", func(t *testing.T) {
		// text, err := Parse("TEXT", []byte(`Hello \n`))
		text, err := ParseFile("test.km")
		if err != nil {
			t.Fatal("Unexpected err: ", err)
		}
		if text != "Hello" {
			t.Fatal("Expected Hello but got: ", text)
		}
	})

	t.Run("HEADING1", func(t *testing.T) {

		h, err := Parse("HEADING1", []byte("# Hi \n\n"))
		if err != nil {
			t.Fatal("Unexpected err: ", err)
		}
		t.Logf("heading of type: %T\n", h)
		t.Logf("heading = %+v", h)
		heading, ok := h.(heading)
		if !ok {
			t.Fatal("Could not typecast to heading")
		}
		t.Logf("heading = %+v", heading)
		if heading.level != 1 || heading.content != "Hi" {
			t.Fatal("Wrong level or content", heading)
		}
	})
}
