package numbers

import (
	"strconv"

	"github.com/tdewolff/parse/v2/js"

	"github.com/leisureprog/veiljs/internal/rng"
	"github.com/leisureprog/veiljs/parser"
)

func Apply(ast *js.AST, r *rng.RNG) {
	parser.Rewrite(ast, func(n js.INode) js.INode {
		lit, ok := n.(*js.LiteralExpr)
		if !ok {
			return n
		}
		v, ok := value(lit)
		if !ok {
			return n
		}
		if r.Intn(2) == 0 {
			a := int32(r.Uint64() & 0x7fffffff)
			b := a ^ v
			return &js.GroupExpr{X: &js.BinaryExpr{
				Op: js.BitXorToken,
				X:  &js.LiteralExpr{TokenType: js.HexadecimalToken, Data: []byte("0x" + strconv.FormatInt(int64(a), 16))},
				Y:  &js.LiteralExpr{TokenType: js.HexadecimalToken, Data: []byte("0x" + strconv.FormatInt(int64(b), 16))},
			}}
		}
		a := int32(r.Intn(int(v) + 1))
		b := v - a
		return &js.GroupExpr{X: &js.BinaryExpr{
			Op: js.AddToken,
			X:  &js.LiteralExpr{TokenType: js.IntegerToken, Data: []byte(strconv.FormatInt(int64(a), 10))},
			Y:  &js.LiteralExpr{TokenType: js.IntegerToken, Data: []byte(strconv.FormatInt(int64(b), 10))},
		}}
	})
}

func value(lit *js.LiteralExpr) (int32, bool) {
	s := string(lit.Data)
	base := 10
	digits := s
	switch lit.TokenType {
	case js.HexadecimalToken:
		base = 16
		digits = s[2:]
	case js.BinaryToken:
		base = 2
		digits = s[2:]
	case js.OctalToken:
		base = 8
		digits = s[2:]
	case js.IntegerToken:
	default:
		return 0, false
	}
	u, err := strconv.ParseUint(digits, base, 32)
	if err != nil || u > 0x7fffffff {
		return 0, false
	}
	return int32(u), true
}
