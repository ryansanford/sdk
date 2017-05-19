package main

import (
	. "fmt"
	"go/ast"
	"os"
	"reflect"
)

// Print a reasonably-friendly string representation of a function signature.
func DebugPrintFuncSig(fn *ast.FuncDecl) {
	name := fn.Name.Name
	parameters := fn.Type.Params.List
	results := fn.Type.Results.List
	str := name

	str += "\n\t"
	for _, paramSet := range parameters {
		for _, param := range paramSet.Names {
			str += param.Name + " " + ExprName(paramSet.Type) + " "
		}
	}

	if len(results) > 0 {
		str += "\n\t"
	}
	for _, resultSet := range results {
		str += ExprName(resultSet.Type) + " "
	}
	Println(str)
}

// Given an expression, print its human-readable type.
func ExprName(expr ast.Expr) string {

	switch n := expr.(type) {

	// Basic identifier
	case *ast.Ident:
		return Sprint(expr)

	// Array
	case *ast.ArrayType:
		return "[]" + ExprName(n.Elt)

	// struct
	case *ast.StructType:
		return "struct"

	// pointer
	case *ast.StarExpr:
		return "*" + ExprName(n.X)

	// Package selector, such as http.Response
	case *ast.SelectorExpr:
		return ExprName(n.X) + "." + n.Sel.Name

	default:
		return Sprint(reflect.TypeOf(expr))
	}
}

// standard array check...
func stringInSlice(str string, slice []string) bool {
	for _, x := range slice {
		if x == str {
			return true
		}
	}
	return false
}

// standard fail on err check...
func check(err error) {
	if err != nil {
		Println(err)
		os.Exit(1)
	}
}
