package simplify

import (
	"testing"

	"veiljs/compact"
	"veiljs/internal/jstest"
	"veiljs/parser"
)

var exprs = []string{
	"1+2", "2*3+4", "10-3", "2*3*4",
	"-5+2", "3*-4", "-2147483648+1", "2147483647+1",
	"17%5", "-17%5", "5%0",
	"5&3", "5|2", "5^1", "1<<4", "16>>2", "16>>>2", "1<<32",
	"3===3", "3!==4", "2<3", "3<=3", "4>3", "3>=4",
	"!true", "!false", "!0", "!5", "!(-3)",
	"true&&5", "false||7", "false&&9", "true||3",
	`typeof "x"`,
}

func TestSimplifyPreservesValue(t *testing.T) {
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
			t.Fatalf("value changed: %s -> %q: %s != %s", e, out, want, got)
		}
	}
}
