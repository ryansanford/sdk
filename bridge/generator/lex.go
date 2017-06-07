package main

import (
	. "fmt"
)

func main() {
	// Create frontend printers
	printers := []Printer{
		NewBasicPrinter("go", "templates/template.go.", "dist/bridge.go"),
		NewBasicPrinter("py", "templates/template.py", "dist/python/flywheel.py"),
		NewBasicPrinter("m", "templates/template.m", "dist/matlab/Bridge.m"),
	}

	// Parse source code into useful sigs
	Println("Parsing SDK...")
	fset, sigs := GenerateSignatures("../api")
	_ = fset

	// Load templates and execute
	for _, printer := range printers {
		Println("Generating", printer.Name(), "code...")

		printer.Init()
		printer.Print(sigs)
	}
}
