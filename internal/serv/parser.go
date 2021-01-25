package serv

import "github.com/alecthomas/participle/v2"

// Parser is a struct that will have all .serv parser element
type Parser struct {
	parser *participle.Parser
}

// NewServParser creates a new ProtobufParser
func NewServParser() (*Parser, error) {
	parser, err := participle.Build(&Serv{})
	if err != nil {
		return nil, err
	}
	return &Parser{
		parser: parser,
	}, nil
}

// Parse will parse a byte and return an AST of Proto
func (p *Parser) Parse(input []byte) (*Serv, error) {
	proto := &Serv{}
	err := p.parser.ParseBytes("", input, proto)
	if err != nil {
		return nil, err
	}
	return proto, nil
}
