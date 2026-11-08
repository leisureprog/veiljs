package jsstr

import "testing"

func TestUnquote(t *testing.T) {
	cases := map[string]string{
		`"abc"`:          "abc",
		`"a\nb"`:         "a\nb",
		`'q"q'`:          `q"q`,
		`"\x41"`:         "A",
		`"\u0041"`:       "A",
		`"\u{1F600}"`:    "\U0001F600",
		`"\uD83D\uDE00"`: "\U0001F600",
	}
	for in, want := range cases {
		if got := Unquote([]byte(in)); got != want {
			t.Fatalf("Unquote(%s)=%q, want %q", in, got, want)
		}
	}
}

func TestQuoteUnquoteRoundtrip(t *testing.T) {
	for _, s := range []string{"", "abc", "a\nb", `q"q`, "\\", "\x00\x1f", "\U0001F600", "world"} {
		if got := Unquote(Quote(s)); got != s {
			t.Fatalf("roundtrip %q -> %q", s, got)
		}
	}
}

func TestEscapeUnicodeRoundtrip(t *testing.T) {
	for _, s := range []string{"abc", "world", "\U0001F600", "a\nb\"c"} {
		if got := Unquote(EscapeUnicode(Quote(s))); got != s {
			t.Fatalf("escape roundtrip %q -> %q", s, got)
		}
	}
}
