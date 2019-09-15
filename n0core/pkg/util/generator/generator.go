package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
)

// referenced github.com/golang/tools/blob/master/cmd/stringer
type GoCodeGenerator struct {
	code bytes.Buffer

	generator string
	pkg       string
}

func NewGoCodeGenerator(generator string, pkg string) *GoCodeGenerator {
	g := &GoCodeGenerator{
		generator: generator,
		pkg:       pkg,
	}

	g.Printf("// Code generated by \"%s\"; DO NOT EDIT.\n", g.generator)
	g.Printf("\n")
	g.Printf("package %s\n", g.pkg)
	g.Printf("\n")

	return g
}

func (g *GoCodeGenerator) Printf(format string, a ...interface{}) (int, error) {
	return g.code.WriteString(fmt.Sprintf(format, a...))
}

func (g GoCodeGenerator) FormattedCode() ([]byte, error) {
	return format.Source(g.code.Bytes())
}

func (g GoCodeGenerator) WriteFile(filename string) error {
	src, err := g.FormattedCode()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, src, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (g GoCodeGenerator) WriteAsTemplateFileName() error {
	return g.WriteFile(fmt.Sprintf("%s.generated.go", g.generator))
}