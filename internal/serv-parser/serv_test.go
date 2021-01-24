package serv

import "testing"

func TestParse(t *testing.T) {
	input := []byte(`
		message TestMessage1 {
			int32 page_number;
		};
		message TestMessage2 {
			int32 page_number;
		};
		service TestService ( TestMessage1, TestMessage2 ) : TestMessage2;
	`)
	res, err := Parse(input)
	if err != nil {
		panic(err)
	}
	println(res.Definitions)
}
