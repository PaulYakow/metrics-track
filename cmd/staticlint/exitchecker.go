package main

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// ExitAnalyzer анализатор, проверяющий есть ли прямой вызов os.Exit в функции main пакета main
var ExitAnalyzer = &analysis.Analyzer{
	Name:     "osexit",
	Doc:      "check for direct call os.Exit in main function of package main",
	Run:      exitAnalyze,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func exitAnalyze(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		if strings.HasSuffix(pass.Pkg.Path(), "test") {
			return
		}

		funcDecl := node.(*ast.FuncDecl)
		if funcDecl.Name.Name != "main" {
			return
		}

		v := visitor{pass: pass}
		for _, stmt := range funcDecl.Body.List {
			ast.Walk(v, stmt)
		}
	})
	return nil, nil
}

type visitor struct {
	pass *analysis.Pass
}

func (v visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return v
	}

	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return v
	}

	if selectorExpr.Sel.Name == "Exit" {
		v.pass.Reportf(node.Pos(), "using os.Exit in main func of main package")
		return nil
	}

	return v
}
