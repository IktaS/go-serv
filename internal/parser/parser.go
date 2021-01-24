package parser

// Parser define what a parser should have
type Parser interface {
	Parse([]byte) interface{}
}
