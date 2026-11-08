package cff

import (
	"testing"

	"github.com/leisureprog/veiljs/compact"
	"github.com/leisureprog/veiljs/internal/jstest"
	"github.com/leisureprog/veiljs/internal/names"
	"github.com/leisureprog/veiljs/internal/rng"
	"github.com/leisureprog/veiljs/parser"
)

var progs = []string{
	`(function(x){var acc=0;if(x>0){acc+=x}else{acc-=x}acc*=2;return acc})(5)`,
	`(function(){var r=0;for(var i=0;i<6;i++){if(i===3)continue;if(i===5)break;r+=i}return r})()`,
	`(function(){var out=[];for(var i=0;i<3;i++){out.push(function(){return i*2}())}return out.join(",")})()`,
	`(function(){var m=0;for(var i=0;i<4;i++){try{if(i===2)throw new Error("e");m+=i}catch(err){m+=100}}return m})()`,
	`(function(){const a=2;let b=3;var s=[];s.push(a+b);s.push(a*b);return s.join(",")})()`,
	`(function(){var fs=[];for(var i=1;i<=3;i++){fs.push(function(){return i})}return fs.map(function(f){return f()}).join("")})()`,
}

func TestCFFPreservesBehavior(t *testing.T) {
	for _, src := range progs {
		ast, err := parser.Parse(src + ";")
		if err != nil {
			t.Fatalf("parse %q: %v", src, err)
		}
		r := rng.New(rng.Seed([]byte(src)))
		Apply(ast, r, names.New(r, ast), 1)
		out := compact.Print(ast)
		if _, err := parser.Parse(out); err != nil {
			t.Fatalf("reparse %q -> %q: %v", src, out, err)
		}
		want, err := jstest.Expr(src)
		if err != nil {
			t.Fatalf("eval %q: %v", src, err)
		}
		got, err := jstest.Expr(out)
		if err != nil {
			t.Fatalf("eval %q: %v", out, err)
		}
		if want != got {
			t.Fatalf("behavior changed: %q -> %q: %s != %s", src, out, want, got)
		}
	}
}
