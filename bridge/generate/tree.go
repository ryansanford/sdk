package main

import (
	"go/ast"
)

type Param struct {
	Name  string
	Type  string
	CType string
}

// An accessible representation of a function signature used for templating.
type Signature struct {
	Name    string
	Params  []*Param
	Results []*Param

	// various template conveniences

	LastResultIndex int
	LastParamIndex  int

	ParamDataName string
	ParamDataType string

	ReturnDataName string
}

// This is a full struct instead of a type alias for text-templating reasons
type ParsedSignatures struct {
	Signatures []*Signature
}

func isStringExpr(ex ast.Expr) bool {
	ident, ok := ex.(*ast.Ident)
	return ok && ident.Name == "string"
}

func isDataExpr(ex ast.Expr) (bool, string) {
	// might be an array of pointers; if so, unwrap
	array, ok := ex.(*ast.ArrayType)
	if ok {
		ex = array.Elt
	}

	pointer, ok := ex.(*ast.StarExpr)
	if !ok {
		return false, ""
	}

	ident, ok := pointer.X.(*ast.Ident)
	if !ok {
		return false, ""
	}
	name := ident.Name

	// Could replace with lexing later
	whitelist := []string{"Acquisition", "Batch", "BatchProposal", "Client", "Config", "ContainerReference", "DeletedResponse", "Error", "FileReference", "Formula", "FormulaResult", "Gear", "GearDoc", "GearSource", "Group", "IdResponse", "Input", "Job", "JobLog", "JobLogStatement", "Key", "ModifiedResponse", "Note", "Origin", "Output", "Permission", "ProgressReader", "Project", "Result", "Session", "Subject", "Target", "UploadResponse", "UploadSource", "User", "Version"}

	if stringInSlice(name, whitelist) {
		return true, name
	} else {
		return false, name
	}
}

func isHttpRespExpr(ex ast.Expr) bool {
	pointer, ok := ex.(*ast.StarExpr)
	if !ok {
		return false
	}

	selector, ok := pointer.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := selector.X.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "http" && selector.Sel.Name == "Response"
}

func isErrorExpr(ex ast.Expr) bool {
	ident, ok := ex.(*ast.Ident)
	return ok && ident.Name == "error"
}
