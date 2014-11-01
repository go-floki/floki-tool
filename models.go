package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	//	"reflect"
)

type ModelCollector struct {
	models map[string]*Model
}

type ParsedNode interface {
}

type Model struct {
	Name   string
	Fields []*Field
}

type Field struct {
	Names []string
	Type  string
	Tag   string
}

type TextNode struct {
	Value string
}

func ParseModels(dir string) map[string]*Model {
	fset := token.NewFileSet()

	packages, err := parser.ParseDir(fset, dir+"/models", nil, 0)
	if err != nil {
		log.Println(err)
		return nil
	}

	m := ModelCollector{}
	m.models = make(map[string]*Model)

	// Print the imports from the file's AST.
	for _, p := range packages {
		for fileName, f := range p.Files {
			log.Println("processing:", fileName)
			m.CollectModels(f)
		}
	}

	return m.models
}

func (m *ModelCollector) CollectModels(file *ast.File) {
	m.walk(file)

}

func (m *ModelCollector) walk(node ast.Node) ParsedNode {
	if node == nil {
		return nil
	}

	//v := reflect.ValueOf(node)
	//fmt.Println("type decl:", v.Type())

	switch n := node.(type) {
	case *ast.Field:
		field := &Field{}

		if n.Doc != nil {
			m.walk(n.Doc)
		}

		for _, x := range n.Names {
			textNode := m.walk(x).(TextNode)
			field.Names = append(field.Names, textNode.Value)
		}

		tn := m.walk(n.Type)
		if tn != nil {
			typeNode := tn.(TextNode)
			field.Type = typeNode.Value
		}

		//v := reflect.ValueOf(n.Type)
		//fmt.Println("type decl:", v.Type())
		m.walk(n.Type)

		if n.Tag != nil {
			textNode := m.walk(n.Tag).(TextNode)
			field.Tag = textNode.Value
		}

		if n.Comment != nil {
			m.walk(n.Comment)
		}

		return field

	case *ast.FieldList:
		fields := make([]*Field, 0)
		for _, f := range n.List {
			fields = append(fields, m.walk(f).(*Field))
		}

		model := Model{"", fields}

		return model

	case *ast.StructType:
		return m.walk(n.Fields)

	case *ast.TypeSpec:
		if n.Doc != nil {
			m.walk(n.Doc)
		}

		idNode := m.walk(n.Name).(TextNode)
		nt := m.walk(n.Type)
		if nt != nil {
			model := nt.(Model)
			model.Name = idNode.Value

			log.Println("found model:", idNode.Value)

			m.models[model.Name] = &model
		}

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
