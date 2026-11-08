package unicode

import (
	"testing"

	"veiljs/compact"
	"veiljs/internal/jstest"
	"veiljs/parser"
)

var exprs = []string{
	`"hello"`,
	`"world"`,
	`"golang"`,
	`"😀"`,
	`"a\nb"`,
	`"quote\"here"`,
}

func TestUnicodePreservesString(t *testing.T) {
	for _, e := range exprs {
		ast, err := parser.Parse(e + ";")
		if err != nil {
			t.Fatalf("parse %s: %v", e, err)
		}
		Apply(ast)
		out := compact.Print(ast)
		want, err := jstest.Expr(e)
		if err != nil {
			t.Fatalf("eval %s: %v", e, err)
		}
		got, err := jstest.Expr(out)
		if err != nil {
			t.Fatalf("eval %q: %v", out, err)
		}
		if want != got {
			t.Fatalf("string changed: %s -> %q: %s != %s", e, out, want, got)
		}
	}
}
