package stringarray

import (
	"strings"
	"testing"

	"github.com/leisureprog/veiljs/compact"
	"github.com/leisureprog/veiljs/internal/jstest"
	"github.com/leisureprog/veiljs/internal/names"
	"github.com/leisureprog/veiljs/internal/rng"
	"github.com/leisureprog/veiljs/parser"
)

func TestEncodingsPreserveStrings(t *testing.T) {
	src := `(function(){return 'hello world 😀'})()`
	for _, enc := range []string{"none", "base64", "hex", "rc4"} {
		ast, err := parser.Parse(src + ";")
		if err != nil {
			t.Fatalf("%s: %v", enc, err)
		}
		r := rng.New(rng.Seed([]byte(enc)))
		Apply(ast, r, names.New(r, ast), Options{Encoding: enc, Threshold: 1})
		out := compact.Print(ast)
		want, err := jstest.Eval(src)
		if err != nil {
			t.Fatalf("%s eval %q: %v", enc, src, err)
		}
		got, err := jstest.Eval(out)
		if err != nil {
			t.Fatalf("%s eval %q: %v", enc, out, err)
		}
		if want != got {
			t.Fatalf("%s: string changed: %q: %s != %s", enc, out, want, got)
		}
	}
}

func TestSkipsModuleArgs(t *testing.T) {
	ast, err := parser.Parse("var x=require('fs');")
	if err != nil {
		t.Fatal(err)
	}
	r := rng.New(rng.Seed([]byte("s")))
	Apply(ast, r, names.New(r, ast), Options{Encoding: "base64", Threshold: 1})
	out := compact.Print(ast)
	if !strings.Contains(out, "fs") {
		t.Fatalf("require argument missing: %s", out)
	}
	if strings.Contains(out, "_0x") {
		t.Fatalf("require argument obfuscated: %s", out)
	}
}
