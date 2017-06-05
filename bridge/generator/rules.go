package main

import (
	"text/template"
	"unicode"
)

// Add these to your template with template.Funcs!
// Then, you can call these in the template with the assigned alias.
var textRules = template.FuncMap{
	"camel2snake": UpperCamelCaseToSnakeCase,
	"camel2lowercamel": UpperCamelCaseToLowerCamelCase,
}

// A great template function guide:
// https://www.calhoun.io/using-functions-inside-go-templates

// Translate a function name from upper camel case, such as "DoStuff", to snake case, such as "do_stuff".
// Written for translating Golang-friendly function names to Python-friendly ones, may be useful elsewhere.
func UpperCamelCaseToSnakeCase(funcName string) string {
	result := ""

	// Iterate over Unicode code points
	for pos, x := range funcName {

		isUpperCase := unicode.IsUpper(x)

		if isUpperCase && pos == 0 {
			// First character; don't prefix with an underscore
			result += string(unicode.ToLower(x))
		} else if isUpperCase && pos != 0 {
			// Not character; don't prefix with an underscore
			result += "_" + string(unicode.ToLower(x))
		} else {
			// Lower case
			result += string(x)
		}
	}

	return result
}

// Translate a function name from upper camel case, such as "DoStuff", to lower camel case, such as "doStuff".
// Written for translating Golang-friendly function names to Matlab-friendly ones, may be useful elsewhere.
func UpperCamelCaseToLowerCamelCase(funcName string) string {
	result := []rune(funcName)

	x := rune(result[0])

	if unicode.IsUpper(x) {
		result[0] = unicode.ToLower(x)
	}

	return string(result)
}
