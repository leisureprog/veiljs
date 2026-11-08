package names

import (
	"strings"
	"testing"

	"github.com/leisureprog/veiljs/internal/rng"
	"github.com/leisureprog/veiljs/parser"
)

func TestNextUniqueAndPrefixed(t *testing.T) {
	ast, err := parser.Parse("var a=1;function b(){}")
	if err != nil {
		t.Fatal(err)
	}
	g := New(rng.New(rng.Seed([]byte("s"))), ast)
	seen := map[string]bool{}
	for i := 0; i < 200; i++ {
		n := g.Next()
		if !strings.HasPrefix(n, "_0x") {
			t.Fatalf("no prefix: %s", n)
		}
		if seen[n] {
			t.Fatalf("duplicate: %s", n)
		}
		seen[n] = true
	}
}

func TestAvoidsExistingNames(t *testing.T) {
	ast, err := parser.Parse("var _0x0=1;function _0x1(){}")
	if err != nil {
		t.Fatal(err)
	}
	g := New(rng.New(rng.Seed([]byte("s"))), ast)
	for i := 0; i < 200; i++ {
		n := g.Next()
		if n == "_0x0" || n == "_0x1" {
			t.Fatalf("collision with existing name: %s", n)
		}
	}
}
