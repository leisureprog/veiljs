package cff

import (
	"strconv"

	"github.com/tdewolff/parse/v2/js"

	"veiljs/internal/names"
	"veiljs/internal/rng"
	"veiljs/parser"
)

func Apply(ast *js.AST, r *rng.RNG, g *names.Gen, threshold float64) {
	parser.Rewrite(ast, func(n js.INode) js.INode {
		switch n := n.(type) {
		case *js.FuncDecl:
			flatten(&n.Body, r, g, threshold)
		case *js.MethodDecl:
			flatten(&n.Body, r, g, threshold)
		case *js.ArrowFunc:
			flatten(&n.Body, r, g, threshold)
		}
		return n
	})
}

func flatten(b *js.BlockStmt, r *rng.RNG, g *names.Gen, threshold float64) {
	if len(b.List) < 3 || r.Float64() >= threshold {
		return
	}

	var hoisted []js.BindingElement
	body := make([]js.IStmt, 0, len(b.List))
	for _, s := range b.List {
		switch s := s.(type) {
		case *js.DirectivePrologueStmt, *js.FuncDecl, *js.ClassDecl:
			return
		case *js.VarDecl:
			if s.TokenType == js.VarToken {
				body = append(body, s)
				continue
			}
			for i := range s.List {
				el := &s.List[i]
				v, ok := el.Binding.(*js.Var)
				if !ok {
					return
				}
				hoisted = append(hoisted, js.BindingElement{Binding: &js.Var{Data: v.Name()}})
				if el.Default != nil {
					body = append(body, &js.ExprStmt{Value: &js.BinaryExpr{Op: js.EqToken, X: &js.Var{Data: v.Name()}, Y: el.Default}})
				}
			}
		default:
			body = append(body, s)
		}
	}
	if len(body) < 3 {
		return
	}

	state := g.Next()
	cases := make([]js.CaseClause, len(body))
	for i, s := range body {
		step := []js.IStmt{s}
		if !terminates(s) {
			step = append(step, &js.ExprStmt{Value: assign(state, i+1)}, &js.BranchStmt{Type: js.ContinueToken})
		}
		cases[i] = js.CaseClause{TokenType: js.CaseToken, Cond: intLit(i), List: step}
	}

	prelude := &js.VarDecl{
		TokenType: js.VarToken,
		List:      append(hoisted, js.BindingElement{Binding: &js.Var{Data: []byte(state)}, Default: intLit(0)}),
	}
	b.List = []js.IStmt{
		prelude,
		&js.WhileStmt{
			Cond: &js.LiteralExpr{TokenType: js.TrueToken, Data: []byte("true")},
			Body: &js.BlockStmt{List: []js.IStmt{
				&js.SwitchStmt{Init: &js.Var{Data: []byte(state)}, List: cases},
				&js.BranchStmt{Type: js.BreakToken},
			}},
		},
	}
}

func terminates(s js.IStmt) bool {
	switch s.(type) {
	case *js.ReturnStmt, *js.ThrowStmt:
		return true
	}
	return false
}

func assign(name string, v int) js.IExpr {
	return &js.BinaryExpr{Op: js.EqToken, X: &js.Var{Data: []byte(name)}, Y: intLit(v)}
}

func intLit(v int) *js.LiteralExpr {
	return &js.LiteralExpr{TokenType: js.IntegerToken, Data: []byte(strconv.Itoa(v))}
}
