package main

import (
	. "fmt"
)

func main() {
	// Create frontend printers
	printers := []Printer {
		&BasicPrinter{ name: "go", templatePath: "templates/template.go.", outputPath: "bridge.go" },
		&BasicPrinter{ name: "py", templatePath: "templates/template.py",  outputPath: "python/flywheel.py" },
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
