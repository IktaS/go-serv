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
	InboundServiceFlag := make(map[string]bool)
	OutboundServiceFlag := make(map[string]bool)
	MessageFlag := make(map[string]bool)
	for _, def := range s.Definitions {
		if def.InboundService != nil {
			_, ok := InboundServiceFlag[def.InboundService.Name]
			if ok {
				return errors.New("Duplicate Inbound Service")
			}
			InboundServiceFlag[def.InboundService.Name] = true
		}
		if def.OutboundService != nil {
			_, ok := OutboundServiceFlag[def.OutboundService.Name]
			if ok {
				return errors.New("Duplicate Outbound Service")
			}
			OutboundServiceFlag[def.OutboundService.Name] = true
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
		if def.InboundService != nil {
			if def.InboundService.Request != nil {
				for _, t := range def.InboundService.Request {
					if t.Reference != "" {
						_, ok := MessageFlag[t.Reference]
						if !ok {
							return errors.New("Reference not found")
						}
					}
				}
			}
			if def.InboundService.Response != nil {
				if def.InboundService.Response.Reference != "" {
					_, ok := MessageFlag[def.InboundService.Response.Reference]
					if !ok {
						return errors.New("Reference not found")
					}
				}
			}
		}
		if def.OutboundService != nil {
			if def.OutboundService.Request != nil {
				for _, t := range def.OutboundService.Request {
					if t.Reference != "" {
						_, ok := MessageFlag[t.Reference]
						if !ok {
							return errors.New("Reference not found")
						}
					}
				}
			}
			if def.OutboundService.Response != nil {
				if def.OutboundService.Response.Reference != "" {
					_, ok := MessageFlag[def.OutboundService.Response.Reference]
					if !ok {
						return errors.New("Reference not found")
					}
				}
			}
		}
	}
	return nil
}
