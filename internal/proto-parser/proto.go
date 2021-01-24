package proto

import (
	"github.com/alecthomas/participle/v2"
)

// Parse will parse a byte and return an AST of Proto
func Parse(input []byte) (*Proto, error) {
	parser, err := participle.Build(&Proto{})
	if err != nil {
		return nil, err
	}
	proto := &Proto{}
	err = parser.ParseBytes("", input, proto)
	if err != nil {
		return nil, err
	}
	return proto, nil
}
