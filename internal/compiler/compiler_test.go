package compiler

import (
	"fmt"
	"testing"

	"github.com/go-serv/internal/compiler/ast"
	"github.com/go-serv/internal/compiler/lexer"
	"github.com/go-serv/internal/compiler/parser"
)

func TestWorld(t *testing.T) {
	input := []byte(`
		service TestService ( TestMessage1 ) : TestMessage1;
	`)
	lex := lexer.NewLexer(input)
	p := parser.NewParser()
	st, err := p.Parse(lex)
	if err != nil {
		panic(err)
	}
	w, ok := st.(*ast.Program)
	if !ok {
		t.Fatalf("This is not a world")
	}
	get := fmt.Sprintf("%v", w)
	if get != `[{TestService}]` {
		t.Fatalf("Wrong world %v", w)
	}
}
