package rename

import (
	"strings"
	"testing"

	"github.com/leisureprog/veiljs/compact"
	"github.com/leisureprog/veiljs/internal/jstest"
	"github.com/leisureprog/veiljs/internal/names"
	"github.com/leisureprog/veiljs/internal/rng"
	"github.com/leisureprog/veiljs/parser"
)

func TestRenameTopLevel(t *testing.T) {
	src := "function greet(name){return 'hi '+name};greet('bob');"
	ast, err := parser.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	Apply(ast, names.New(rng.New(rng.Seed([]byte("s"))), ast))
	out := compact.Print(ast)
	if strings.Contains(out, "greet") {
		t.Fatalf("name not renamed: %s", out)
	}
	if !strings.Contains(out, "_0x") {
		t.Fatalf("no new names: %s", out)
	}
	want, err := jstest.Eval(src)
	if err != nil {
		t.Fatalf("eval %q: %v", src, err)
	}
	got, err := jstest.Eval(out)
	if err != nil {
		t.Fatalf("eval %q: %v", out, err)
	}
	if want != got {
		t.Fatalf("behavior changed: %q -> %q: %s != %s", src, out, want, got)
	}
}

func TestProtectImports(t *testing.T) {
	src := "import foo from 'mod';var bar=1;"
	ast, err := parser.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	Apply(ast, names.New(rng.New(rng.Seed([]byte("s"))), ast))
	out := compact.Print(ast)
	if !strings.Contains(out, "foo") {
		t.Fatalf("import renamed: %s", out)
	}
	if strings.Contains(out, "bar") {
		t.Fatalf("top-level not renamed: %s", out)
	}
}
