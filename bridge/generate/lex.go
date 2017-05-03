package main

import (
	"bytes"
	"encoding/json"
	"errors"
	. "fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"text/template"
)

func check(err error) {
	if err != nil {
		Println(err)
		os.Exit(1)
	}
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

// Scan & parse a package of source code and return the set of function signatures we want to compile.
func GetRelevantFunctionSignatures(path string) (*token.FileSet, []*ast.FuncDecl) {
	fset := token.NewFileSet()

	// Wall-clock warning: discover and parse source code from disk
	packages, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	check(err)

	apiPackage, ok := packages["api"]
	if !ok {
		check(errors.New("Package API not found. Wrong directory path?"))
	}

	funcs := []*ast.FuncDecl{}

	// Inspect the AST, finding all function declarations and filtering
	ast.Inspect(apiPackage, func(n ast.Node) bool {

		// Filter on function declarations
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		// Ignore closures
		if fn.Name == nil || fn.Type.Params == nil || fn.Type.Results == nil {
			return true
		}

		// Ignore functions that are not attached to a receiver
		if fn.Recv == nil || len(fn.Recv.List) != 1 {
			return true
		}

		// Check that the function receiver is *api.Client
		pointer, ok := fn.Recv.List[0].Type.(*ast.StarExpr)
		if !ok {
			return true
		}
		identifier, ok := pointer.X.(*ast.Ident)
		if !ok {
			return true
		}
		if identifier.Name != "Client" {
			return true
		}

		name := fn.Name.Name

		// Ignore unexported functions
		if strings.ToLower(string(name[0])) == string(name[0]) {
			return true
		}

		// Blacklist file I/O
		if strings.Contains(name, "Upload") || strings.Contains(name, "Download") {
			return false
		}

		// Troublesome functions that need to either be fixed or ignored
		blacklist := []string{
			// boolean parameter
			"ModifyJob",

			// variadic string array
			"StartNextPendingJob",

			// JobState enum
			"ChangeJobState",

			// map string -> interface
			"ProposeBatch",

			// api.JobLogStatement instead of []*JobLogStatement
			"AddJobLogs",

			// Doesn't detect int return
			"CancelBatch",

			// Doesn't detect map[string]interface{} return
			"GetGearInvocation",
		}
		if stringInSlice(name, blacklist) {
			return false
		}

		// Println(name)

		funcs = append(funcs, fn)
		return false // don't need to go deeper; we found what we're looking for.
	})

	return fset, funcs
}

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

func stringInSlice(str string, slice []string) bool {
	for _, x := range slice {
		if x == str {
			return true
		}
	}
	return false
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

type Param struct {
	Name  string
	Type  string
	CType string
}

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

func main() {
	fset, funcs := GetRelevantFunctionSignatures("../api/")
	_ = fset

	goTemplateRaw, err := ioutil.ReadFile("generate/template.go.")
	check(err)
	goTemplate, err := template.New("bridgeFuncs").Parse(string(goTemplateRaw))
	check(err)
	goPreamble, err := ioutil.ReadFile("generate/preamble.go")
	check(err)
	goBridgeBuffer := bytes.NewBuffer(goPreamble)

	pyTemplateRaw, err := ioutil.ReadFile("generate/template.py")
	check(err)
	pyTemplate, err := template.New("pyClientFuncs").Parse(string(pyTemplateRaw))
	check(err)
	pyPreamble, err := ioutil.ReadFile("generate/preamble.py")
	check(err)
	pyBridgeBuffer := bytes.NewBuffer(pyPreamble)

	for _, fn := range funcs {
		// DebugPrintFuncSig(fn)
		// ast.Print(fset, fn)

		name := fn.Name.Name

		signature := &Signature{
			Name:           name,
			Params:         []*Param{},
			Results:        []*Param{},
			ParamDataName:  "",
			ReturnDataName: "nil",
		}

		parameters := fn.Type.Params.List
		for _, paramSet := range parameters {
			// ast.Print(fset, paramSet.Type)

			pt := "unknown"
			cpt := "unknown"

			dExpr, dname := isDataExpr(paramSet.Type)

			// Check param type
			if isStringExpr(paramSet.Type) {
				pt = "string"
				cpt = "*C.char"
			} else if dExpr {
				pt = "data"
				cpt = "*C.char"
			}

			for _, param := range paramSet.Names {
				signature.Params = append(signature.Params, &Param{
					Name:  param.Name,
					Type:  pt,
					CType: cpt,
				})

				if pt == "data" {
					signature.ParamDataName = param.Name
					signature.ParamDataType = dname
				}
			}
		}
		signature.LastParamIndex = len(signature.Params) - 1

		results := fn.Type.Results.List
		for _, field := range results {
			pt := "unknown"
			name := "unknown"

			dExpr, dname := isDataExpr(field.Type)
			_ = dname

			if isStringExpr(field.Type) {
				pt = "string"
				name = "data"
				signature.ReturnDataName = "data"
			} else if dExpr {
				pt = "data"
				name = "data"
				signature.ReturnDataName = "data"
			} else if isHttpRespExpr(field.Type) {
				pt = "http"
				name = "_"
			} else if isErrorExpr(field.Type) {
				pt = "error"
				name = "err"
			} else {

				ast.Print(fset, field.Type)
			}

			// Println(pt, name)

			signature.Results = append(signature.Results, &Param{
				Name: name,
				Type: pt,
			})
		}
		signature.LastResultIndex = len(signature.Results) - 1

		x, _ := json.MarshalIndent(signature, "", "\t")
		_ = x
		// Println(string(x))

		err = goTemplate.Execute(goBridgeBuffer, signature)
		check(err)
		err = pyTemplate.Execute(pyBridgeBuffer, signature)
		check(err)
	}

	goBridgeResult := goBridgeBuffer.Bytes()
	// Println(string(goBridgeResult))
	err = ioutil.WriteFile("bridge.go", goBridgeResult, 0644)
	check(err)

	pyBridgeResult := pyBridgeBuffer.Bytes()
	err = ioutil.WriteFile("python/flywheel.py", pyBridgeResult, 0644)
	check(err)
}
