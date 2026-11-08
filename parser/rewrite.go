package parser

import "github.com/tdewolff/parse/v2/js"

type RewriteFunc func(js.INode) js.INode

func Rewrite(root js.INode, f RewriteFunc) js.INode {
	return rw(root, f)
}

func rwStmt(n js.IStmt, f RewriteFunc) js.IStmt {
	if n == nil {
		return nil
	}
	r := rw(n, f)
	if r == nil {
		return nil
	}
	return r.(js.IStmt)
}

func rwStmts(list []js.IStmt, f RewriteFunc) []js.IStmt {
	for i, s := range list {
		list[i] = rwStmt(s, f)
	}
	return list
}

func rwExpr(n js.IExpr, f RewriteFunc) js.IExpr {
	if n == nil {
		return nil
	}
	r := rw(n, f)
	if r == nil {
		return nil
	}
	return r.(js.IExpr)
}

func rwExprs(list []js.IExpr, f RewriteFunc) []js.IExpr {
	for i, e := range list {
		list[i] = rwExpr(e, f)
	}
	return list
}

func rwBinding(n js.IBinding, f RewriteFunc) js.IBinding {
	if n == nil {
		return nil
	}
	r := rw(n, f)
	if r == nil {
		return nil
	}
	return r.(js.IBinding)
}

func rwVar(n *js.Var, f RewriteFunc) *js.Var {
	if n == nil {
		return nil
	}
	r := f(n)
	if v, ok := r.(*js.Var); ok {
		return v
	}
	return n
}

func rwProps(name *js.PropertyName, f RewriteFunc) *js.PropertyName {
	if name == nil {
		return nil
	}
	if name.Computed != nil {
		name.Computed = rwExpr(name.Computed, f)
	}
	return name
}

func rwBindingElement(e *js.BindingElement, f RewriteFunc) {
	e.Binding = rwBinding(e.Binding, f)
	e.Default = rwExpr(e.Default, f)
}

func rwBindingElements(list []js.BindingElement, f RewriteFunc) {
	for i := range list {
		rwBindingElement(&list[i], f)
	}
}

func rwClassElementName(n *js.ClassElementName, f RewriteFunc) {
	if n.Private != nil {
		rwVar(n.Private, f)
	}
	rwProps(&n.PropertyName, f)
}

func rw(n js.INode, f RewriteFunc) js.INode {
	switch n := n.(type) {
	case *js.AST:
		n.List = rwStmts(n.List, f)
	case *js.Var:
	case *js.BlockStmt:
		n.List = rwStmts(n.List, f)
	case *js.EmptyStmt:
	case *js.ExprStmt:
		n.Value = rwExpr(n.Value, f)
	case *js.IfStmt:
		n.Cond = rwExpr(n.Cond, f)
		n.Body = rwStmt(n.Body, f)
		n.Else = rwStmt(n.Else, f)
	case *js.DoWhileStmt:
		n.Cond = rwExpr(n.Cond, f)
		n.Body = rwStmt(n.Body, f)
	case *js.WhileStmt:
		n.Cond = rwExpr(n.Cond, f)
		n.Body = rwStmt(n.Body, f)
	case *js.ForStmt:
		n.Init = rwExpr(n.Init, f)
		n.Cond = rwExpr(n.Cond, f)
		n.Post = rwExpr(n.Post, f)
		if n.Body != nil {
			n.Body = rw(n.Body, f).(*js.BlockStmt)
		}
	case *js.ForInStmt:
		n.Init = rwExpr(n.Init, f)
		n.Value = rwExpr(n.Value, f)
		if n.Body != nil {
			n.Body = rw(n.Body, f).(*js.BlockStmt)
		}
	case *js.ForOfStmt:
		n.Init = rwExpr(n.Init, f)
		n.Value = rwExpr(n.Value, f)
		if n.Body != nil {
			n.Body = rw(n.Body, f).(*js.BlockStmt)
		}
	case *js.SwitchStmt:
		n.Init = rwExpr(n.Init, f)
		for i := range n.List {
			n.List[i].Cond = rwExpr(n.List[i].Cond, f)
			n.List[i].List = rwStmts(n.List[i].List, f)
		}
	case *js.CaseClause:
		n.Cond = rwExpr(n.Cond, f)
		n.List = rwStmts(n.List, f)
	case *js.BranchStmt:
	case *js.ReturnStmt:
		n.Value = rwExpr(n.Value, f)
	case *js.WithStmt:
		n.Cond = rwExpr(n.Cond, f)
		n.Body = rwStmt(n.Body, f)
	case *js.LabelledStmt:
		n.Value = rwStmt(n.Value, f)
	case *js.ThrowStmt:
		n.Value = rwExpr(n.Value, f)
	case *js.TryStmt:
		if n.Body != nil {
			n.Body = rw(n.Body, f).(*js.BlockStmt)
		}
		if n.Catch != nil {
			n.Catch = rw(n.Catch, f).(*js.BlockStmt)
		}
		if n.Finally != nil {
			n.Finally = rw(n.Finally, f).(*js.BlockStmt)
		}
		n.Binding = rwBinding(n.Binding, f)
	case *js.DebuggerStmt:
	case *js.Comment:
	case *js.ImportStmt:
	case *js.ExportStmt:
		n.Decl = rwExpr(n.Decl, f)
	case *js.DirectivePrologueStmt:
	case *js.PropertyName:
		if n.Computed != nil {
			n.Computed = rwExpr(n.Computed, f)
		}
	case *js.BindingArray:
		rwBindingElements(n.List, f)
		n.Rest = rwBinding(n.Rest, f)
	case *js.BindingObject:
		for i := range n.List {
			n.List[i].Key = rwProps(n.List[i].Key, f)
			rwBindingElement(&n.List[i].Value, f)
		}
		if n.Rest != nil {
			n.Rest = rwVar(n.Rest, f)
		}
	case *js.BindingObjectItem:
		n.Key = rwProps(n.Key, f)
		rwBindingElement(&n.Value, f)
	case *js.BindingElement:
		rwBindingElement(n, f)
	case *js.VarDecl:
		rwBindingElements(n.List, f)
	case *js.Params:
		rwBindingElements(n.List, f)
		n.Rest = rwBinding(n.Rest, f)
	case *js.FuncDecl:
		rwVar(n.Name, f)
		rwBindingElements(n.Params.List, f)
		n.Body.List = rwStmts(n.Body.List, f)
	case *js.MethodDecl:
		rwClassElementName(&n.Name, f)
		rwBindingElements(n.Params.List, f)
		n.Body.List = rwStmts(n.Body.List, f)
	case *js.Field:
		rwClassElementName(&n.Name, f)
		n.Init = rwExpr(n.Init, f)
	case *js.ClassDecl:
		rwVar(n.Name, f)
		n.Extends = rwExpr(n.Extends, f)
		for i := range n.List {
			el := &n.List[i]
			if el.StaticBlock != nil {
				el.StaticBlock.List = rwStmts(el.StaticBlock.List, f)
			} else if el.Method != nil {
				rw(el.Method, f)
			} else {
				rw(&el.Field, f)
			}
		}
	case *js.ClassElement:
	case *js.ClassElementName:
	case *js.LiteralExpr:
	case *js.Element:
		n.Value = rwExpr(n.Value, f)
	case *js.ArrayExpr:
		for i := range n.List {
			n.List[i].Value = rwExpr(n.List[i].Value, f)
		}
	case *js.Property:
		n.Name = rwProps(n.Name, f)
		n.Value = rwExpr(n.Value, f)
		n.Init = rwExpr(n.Init, f)
	case *js.ObjectExpr:
		for i := range n.List {
			n.List[i].Name = rwProps(n.List[i].Name, f)
			n.List[i].Value = rwExpr(n.List[i].Value, f)
			n.List[i].Init = rwExpr(n.List[i].Init, f)
		}
	case *js.TemplatePart:
		n.Expr = rwExpr(n.Expr, f)
	case *js.TemplateExpr:
		for i := range n.List {
			n.List[i].Expr = rwExpr(n.List[i].Expr, f)
		}
		n.Tag = rwExpr(n.Tag, f)
	case *js.GroupExpr:
		n.X = rwExpr(n.X, f)
	case *js.IndexExpr:
		n.X = rwExpr(n.X, f)
		n.Y = rwExpr(n.Y, f)
	case *js.DotExpr:
		n.X = rwExpr(n.X, f)
		n.Y = rwExpr(n.Y, f)
	case *js.NewTargetExpr:
	case *js.ImportMetaExpr:
	case *js.Arg:
		n.Value = rwExpr(n.Value, f)
	case *js.Args:
		for i := range n.List {
			n.List[i].Value = rwExpr(n.List[i].Value, f)
		}
	case *js.NewExpr:
		n.X = rwExpr(n.X, f)
		if n.Args != nil {
			for i := range n.Args.List {
				n.Args.List[i].Value = rwExpr(n.Args.List[i].Value, f)
			}
		}
	case *js.CallExpr:
		n.X = rwExpr(n.X, f)
		for i := range n.Args.List {
			n.Args.List[i].Value = rwExpr(n.Args.List[i].Value, f)
		}
	case *js.UnaryExpr:
		n.X = rwExpr(n.X, f)
	case *js.BinaryExpr:
		n.X = rwExpr(n.X, f)
		n.Y = rwExpr(n.Y, f)
	case *js.CondExpr:
		n.Cond = rwExpr(n.Cond, f)
		n.X = rwExpr(n.X, f)
		n.Y = rwExpr(n.Y, f)
	case *js.YieldExpr:
		n.X = rwExpr(n.X, f)
	case *js.ArrowFunc:
		rwBindingElements(n.Params.List, f)
		n.Body.List = rwStmts(n.Body.List, f)
	case *js.CommaExpr:
		n.List = rwExprs(n.List, f)
	}
	return f(n)
}
