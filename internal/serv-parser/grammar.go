//revive:disable
package serv

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

//Serv is entry point to our parser
type Serv struct {
	Pos lexer.Position

	Definitions []*Definition `parser:"( @@ ";"* )*"`
}

//Definition is each definition in our service
type Definition struct {
	Pos lexer.Position

	Service *Service `parser:"@@*"`
	Message *Message `parser:"@@*"`
}

// Service is a definition of a service
type Service struct {
	Pos lexer.Position

	Name     string  `"service" @Ident`
	Request  []*Type `"(" ( @@ ","* )* ")"`
	Response *Type   `":" @@ `
}

// Message is a definition of a message
type Message struct {
	Pos lexer.Position

	Name        string               `"message" @Ident`
	Definitions []*MessageDefinition `"{" @@* "}"`
}

// MessageDefinition is a definition for messages
type MessageDefinition struct {
	Enum    *Enum    `( @@`
	Message *Message ` | @@`
	Field   *Field   ` | @@ ) ";"*`
}

// Enum is enum
type Enum struct {
	Pos lexer.Position

	Name   string            `"enum" @Ident`
	Values []*EnumDefinition `"{" ( @@ ( ";" )* )* "}"`
}

// Enum is definition of an enum
type EnumDefinition struct {
	Pos   lexer.Position
	Value *EnumValue ` @@`
}

// EnumValue is enum value
type EnumValue struct {
	Pos lexer.Position

	Key   string `@Ident`
	Value int    `"=" @( [ "-" ] Int )`
}

// Field is a field
type Field struct {
	Pos lexer.Position

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

var scalarToString = map[Scalar]string{
	None: "None", Double: "Double", Float: "Float", Int32: "Int32", Int64: "Int64", Uint32: "Uint32",
	Uint64: "Uint64", Sint32: "Sint32", Sint64: "Sint64", Fixed32: "Fixed32", Fixed64: "Fixed64",
	SFixed32: "SFixed32", SFixed64: "SFixed64", Bool: "Bool", String: "String", Bytes: "Bytes",
}

func (s Scalar) GoString() string { return scalarToString[s] }

var stringToScalar = map[string]Scalar{
	"double": Double, "float": Float, "int32": Int32, "int64": Int64, "uint32": Uint32, "uint64": Uint64,
	"sint32": Sint32, "sint64": Sint64, "fixed32": Fixed32, "fixed64": Fixed64, "sfixed32": SFixed32,
	"sfixed64": SFixed64, "bool": Bool, "string": String, "bytes": Bytes,
}

func (s *Scalar) Parse(lex *lexer.PeekingLexer) error {
	token, err := lex.Peek(0)
	if err != nil {
		return err
	}
	v, ok := stringToScalar[token.Value]
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
	Pos lexer.Position

	Scalar    Scalar   `  @@`
	Map       *MapType `| @@`
	Reference string   `| @(Ident ( "." Ident)*)`
}

type MapType struct {
	Pos lexer.Position

	Key   *Type `"map" "<" @@`
	Value *Type `"," @@ ">"`
}
