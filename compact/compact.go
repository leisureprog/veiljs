package compact

import (
	"strings"

	"github.com/tdewolff/parse/v2/js"
)

type printer struct {
	sb   strings.Builder
	last byte
}

func identByte(c byte) bool {
	return c == '$' || c == '_' || ('0' <= c && c <= '9') || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c >= 0x80
}

func needSpace(prev, next byte) bool {
	if identByte(prev) && identByte(next) {
		return true
	}
	if prev == '+' && next == '+' {
		return true
	}
	if prev == '-' && next == '-' {
		return true
	}
	if prev == '/' && (next == '/' || next == '*') {
		return true
	}
	if prev == '.' && '0' <= next && next <= '9' {
		return true
	}
	return false
}

func (p *printer) w(b []byte) {
	if len(b) == 0 {
		return
	}
	if p.last != 0 && needSpace(p.last, b[0]) {
		p.sb.WriteByte(' ')
	}
	p.sb.Write(b)
	p.last = b[len(b)-1]
}

func (p *printer) c(ch byte) {
	if p.last != 0 && needSpace(p.last, ch) {
		p.sb.WriteByte(' ')
	}
	p.sb.WriteByte(ch)
	p.last = ch
}

func Print(ast *js.AST) string {
	p := &printer{}
	p.stmts(ast.List)
	return p.sb.String()
}

func (p *printer) stmts(list []js.IStmt) {
	for _, s := range list {
		p.stmt(s)
		if p.last != ';' {
			p.c(';')
		}
	}
}

func (p *printer) stmtBody(s js.IStmt) {
	p.stmt(s)
	if _, ok := s.(*js.BlockStmt); !ok {
		if p.last != ';' {
			p.c(';')
		}
	}
}

func (p *printer) block(b *js.BlockStmt) {
	if b == nil {
		p.w([]byte("{}"))
		return
	}
	if len(b.List) == 0 {
		p.w([]byte("{}"))
		return
	}
	p.c('{')
	p.stmts(b.List)
	p.c('}')
}

func (p *printer) stmt(n js.IStmt) {
	switch n := n.(type) {
	case *js.BlockStmt:
		p.block(n)
	case *js.EmptyStmt:
		p.c(';')
	case *js.Comment:
		p.w(n.Value)
	case *js.DirectivePrologueStmt:
		p.w(n.Value)
		p.c(';')
	case *js.ExprStmt:
		p.expr(n.Value)
		p.c(';')
	case *js.VarDecl:
		p.varDecl(n)
		p.c(';')
	case *js.FuncDecl:
		p.funcDecl(n)
	case *js.ClassDecl:
		p.classDecl(n)
	case *js.IfStmt:
		p.w([]byte("if"))
		p.c('(')
		p.expr(n.Cond)
		p.c(')')
		p.stmtBody(n.Body)
		if n.Else != nil {
			p.w([]byte("else"))
			if _, ok := n.Else.(*js.IfStmt); ok {
				p.stmt(n.Else)
			} else {
				p.stmtBody(n.Else)
			}
		}
	case *js.DoWhileStmt:
		p.w([]byte("do"))
		p.stmtBody(n.Body)
		p.w([]byte("while"))
		p.c('(')
		p.expr(n.Cond)
		p.c(')')
		p.c(';')
	case *js.WhileStmt:
		p.w([]byte("while"))
		p.c('(')
		p.expr(n.Cond)
		p.c(')')
		p.stmtBody(n.Body)
	case *js.ForStmt:
		p.w([]byte("for"))
		p.c('(')
		if n.Init != nil {
			p.expr(n.Init)
		}
		p.c(';')
		if n.Cond != nil {
			p.expr(n.Cond)
		}
		p.c(';')
		if n.Post != nil {
			p.expr(n.Post)
		}
		p.c(')')
		p.block(n.Body)
	case *js.ForInStmt:
		p.w([]byte("for"))
		p.c('(')
		p.expr(n.Init)
		p.w([]byte("in"))
		p.expr(n.Value)
		p.c(')')
		p.block(n.Body)
	case *js.ForOfStmt:
		p.w([]byte("for"))
		if n.Await {
			p.w([]byte("await"))
		}
		p.c('(')
		p.expr(n.Init)
		p.w([]byte("of"))
		p.expr(n.Value)
		p.c(')')
		p.block(n.Body)
	case *js.SwitchStmt:
		p.w([]byte("switch"))
		p.c('(')
		p.expr(n.Init)
		p.c(')')
		p.c('{')
		for i := range n.List {
			p.caseClause(&n.List[i])
		}
		p.c('}')
	case *js.BranchStmt:
		p.w(n.Type.Bytes())
		if n.Label != nil {
			p.w(n.Label)
		}
		p.c(';')
	case *js.ReturnStmt:
		p.w([]byte("return"))
		if n.Value != nil {
			p.expr(n.Value)
		}
		p.c(';')
	case *js.WithStmt:
		p.w([]byte("with"))
		p.c('(')
		p.expr(n.Cond)
		p.c(')')
		p.stmtBody(n.Body)
	case *js.LabelledStmt:
		p.w(n.Label)
		p.c(':')
		p.stmt(n.Value)
	case *js.ThrowStmt:
		p.w([]byte("throw"))
		p.expr(n.Value)
		p.c(';')
	case *js.TryStmt:
		p.w([]byte("try"))
		p.block(n.Body)
		if n.Catch != nil {
			p.w([]byte("catch"))
			if n.Binding != nil {
				p.c('(')
				p.binding(n.Binding)
				p.c(')')
			}
			p.block(n.Catch)
		}
		if n.Finally != nil {
			p.w([]byte("finally"))
			p.block(n.Finally)
		}
	case *js.DebuggerStmt:
		p.w([]byte("debugger"))
		p.c(';')
	case *js.ImportStmt:
		p.importStmt(n)
	case *js.ExportStmt:
		p.exportStmt(n)
	}
}

func (p *printer) caseClause(n *js.CaseClause) {
	if n.Cond != nil {
		p.w([]byte("case"))
		p.expr(n.Cond)
	} else {
		p.w([]byte("default"))
	}
	p.c(':')
	p.stmts(n.List)
}

func (p *printer) varDecl(n *js.VarDecl) {
	p.w(n.TokenType.Bytes())
	for i := range n.List {
		if i != 0 {
			p.c(',')
		}
		p.bindingEl(&n.List[i])
	}
}

func (p *printer) bindingEl(e *js.BindingElement) {
	if e.Binding == nil {
		return
	}
	p.binding(e.Binding)
	if e.Default != nil {
		p.c('=')
		p.expr(e.Default)
	}
}

func (p *printer) binding(b js.IBinding) {
	switch b := b.(type) {
	case *js.Var:
		p.w(b.Name())
	case *js.BindingArray:
		p.c('[')
		for i := range b.List {
			if i != 0 {
				p.c(',')
			}
			p.bindingEl(&b.List[i])
		}
		if b.Rest != nil {
			if len(b.List) != 0 {
				p.c(',')
			}
			p.w([]byte("..."))
			p.binding(b.Rest)
		} else if 0 < len(b.List) && b.List[len(b.List)-1].Binding == nil {
			p.c(',')
		}
		p.c(']')
	case *js.BindingObject:
		p.c('{')
		for i := range b.List {
			if i != 0 {
				p.c(',')
			}
			item := &b.List[i]
			if item.Key != nil {
				write := true
				if v, ok := item.Value.Binding.(*js.Var); ok && item.Key.IsIdent(v.Data) {
					write = false
				}
				if write {
					p.propName(item.Key)
					p.c(':')
				}
			}
			p.bindingEl(&item.Value)
		}
		if b.Rest != nil {
			if len(b.List) != 0 {
				p.c(',')
			}
			p.w([]byte("..."))
			p.w(b.Rest.Name())
		}
		p.c('}')
	}
}

func (p *printer) propName(n *js.PropertyName) {
	if n.Computed != nil {
		p.c('[')
		p.expr(n.Computed)
		p.c(']')
		return
	}
	p.w(n.Literal.Data)
}

func (p *printer) importStmt(n *js.ImportStmt) {
	p.w([]byte("import"))
	if n.Default != nil {
		p.w(n.Default)
		if n.List != nil {
			p.c(',')
		}
	}
	if len(n.List) == 1 && len(n.List[0].Name) == 1 && n.List[0].Name[0] == '*' {
		p.alias(&n.List[0])
	} else if n.List != nil {
		p.c('{')
		for i := range n.List {
			if i != 0 {
				p.c(',')
			}
			if n.List[i].Binding != nil {
				p.alias(&n.List[i])
			}
		}
		p.c('}')
	}
	if n.Default != nil || n.List != nil {
		p.w([]byte("from"))
	}
	p.w(n.Module)
	p.c(';')
}

func (p *printer) alias(a *js.Alias) {
	if a.Name != nil {
		p.w(a.Name)
		p.w([]byte("as"))
	}
	p.w(a.Binding)
}

func (p *printer) exportStmt(n *js.ExportStmt) {
	p.w([]byte("export"))
	if n.Decl != nil {
		if n.Default {
			p.w([]byte("default"))
		}
		p.expr(n.Decl)
		p.c(';')
		return
	}
	if len(n.List) == 1 && (len(n.List[0].Name) == 1 && n.List[0].Name[0] == '*' || n.List[0].Name == nil && len(n.List[0].Binding) == 1 && n.List[0].Binding[0] == '*') {
		p.alias(&n.List[0])
	} else {
		p.c('{')
		for i := range n.List {
			if i != 0 {
				p.c(',')
			}
			if n.List[i].Binding != nil {
				p.alias(&n.List[i])
			}
		}
		p.c('}')
	}
	if n.Module != nil {
		p.w([]byte("from"))
		p.w(n.Module)
	}
	p.c(';')
}

func (p *printer) funcDecl(n *js.FuncDecl) {
	if n.Async {
		p.w([]byte("async"))
	}
	p.w([]byte("function"))
	if n.Generator {
		p.c('*')
	}
	if n.Name != nil {
		p.w(n.Name.Name())
	}
	p.params(&n.Params)
	p.block(&n.Body)
}

func (p *printer) params(pr *js.Params) {
	p.c('(')
	for i := range pr.List {
		if i != 0 {
			p.c(',')
		}
		p.bindingEl(&pr.List[i])
	}
	if pr.Rest != nil {
		if len(pr.List) != 0 {
			p.c(',')
		}
		p.w([]byte("..."))
		p.binding(pr.Rest)
	}
	p.c(')')
}

func (p *printer) classDecl(n *js.ClassDecl) {
	p.w([]byte("class"))
	if n.Name != nil {
		p.w(n.Name.Name())
	}
	if n.Extends != nil {
		p.w([]byte("extends"))
		p.expr(n.Extends)
	}
	p.c('{')
	for i := range n.List {
		p.classElement(&n.List[i])
	}
	p.c('}')
}

func (p *printer) classElement(n *js.ClassElement) {
	if n.StaticBlock != nil {
		p.w([]byte("static"))
		p.block(n.StaticBlock)
		if p.last != ';' {
			p.c(';')
		}
		return
	}
	if n.Method != nil {
		p.methodDecl(n.Method)
		return
	}
	if n.Field.Static {
		p.w([]byte("static"))
	}
	p.classElementName(&n.Field.Name)
	if n.Field.Init != nil {
		p.c('=')
		p.expr(n.Field.Init)
	}
	p.c(';')
}

func (p *printer) classElementName(n *js.ClassElementName) {
	if n.Private != nil {
		p.w(n.Private.Name())
		return
	}
	p.propName(&n.PropertyName)
}

func (p *printer) methodDecl(n *js.MethodDecl) {
	if n.Static {
		p.w([]byte("static"))
	}
	if n.Async {
		p.w([]byte("async"))
	}
	if n.Generator {
		p.c('*')
	}
	if n.Get {
		p.w([]byte("get"))
	}
	if n.Set {
		p.w([]byte("set"))
	}
	p.classElementName(&n.Name)
	p.params(&n.Params)
	p.block(&n.Body)
}

func (p *printer) expr(n js.IExpr) {
	switch n := n.(type) {
	case *js.Var:
		p.w(n.Name())
	case *js.LiteralExpr:
		p.w(n.Data)
	case js.LiteralExpr:
		p.w(n.Data)
	case *js.GroupExpr:
		p.c('(')
		p.expr(n.X)
		p.c(')')
	case *js.TemplateExpr:
		p.template(n)
	case *js.ArrayExpr:
		p.c('[')
		for i := range n.List {
			if i != 0 {
				p.c(',')
			}
			el := &n.List[i]
			if el.Value != nil {
				if el.Spread {
					p.w([]byte("..."))
				}
				p.expr(el.Value)
			}
		}
		if 0 < len(n.List) && n.List[len(n.List)-1].Value == nil {
			p.c(',')
		}
		p.c(']')
	case *js.ObjectExpr:
		p.c('{')
		for i := range n.List {
			if i != 0 {
				p.c(',')
			}
			prop := &n.List[i]
			if prop.Name != nil {
				write := true
				if v, ok := prop.Value.(*js.Var); ok && prop.Name.IsIdent(v.Data) {
					write = false
				}
				if write {
					p.propName(prop.Name)
					p.c(':')
				}
			} else if prop.Spread {
				p.w([]byte("..."))
			}
			p.expr(prop.Value)
			if prop.Init != nil {
				p.c('=')
				p.expr(prop.Init)
			}
		}
		p.c('}')
	case *js.DotExpr:
		p.dotX(n.X)
		if n.Optional {
			p.w([]byte("?."))
		} else {
			p.c('.')
		}
		p.expr(n.Y)
	case *js.IndexExpr:
		p.expr(n.X)
		if n.Optional {
			p.w([]byte("?.["))
		} else {
			p.c('[')
		}
		p.expr(n.Y)
		p.c(']')
	case *js.NewTargetExpr:
		p.w([]byte("new.target"))
	case *js.ImportMetaExpr:
		p.w([]byte("import.meta"))
	case *js.NewExpr:
		p.w([]byte("new"))
		p.expr(n.X)
		if n.Args != nil {
			p.args(n.Args)
		} else {
			p.w([]byte("()"))
		}
	case *js.CallExpr:
		p.expr(n.X)
		if n.Optional {
			p.w([]byte("?."))
		}
		p.args(&n.Args)
	case *js.UnaryExpr:
		p.unary(n)
	case *js.BinaryExpr:
		p.expr(n.X)
		p.w(n.Op.Bytes())
		p.expr(n.Y)
	case *js.CondExpr:
		p.expr(n.Cond)
		p.c('?')
		p.expr(n.X)
		p.c(':')
		p.expr(n.Y)
	case *js.YieldExpr:
		p.w([]byte("yield"))
		if n.X != nil {
			if n.Generator {
				p.c('*')
			}
			p.expr(n.X)
		}
	case *js.ArrowFunc:
		if n.Async {
			p.w([]byte("async"))
		}
		p.params(&n.Params)
		p.w([]byte("=>"))
		p.block(&n.Body)
	case *js.CommaExpr:
		for i := range n.List {
			if i != 0 {
				p.c(',')
			}
			p.expr(n.List[i])
		}
	case *js.FuncDecl:
		p.funcDecl(n)
	case *js.ClassDecl:
		p.classDecl(n)
	case *js.MethodDecl:
		p.methodDecl(n)
	case *js.VarDecl:
		p.varDecl(n)
	}
}

func (p *printer) template(n *js.TemplateExpr) {
	if n.Tag != nil {
		p.dotX(n.Tag)
		if n.Optional {
			p.w([]byte("?."))
		}
	}
	for i := range n.List {
		p.w(n.List[i].Value)
		p.expr(n.List[i].Expr)
	}
	p.w(n.Tail)
}

func (p *printer) dotX(x js.IExpr) {
	wrap := false
	switch v := x.(type) {
	case *js.LiteralExpr:
		if v.TokenType == js.DecimalToken || v.TokenType == js.IntegerToken {
			wrap = true
		}
	case *js.UnaryExpr:
		if v.Op == js.NegToken || v.Op == js.PosToken {
			if lit, ok := v.X.(*js.LiteralExpr); ok && (lit.TokenType == js.DecimalToken || lit.TokenType == js.IntegerToken) {
				wrap = true
			}
		}
	}
	if wrap {
		p.c('(')
		p.expr(x)
		p.c(')')
		return
	}
	p.expr(x)
}

func (p *printer) unary(n *js.UnaryExpr) {
	if n.Op == js.PostIncrToken || n.Op == js.PostDecrToken {
		p.expr(n.X)
		p.w(n.Op.Bytes())
		return
	}
	p.w(n.Op.Bytes())
	p.expr(n.X)
}

func (p *printer) args(a *js.Args) {
	p.c('(')
	for i := range a.List {
		if i != 0 {
			p.c(',')
		}
		if a.List[i].Rest {
			p.w([]byte("..."))
		}
		p.expr(a.List[i].Value)
	}
	p.c(')')
}
