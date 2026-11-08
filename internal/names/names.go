package names

import (
	"strconv"

	"github.com/tdewolff/parse/v2/js"

	"veiljs/internal/rng"
)

type Gen struct {
	r    *rng.RNG
	used map[string]bool
}

func New(r *rng.RNG, root js.INode) *Gen {
	g := &Gen{r: r, used: map[string]bool{}}
	js.Walk(&collector{set: g.used}, root)
	return g
}

func (g *Gen) Next() string {
	for {
		name := "_0x" + strconv.FormatInt(int64(g.r.Uint64()&0xffffff), 36)
		if !g.used[name] {
			g.used[name] = true
			return name
		}
	}
}

type collector struct {
	set map[string]bool
}

func (c *collector) Enter(n js.INode) js.IVisitor {
	switch n := n.(type) {
	case *js.Var:
		c.set[string(n.Name())] = true
	case *js.Alias:
		if n.Name != nil {
			c.set[string(n.Name)] = true
		}
		if n.Binding != nil {
			c.set[string(n.Binding)] = true
		}
	case *js.ImportStmt:
		if n.Default != nil {
			c.set[string(n.Default)] = true
		}
	case *js.ExportStmt:
		if n.Module != nil {
			c.set[string(n.Module)] = true
		}
	case *js.BranchStmt:
		if n.Label != nil {
			c.set[string(n.Label)] = true
		}
	case *js.LabelledStmt:
		c.set[string(n.Label)] = true
	case *js.PropertyName:
		if n.Computed == nil && n.Literal.TokenType == js.IdentifierToken {
			c.set[string(n.Literal.Data)] = true
		}
	}
	return c
}

func (c *collector) Exit(js.INode) {}
