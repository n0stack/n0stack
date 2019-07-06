package generator

import (
	"os"
	"testing"
)

func TestGoCodeGenerator(t *testing.T) {
	g := NewGoCodeGenerator("test_generator", "generator")
	if err := g.WriteWithTemplateFileName(); err != nil {
		t.Errorf("err=%s, src=\n%s", err.Error(), g.code.String())
	}
	defer os.Remove("test_generator.go")
}
