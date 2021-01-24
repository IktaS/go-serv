package serv

import "github.com/alecthomas/participle/v2"

// Parse will parse bytes and return and ast of Serv
func Parse(input []byte) (*Serv, error) {
	parser, err := participle.Build(&Serv{})
	if err != nil {
		return nil, err
	}
	serv := &Serv{}
	err = parser.ParseBytes("", input, serv)
	if err != nil {
		return nil, err
	}
	return serv, nil
}
