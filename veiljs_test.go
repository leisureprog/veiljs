package veiljs

import (
	"strings"
	"sync"
	"testing"

	"veiljs/internal/jstest"
)

var programs = []string{
	`(function(a,b){return a+b})(2,3)`,
	`(function(){var s=0;for(var i=0;i<10;i++){if(i%2===0)continue;s+=i}return s})()`,
	`(function(){class A{constructor(x){this.x=x}add(y){return this.x+y}}return new A(4).add(6)})()`,
	`(function(){var o={a:1,b:2};var s=0;for(var k in o){s+=o[k]}return s})()`,
	`(function(){try{throw new Error("boom")}catch(e){return e.message}})()`,
	`(function(){var a=[3,1,2];return a.sort().join("")})()`,
}

const bundle = `
(function(){
  function fib(n){return n<2?n:fib(n-1)+fib(n-2)}
  function sum(a){var s=0;for(var i=0;i<a.length;i++){s+=a[i]}return s}
  var lib={name:"demo",version:"1.0.0"};
  class Counter{constructor(){this.n=0}inc(){this.n++;return this.n}get value(){return this.n}}
  var c=new Counter();c.inc();c.inc();c.inc();
  return JSON.stringify({lib:lib.name,fib:fib(10),sum:sum([1,2,3,4]),count:c.value});
})()
`

func TestObfuscatePreservesBehavior(t *testing.T) {
	for _, src := range programs {
		out, err := Obfuscate([]byte(src))
		if err != nil {
			t.Fatalf("obfuscate %q: %v", src, err)
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
}

func TestConfigVariants(t *testing.T) {
	cfgs := []Config{
		DefaultConfig(),
		{Compact: true, StringArray: true, StringArrayEncoding: []string{"hex", "base64", "rc4"}, StringArrayThreshold: 1, ControlFlowFlattening: true, ControlFlowFlatteningThreshold: 1, NumbersToExpressions: true, Simplify: true, UnicodeEscapeSequence: true, RenameGlobals: true},
	}
	for i, cfg := range cfgs {
		for _, src := range programs {
			out, err := ObfuscateWithConfig([]byte(src), cfg)
			if err != nil {
				t.Fatalf("cfg %d obfuscate %q: %v", i, src, err)
			}
			want, err := jstest.Eval(src)
			if err != nil {
				t.Fatalf("cfg %d eval %q: %v", i, src, err)
			}
			got, err := jstest.Eval(out)
			if err != nil {
				t.Fatalf("cfg %d eval %q: %v", i, out, err)
			}
			if want != got {
				t.Fatalf("cfg %d behavior changed: %q -> %q: %s != %s", i, src, out, want, got)
			}
		}
	}
}

func TestIdempotent(t *testing.T) {
	for _, src := range programs {
		a, err := Obfuscate([]byte(src))
		if err != nil {
			t.Fatal(err)
		}
		b, err := Obfuscate([]byte(src))
		if err != nil {
			t.Fatal(err)
		}
		if a != b {
			t.Fatalf("not deterministic:\n%s\n%s", a, b)
		}
	}
}

func TestEmptyInput(t *testing.T) {
	out, err := Obfuscate(nil)
	if err != nil {
		t.Fatal(err)
	}
	if out != "" {
		t.Fatalf("expected empty string, got %q", out)
	}
}

func TestInvalidInputFallback(t *testing.T) {
	src := "function("
	out, err := Obfuscate([]byte(src))
	if err != nil {
		t.Fatalf("error must not be returned: %v", err)
	}
	if out != src {
		t.Fatalf("expected original code, got %q", out)
	}
}

func TestObfuscateJS(t *testing.T) {
	src := `(function(){return "ok"})()`
	if got := ObfuscateJS(src); got == "" {
		t.Fatal("empty result")
	}
}

func TestConcurrent(t *testing.T) {
	src := `(function(){var s=0;for(var i=0;i<20;i++){s+=i}return s})()`
	want, err := Obfuscate([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	results := make([]string, 16)
	for i := range results {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i], _ = Obfuscate([]byte(src))
		}(i)
	}
	wg.Wait()
	for i, got := range results {
		if got != want {
			t.Fatalf("goroutine %d produced a different result", i)
		}
	}
}

func TestBundle(t *testing.T) {
	out, err := Obfuscate([]byte(bundle))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "_0x") {
		t.Fatalf("bundle not obfuscated: %s", out)
	}
	want, err := jstest.Eval(bundle)
	if err != nil {
		t.Fatalf("eval bundle: %v", err)
	}
	got, err := jstest.Eval(out)
	if err != nil {
		t.Fatalf("eval obfuscated: %v", err)
	}
	if want != got {
		t.Fatalf("bundle behavior changed: %s != %s", want, got)
	}
}
