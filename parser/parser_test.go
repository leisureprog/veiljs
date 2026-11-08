package parser

import (
	"strings"
	"testing"

	"github.com/tdewolff/parse/v2/js"
)

func TestParsePrint(t *testing.T) {
	code := "function f(a){return a+1}f(2);"
	ast, err := Parse(code)
	if err != nil {
		t.Fatal(err)
	}
	out := Print(ast)
	if _, err := Parse(out); err != nil {
		t.Fatalf("re-parse %q: %v", out, err)
	}
}

func TestParseError(t *testing.T) {
	if _, err := Parse("function("); err == nil {
		t.Fatal("expected a parse error")
	}
}

func TestRewrite(t *testing.T) {
	ast, err := Parse("var x=1;")
	if err != nil {
		t.Fatal(err)
	}
	Rewrite(ast, func(n js.INode) js.INode {
		if lit, ok := n.(*js.LiteralExpr); ok && string(lit.Data) == "1" {
			lit.Data = []byte("2")
		}
		return n
	})
	out := Print(ast)
	if _, err := Parse(out); err != nil {
		t.Fatalf("re-parse %q: %v", out, err)
	}
	if !strings.Contains(out, "2") {
		t.Fatalf("replacement not applied: %s", out)
	}
}
