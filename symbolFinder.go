package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	//	"reflect"
)

type SymbolCollector struct {
	symbols map[string]string
}

func ParseSymbols(dir string) map[string]string {
	fset := token.NewFileSet()

	packages, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	m := SymbolCollector{}
	m.symbols = make(map[string]string)

	// Print the imports from the file's AST.
	for _, p := range packages {
		m.walk(p)
	}

	return m.symbols
}

func (m *SymbolCollector) walk(node ast.Node) ParsedNode {
	if node == nil {
		return nil
	}

	//v := reflect.ValueOf(node)
	//fmt.Println("type decl:", v.Type())

	switch n := node.(type) {
	case *ast.Field:
		if n.Doc != nil {
			m.walk(n.Doc)
		}

		for _, x := range n.Names {
			m.walk(x)
		}

		m.walk(n.Type)

		//v := reflect.ValueOf(n.Type)
		//fmt.Println("type decl:", v.Type())
		m.walk(n.Type)

		if n.Tag != nil {
			m.walk(n.Tag)
		}

		if n.Comment != nil {
			m.walk(n.Comment)
		}

	case *ast.FuncDecl:
		node := m.walk(n.Name).(TextNode)
		m.symbols[node.Value] = "1"

	case *ast.FieldList:
		for _, f := range n.List {
			m.walk(f)
		}

	case *ast.StructType:
		return m.walk(n.Fields)

	case *ast.TypeSpec:
		if n.Doc != nil {
			m.walk(n.Doc)
		}

		tn := m.walk(n.Name).(TextNode)
		m.symbols[tn.Value] = "1"
		m.walk(n.Type)

		if n.Comment != nil {
			m.walk(n.Comment)
		}

	case *ast.SelectorExpr:
		x := m.walk(n.X).(TextNode)
		sel := m.walk(n.Sel).(TextNode)

		return TextNode{x.Value + "." + sel.Value}

	case *ast.Ident:
		return TextNode{n.Name}

	case *ast.BasicLit:
		return TextNode{n.Value}

	case *ast.ImportSpec:
		if n.Doc != nil {
			m.walk(n.Doc)
		}
		if n.Name != nil {
			m.walk(n.Name)
		}
		m.walk(n.Path)
		if n.Comment != nil {
			m.walk(n.Comment)
		}

	case *ast.GenDecl:
		for _, s := range n.Specs {
			m.walk(s)
		}

	case *ast.File:
		if n.Doc != nil {
			m.walk(n.Doc)
		}
		m.walk(n.Name)
		for _, typeDecl := range n.Decls {
			m.walk(typeDecl)
		}

	case *ast.Package:
		for _, f := range n.Files {
			m.walk(f)
		}
	}

	return nil
}
