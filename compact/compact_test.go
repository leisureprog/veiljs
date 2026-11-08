package compact

import (
	"testing"

	"github.com/leisureprog/veiljs/internal/jstest"
	"github.com/leisureprog/veiljs/parser"
)

var snippets = []string{
	`(function(a,b){return a+b})(1,2)`,
	`(function(){var o={a:1,b:[1,2,3],c:{d:4}};return JSON.stringify(o)})()`,
	`(function(){class A{constructor(x){this.x=x}get y(){return this.x+1}static z(){return 5}}var a=new A(3);return a.y+A.z()})()`,
	`(function(){var s=0;for(var i=0;i<5;i++){if(i===2)continue;if(i===4)break;s+=i}return s})()`,
	`(function(){try{throw new Error("x")}catch(e){return e.message}})()`,
	`(function(){var a=[1,2,3];return a.map(function(x){return x*2}).join(",")})()`,
	"(function(){var n=\"w\";return `hi ${n}`})()",
	`(function(){var o={a:{b:5}};return o?.a?.b})()`,
	`(function(){var [a,b]=[1,2];var {c}={c:3};return a+b+c})()`,
	`(function(){var a=[1,2];return [0].concat(a,3).join("")})()`,
	`(function(){var x=2;switch(x){case 1:return "a";case 2:return "b";default:return "c"}})()`,
	`(function(){return 1+2*3-4/2})()`,
	`(function(){return "hello"+" "+"world"})()`,
}

func TestCompactSemantics(t *testing.T) {
	for _, src := range snippets {
		ast, err := parser.Parse(src + ";")
		if err != nil {
			t.Fatalf("parse %q: %v", src, err)
		}
		out := Print(ast)
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
			t.Fatalf("compact changed behavior: %q -> %q: %s != %s", src, out, want, got)
		}
	}
}
