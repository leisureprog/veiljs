package parser

import (
	"strings"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/js"
)

func Parse(code string) (*js.AST, error) {
	return js.Parse(parse.NewInputString(code), js.Options{})
}

func Print(ast *js.AST) string {
	var sb strings.Builder
	ast.JS(&sb)
	return sb.String()
}
