# veiljs

A native JavaScript obfuscator written in pure Go. It mirrors the behavior of
`javascript-obfuscator`, but runs without Node.js, CGO, `npx`, external binaries,
or `exec.Command` — only Go and the standard library.

Author: **Leisure**

## Features

- **String Array** — move strings into an array with `none`, `base64`, `hex`, `rc4` encoding.
- **Control Flow Flattening** — flatten control flow through a switch dispatcher.
- **Numbers To Expressions** — replace numeric literals with arithmetic expressions.
- **Simplify** — constant folding and dead code elimination.
- **Unicode Escape Sequence** — escape strings as `\uXXXX`.
- **Rename Globals** — rename global identifiers.
- **Compact** — minified code generator.

Transformations are applied in order: String Array → Control Flow Flattening →
Numbers To Expressions → Simplify → Unicode Escape Sequence → Rename Globals →
Compact.

## Installation

```bash
go get github.com/leisureprog/veiljs
```

Dependencies: [`github.com/tdewolff/parse/v2`](https://github.com/tdewolff/parse/v2) and the standard library.

## Usage

```go
package main

import (
	"fmt"

	"github.com/leisureprog/veiljs"
)

func main() {
	out, err := veiljs.Obfuscate([]byte(`console.log("hello world");`))
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}
```

With a custom configuration:

```go
cfg := veiljs.DefaultConfig()
cfg.StringArrayEncoding = []string{"base64", "rc4"}
cfg.ControlFlowFlatteningThreshold = 0.5

out, err := veiljs.ObfuscateWithConfig([]byte(code), cfg)
```

Without error handling:

```go
out := veiljs.ObfuscateJS(code)
```

## API

```go
func Obfuscate(data []byte) (string, error)
func ObfuscateWithConfig(data []byte, cfg Config) (string, error)
func ObfuscateJS(code string) string
func DefaultConfig() Config
```

Error behavior:

- empty input → `("", nil)`;
- parse or transformation failure → the original code is returned.

## Configuration

```go
type Config struct {
	Compact                        bool
	ControlFlowFlattening          bool
	ControlFlowFlatteningThreshold float64
	NumbersToExpressions           bool
	Simplify                       bool
	StringArray                    bool
	StringArrayEncoding            []string
	StringArrayThreshold           float64
	UnicodeEscapeSequence          bool
	RenameGlobals                  bool
}
```

| Field                            | Default      | Description                                     |
| -------------------------------- | ------------ | ----------------------------------------------- |
| `Compact`                        | `true`       | Minified output                                 |
| `ControlFlowFlattening`          | `true`       | Control flow flattening                         |
| `ControlFlowFlatteningThreshold` | `0.75`       | Share of nodes affected by CFF                  |
| `NumbersToExpressions`           | `false`      | Numbers to expressions                          |
| `Simplify`                       | `true`       | Constant folding and dead code elimination      |
| `StringArray`                    | `true`       | Move strings into an array                      |
| `StringArrayEncoding`            | `["base64"]` | String encoding: `none`, `base64`, `hex`, `rc4` |
| `StringArrayThreshold`           | `0.75`       | Share of strings moved into the array           |
| `UnicodeEscapeSequence`          | `false`      | Escape strings as `\uXXXX`                      |
| `RenameGlobals`                  | `false`      | Rename global identifiers                       |

## Guarantees

- **Determinism** — identical input and configuration produce identical output
  (PRNG seeded from the input data and configuration).
- **Thread safety** — calls are independent and share no state.
- **No external processes** — pure Go, no CGO and no Node.js.

## Testing

```bash
go build ./...
go vet ./...
go test ./...
go test -race ./...
```

Tests verify the semantics of each transformation: obfuscated code is executed in
[`github.com/dop251/goja`](https://github.com/dop251/goja) and compared against the
original result, plus idempotency and behavior on a real bundle.

## Layout

```
veiljs
├── config.go           public Config and DefaultConfig
├── veiljs.go           transformation orchestration and public API
├── parser/             AST parsing and printing
├── compact/            minified code generator
├── stringarray/        string array and encoders
├── cff/                control flow flattening
├── numbers/            numbers to expressions
├── simplify/           constant folding and dead code elimination
├── unicode/            unicode escape sequences
├── rename/             global renaming
└── internal/           rng, names, jsstr, jstest
```

## License

MIT © Leisure
