package proto

import (
	"testing"
)

func TestParse(t *testing.T) {
	input := []byte(`
		service TestService ( TestMessage1 ) : TestMessage1;
	`)
	println(input)
}
