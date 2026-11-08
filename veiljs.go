package veiljs

import (
	"strconv"

	"veiljs/cff"
	"veiljs/compact"
	"veiljs/internal/names"
	"veiljs/internal/rng"
	"veiljs/numbers"
	"veiljs/parser"
	"veiljs/rename"
	"veiljs/simplify"
	"veiljs/stringarray"
	"veiljs/unicode"
)

func Obfuscate(data []byte) (string, error) {
	return ObfuscateWithConfig(data, DefaultConfig())
}

func ObfuscateJS(code string) string {
	out, err := Obfuscate([]byte(code))
	if err != nil {
		return code
	}
	return out
}

func ObfuscateWithConfig(data []byte, cfg Config) (string, error) {
	if len(data) == 0 {
		return "", nil
	}
	code := string(data)
	ast, err := parser.Parse(code)
	if err != nil {
		return code, nil
	}

	r := rng.New(rng.Seed(data, configSeed(cfg)))
	g := names.New(r, ast)

	if cfg.StringArray {
		stringarray.Apply(ast, r, g, stringarray.Options{
			Encoding:  pickEncoding(r, cfg.StringArrayEncoding),
			Threshold: cfg.StringArrayThreshold,
		})
	}
	if cfg.ControlFlowFlattening {
		cff.Apply(ast, r, g, cfg.ControlFlowFlatteningThreshold)
	}
	if cfg.Simplify {
		simplify.Apply(ast)
	}
	if cfg.NumbersToExpressions {
		numbers.Apply(ast, r)
	}
	if cfg.UnicodeEscapeSequence {
		unicode.Apply(ast)
	}
	if cfg.RenameGlobals {
		rename.Apply(ast, g)
	}

	if cfg.Compact {
		return compact.Print(ast), nil
	}
	return parser.Print(ast), nil
}

func pickEncoding(r *rng.RNG, list []string) string {
	if len(list) == 0 {
		return "none"
	}
	return list[r.Intn(len(list))]
}

func configSeed(cfg Config) []byte {
	var b []byte
	for _, f := range []bool{
		cfg.Compact,
		cfg.ControlFlowFlattening,
		cfg.NumbersToExpressions,
		cfg.Simplify,
		cfg.StringArray,
		cfg.UnicodeEscapeSequence,
		cfg.RenameGlobals,
	} {
		if f {
			b = append(b, 1)
		} else {
			b = append(b, 0)
		}
	}
	b = append(b, []byte(strconv.FormatFloat(cfg.ControlFlowFlatteningThreshold, 'g', -1, 64))...)
	b = append(b, 0)
	b = append(b, []byte(strconv.FormatFloat(cfg.StringArrayThreshold, 'g', -1, 64))...)
	b = append(b, 0)
	for _, e := range cfg.StringArrayEncoding {
		b = append(b, []byte(e)...)
		b = append(b, 0)
	}
	return b
}
