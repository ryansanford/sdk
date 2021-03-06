package main

import (
	"bytes"
	"io/ioutil"
	"text/template"
)

// Printer specifies a struct that can generate an SDK frontend.
type Printer interface {

	// Name returns the name of the printer.
	Name() string

	// Init loads any templates or other one-time work necessary.
	Init()

	// Print will generate the frontend and write the resultant file(s).
	Print(*ParsedSignatures)
}

// BasicPrinter reads a single template and prints to a single file.
type BasicPrinter struct {
	name         string
	templatePath string
	outputPath   string

	template *template.Template
	buffer   bytes.Buffer
}

var _ Printer = (*BasicPrinter)(nil)

func NewBasicPrinter(name, templatePath, outputPath string) Printer {
	return &BasicPrinter{
		name:         name,
		templatePath: templatePath,
		outputPath:   outputPath,
	}
}

func (p *BasicPrinter) Name() string {
	return p.name
}

func (p *BasicPrinter) Init() {
	templateBytes, err := ioutil.ReadFile(p.templatePath)
	check(err)

	temp, err := template.
		New("printerTemplate").
		Funcs(textRules). // load rules defined in rules.go
		Parse(string(templateBytes))
	check(err)

	p.template = temp
}

func (p *BasicPrinter) Print(sigs *ParsedSignatures) {

	err := p.template.Execute(&p.buffer, sigs)
	check(err)

	bufferBytes := p.buffer.Bytes()

	err = ioutil.WriteFile(p.outputPath, bufferBytes, 0644)
	check(err)
}
