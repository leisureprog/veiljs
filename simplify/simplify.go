package simplify

import (
	"math"
	"strconv"
	"strings"

	"github.com/tdewolff/parse/v2/js"

	"veiljs/parser"
)

func Apply(ast *js.AST) {
	parser.Rewrite(ast, func(n js.INode) js.INode {
		switch n := n.(type) {
		case *js.BinaryExpr:
			return foldBinary(n)
		case *js.UnaryExpr:
			return foldUnary(n)
		}
		return n
	})
}

func foldBinary(n *js.BinaryExpr) js.INode {
	if n.Op == js.AndToken || n.Op == js.OrToken {
		return foldLogic(n)
	}
	x, okx := asInt(n.X)
	y, oky := asInt(n.Y)
	if !okx || !oky {
		return n
	}
	if v, ok := arith(n.Op, x, y); ok {
		return intLit(v)
	}
	if v, ok := bits(n.Op, x, y); ok {
		return intLit(v)
	}
	if v, ok := compare(n.Op, x, y); ok {
		return boolLit(v)
	}
	return n
}

func foldLogic(n *js.BinaryExpr) js.INode {
	b, ok := asBool(n.X)
	if !ok {
		return n
	}
	if n.Op == js.AndToken {
		if b {
			return n.Y
		}
		return boolLit(false)
	}
	if b {
		return boolLit(true)
	}
	return n.Y
}

func foldUnary(n *js.UnaryExpr) js.INode {
	switch n.Op {
	case js.NotToken:
		if b, ok := asBool(n.X); ok {
			return boolLit(!b)
		}
		if v, ok := asInt(n.X); ok {
			return boolLit(v == 0)
		}
	case js.TypeofToken:
		if lit, ok := n.X.(*js.LiteralExpr); ok && lit.TokenType == js.StringToken {
			return &js.LiteralExpr{TokenType: js.StringToken, Data: []byte(`"string"`)}
		}
	}
	return n
}

func arith(op js.TokenType, x, y int32) (int32, bool) {
	var v int64
	switch op {
	case js.AddToken:
		v = int64(x) + int64(y)
	case js.SubToken:
		v = int64(x) - int64(y)
	case js.MulToken:
		v = int64(x) * int64(y)
	default:
		return 0, false
	}
	if v < math.MinInt32 || v > math.MaxInt32 {
		return 0, false
	}
	return int32(v), true
}

func bits(op js.TokenType, x, y int32) (int32, bool) {
	switch op {
	case js.BitAndToken:
		return x & y, true
	case js.BitOrToken:
		return x | y, true
	case js.BitXorToken:
		return x ^ y, true
	case js.ModToken:
		if y != 0 {
			return x % y, true
		}
	case js.LtLtToken:
		if uint(y) < 32 {
			return x << y, true
		}
	case js.GtGtToken:
		if uint(y) < 32 {
			return x >> y, true
		}
	case js.GtGtGtToken:
		if uint(y) < 32 {
			return int32(uint32(x) >> uint(y)), true
		}
	}
	return 0, false
}

func compare(op js.TokenType, x, y int32) (bool, bool) {
	switch op {
	case js.EqEqEqToken:
		return x == y, true
	case js.NotEqEqToken:
		return x != y, true
	case js.LtToken:
		return x < y, true
	case js.LtEqToken:
		return x <= y, true
	case js.GtToken:
		return x > y, true
	case js.GtEqToken:
		return x >= y, true
	}
	return false, false
}

func asInt(e js.IExpr) (int32, bool) {
	neg := false
	if u, ok := e.(*js.UnaryExpr); ok && u.Op == js.NegToken {
		e, neg = u.X, true
	}
	lit, ok := e.(*js.LiteralExpr)
	if !ok {
		return 0, false
	}
	s := string(lit.Data)
	base := 10
	switch lit.TokenType {
	case js.HexadecimalToken:
		s, base = s[2:], 16
	case js.BinaryToken:
		s, base = s[2:], 2
	case js.OctalToken:
		s, base = s[2:], 8
	case js.IntegerToken:
	default:
		return 0, false
	}
	v, err := strconv.ParseInt(strings.ReplaceAll(s, "_", ""), base, 64)
	if err != nil {
		return 0, false
	}
	if neg {
		v = -v
	}
	if v < math.MinInt32 || v > math.MaxInt32 {
		return 0, false
	}
	return int32(v), true
}

func asBool(e js.IExpr) (bool, bool) {
	lit, ok := e.(*js.LiteralExpr)
	if !ok {
		return false, false
	}
	switch lit.TokenType {
	case js.TrueToken:
		return true, true
	case js.FalseToken:
		return false, true
	}
	return false, false
}

func boolLit(v bool) *js.LiteralExpr {
	if v {
		return &js.LiteralExpr{TokenType: js.TrueToken, Data: []byte("true")}
	}
	return &js.LiteralExpr{TokenType: js.FalseToken, Data: []byte("false")}
}

func intLit(v int32) js.IExpr {
	if v < 0 {
		return &js.UnaryExpr{
			Op: js.NegToken,
			X:  &js.LiteralExpr{TokenType: js.IntegerToken, Data: []byte(strconv.FormatInt(-int64(v), 10))},
		}
	}
	return &js.LiteralExpr{TokenType: js.IntegerToken, Data: []byte(strconv.FormatInt(int64(v), 10))}
}
