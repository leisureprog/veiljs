package numbers

import (
	"testing"

	"github.com/leisureprog/veiljs/compact"
	"github.com/leisureprog/veiljs/internal/jstest"
	"github.com/leisureprog/veiljs/internal/rng"
	"github.com/leisureprog/veiljs/parser"
)

var exprs = []string{
	"0", "1", "2", "7", "42", "255", "256", "1000", "65535", "123456", "2147483647",
	"0xff", "0b1010", "0o17", "0x0",
}

func TestNumbersPreserveValue(t *testing.T) {
	for _, e := range exprs {
		ast, err := parser.Parse(e + ";")
		if err != nil {
			t.Fatalf("parse %s: %v", e, err)
		}
		Apply(ast, rng.New(rng.Seed([]byte(e))))
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
			t.Fatalf("value changed: %s -> %q: %s != %s", e, out, want, got)
		}
	}
}
