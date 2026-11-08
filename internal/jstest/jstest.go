package jstest

import (
	"strings"

	"github.com/dop251/goja"
)

func Eval(code string) (string, error) {
	v, err := goja.New().RunString(code)
	if err != nil {
		return "", err
	}
	return v.String(), nil
}

func Expr(expr string) (string, error) {
	return Eval("JSON.stringify(" + strings.Trim(expr, "; \n\t") + ")")
}
