//revive:disable
package serv

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

//Gserv is entry point to our parser
type Gserv struct {
	Definitions []*Definition `( @@ ";"* )*`
}

//Definition is each definition in our service
type Definition struct {
	InboundService  *InboundService  `@@*`
	OutboundService *OutboundService `@@*`
	Message         *Message         `@@*`
}

// Service is a definition of a service
type InboundService struct {
	Name     string  `"inbound" @Ident`
	Request  []*Type `"(" ( @@ ","* )* ")"`
	Response *Type   `":"? @@?`
}

// Service is a definition of a service
type OutboundService struct {
	Name     string  `"outbound" @Ident`
	Request  []*Type `"(" ( @@ ","* )* ")"`
	Response *Type   `":"? @@?`
}

// Message is a definition of a message
type Message struct {
	Name        string               `"message" @Ident`
	Definitions []*MessageDefinition `"{" @@* "}"`
}

// MessageDefinition is a definition for messages
type MessageDefinition struct {
	// Message *Message `( @@`
	// Enum    *Enum    ` | @@`
	Field *Field `( @@ ) ";"*`
}

// // Enum is enum
// type Enum struct {
// 	Name   string       `"enum" @Ident`
// 	Values []*EnumValue `"{" ( @@ ( ";" )* )* "}"`
// }

// // EnumValue is enum value
// type EnumValue struct {
// 	Key   string `@Ident`
// 	Value int    `"=" @( [ "-" ] Int )`
// }

// Field is a field
type Field struct {
	Optional bool `( @"optional"`
	Required bool ` | @"required" )?`

	Type *Type  `@@`
	Name string `@Ident`
}

type Scalar int

const (
	None Scalar = iota
	Double
	Float
	Int32
	Int64
	Uint32
	Uint64
	Sint32
	Sint64
	Fixed32
	Fixed64
	SFixed32
	SFixed64
	Bool
	String
	Bytes
)

var scalarToGoString = map[Scalar]string{
	None: "None", Double: "Double", Float: "Float", Int32: "Int32", Int64: "Int64", Uint32: "Uint32",
	Uint64: "Uint64", Sint32: "Sint32", Sint64: "Sint64", Fixed32: "Fixed32", Fixed64: "Fixed64",
	SFixed32: "SFixed32", SFixed64: "SFixed64", Bool: "Bool", String: "String", Bytes: "Bytes",
}

var scalarToString = map[Scalar]string{
	None: "None", Double: "float64", Float: "float32", Int32: "int32", Int64: "int64", Uint32: "uint32",
	Uint64: "uint64", Sint32: "int32", Sint64: "int64", Fixed32: "uint32", Fixed64: "uint64",
	SFixed32: "int32", SFixed64: "int64", Bool: "bool", String: "string", Bytes: "[]byte",
}

func (s Scalar) GoString() string { return scalarToGoString[s] }
func (s Scalar) String() string   { return scalarToString[s] }

var StringToScalar = map[string]Scalar{
	"double": Double, "float": Float, "int32": Int32, "int64": Int64, "uint32": Uint32, "uint64": Uint64,
	"sint32": Sint32, "sint64": Sint64, "fixed32": Fixed32, "fixed64": Fixed64, "sfixed32": SFixed32,
	"sfixed64": SFixed64, "bool": Bool, "string": String, "bytes": Bytes,
}

func (s *Scalar) Parse(lex *lexer.PeekingLexer) error {
	token, err := lex.Peek(0)
	if err != nil {
		return err
	}
	v, ok := StringToScalar[token.Value]
	if !ok {
		return participle.NextMatch
	}
	_, err = lex.Next()
	if err != nil {
		return err
	}
	*s = v
	return nil
}

type Type struct {
	Scalar    Scalar `  @@`
	Reference string `| @(Ident ( "." Ident)*)`
}

func (s Type) String() string {
	if s.Reference == "" {
		return s.Scalar.String()
	} else {
		return s.Reference
	}
}
