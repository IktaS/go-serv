package serv

import (
	"errors"

	"github.com/alecthomas/participle/v2"
)

// Parser is a struct that will have all .serv parser element
type Parser struct {
	parser *participle.Parser
}

// NewServParser creates a new ProtobufParser
func NewServParser() (*Parser, error) {
	parser, err := participle.Build(&Gserv{})
	if err != nil {
		return nil, err
	}
	return &Parser{
		parser: parser,
	}, nil
}

// Parse will parse a byte and return an AST of Proto
func (p *Parser) Parse(input []byte) (*Gserv, error) {
	proto := &Gserv{}
	err := p.parser.ParseBytes("", input, proto)
	if err != nil {
		return nil, err
	}
	err = checkDuplicateReference(proto)
	if err != nil {
		return nil, err
	}
	return proto, nil
}

func checkDuplicateReference(s *Gserv) error {
	ServiceFlag := make(map[string]bool)
	MessageFlag := make(map[string]bool)
	for _, def := range s.Definitions {
		if def.Service != nil {
			_, ok := ServiceFlag[def.Service.Name]
			if ok {
				return errors.New("Duplicate Service")
			}
			ServiceFlag[def.Service.Name] = true
		}
		if def.Message != nil {
			_, ok := MessageFlag[def.Message.Name]
			if ok {
				return errors.New("Duplicate Message")
			}
			MessageFlag[def.Message.Name] = true
		}
	}
	for _, def := range s.Definitions {
		if def.Service != nil {
			if def.Service.Request != nil {
				for _, t := range def.Service.Request {
					if t.Reference != "" {
						_, ok := MessageFlag[t.Reference]
						if !ok {
							return errors.New("Reference not found")
						}
					}
				}
			}
			if def.Service.Response != nil {
				if def.Service.Response.Reference != "" {
					_, ok := MessageFlag[def.Service.Response.Reference]
					if !ok {
						return errors.New("Reference not found")
					}
				}
			}
		}
	}
	return nil
}
