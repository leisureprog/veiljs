package unicode

import (
	"github.com/tdewolff/parse/v2/js"

	"veiljs/internal/jsstr"
	"veiljs/parser"
)

func Apply(ast *js.AST) {
	parser.Rewrite(ast, func(n js.INode) js.INode {
		lit, ok := n.(*js.LiteralExpr)
		if !ok || lit.TokenType != js.StringToken || len(lit.Data) == 0 {
			return n
		}
		lit.Data = jsstr.EscapeUnicode(lit.Data)
		return lit
	})
}
