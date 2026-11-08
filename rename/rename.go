package rename

import (
	"github.com/tdewolff/parse/v2/js"

	"veiljs/internal/names"
)

func Apply(ast *js.AST, g *names.Gen) {
	protected := map[string]bool{}
	js.Walk(&protector{set: protected}, ast)

	for _, v := range ast.BlockStmt.Scope.Declared {
		if protected[string(v.Name())] {
			continue
		}
		v.Data = []byte(g.Next())
	}
}

type protector struct {
	set map[string]bool
}

func (p *protector) Enter(n js.INode) js.IVisitor {
	switch n := n.(type) {
	case *js.ImportStmt:
		if n.Default != nil {
			p.set[string(n.Default)] = true
		}
		for i := range n.List {
			if n.List[i].Name != nil {
				p.set[string(n.List[i].Name)] = true
			}
			if n.List[i].Binding != nil {
				p.set[string(n.List[i].Binding)] = true
			}
		}
	case *js.ExportStmt:
		for i := range n.List {
			if n.List[i].Name != nil {
				p.set[string(n.List[i].Name)] = true
			}
		}
	}
	return p
}

func (p *protector) Exit(js.INode) {}
