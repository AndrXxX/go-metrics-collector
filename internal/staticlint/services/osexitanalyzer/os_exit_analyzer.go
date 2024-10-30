package osexitanalyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var OSExitAnalyzer = &analysis.Analyzer{
	Name: "osexitanalyzer",
	Doc:  "check for usage os.Exit() method",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.File: // пакет
				return x.Name.Name == "main"
			case *ast.FuncDecl: // объявление функции
				return x.Name.Name == "main"
			case *ast.CallExpr: // выражение
				if isOSExitFunc(x) {
					pass.Reportf(x.Pos(), "exit functiondirect call of function os.Exit()")
				}
			}
			return true
		})
	}
	return nil, nil
}

// isOSExitFunc возвращает true, если вызываемый метод - это os.Exit()
func isOSExitFunc(call *ast.CallExpr) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	xIdent, ok := selector.X.(*ast.Ident)
	if !ok || xIdent.Name != "os" {
		return false
	}
	if selector.Sel.Name != "Exit" {
		return false
	}
	return true
}
