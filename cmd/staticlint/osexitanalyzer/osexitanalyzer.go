// Package osexitanalyzer checks for using os.Exit in main function of main package.
package osexitanalyzer

import (
	"go/ast"
	"log"

	"golang.org/x/tools/go/analysis"
)

var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexitmain",
	Doc:  "check for calling os exit in main",
	Run:  run,
}

func isOsExitCall(callExpr *ast.CallExpr) bool {
	selectorExpression, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if selectorExpression.Sel.Name != "Exit" {
		return false
	}
	pkgIdent, ok := selectorExpression.X.(*ast.Ident)
	if !ok {
		return false
	}
	return pkgIdent.Name == "os"
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}
	for _, file := range pass.Files {
		var mainFunction bool
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.FuncDecl:
				mainFunction = x.Name.Name == "main"
			case *ast.CallExpr:
				if mainFunction && isOsExitCall(x) {
					log.Println("Found a os.Exit in main function")
					pass.Reportf(x.Pos(), "call to os.Exit in main function is forbidden")
				}
			}
			return true
		})
	}
	return nil, nil
}
