package project

import (
	"testing"

	stdapi "n0st.ac/n0stack/n0core/pkg/api/standard_api"
	"n0st.ac/n0stack/n0core/pkg/util/generator"
)

func TestGenerate(t *testing.T) {
	g := generator.NewGoCodeGenerator("standard_api", "project")
	stdapi.GenerateTemplateAPI(g, "iam", "Project", "v1alpha")
	stdapi.GenerateTemplatedListAPI(g, "iam", "Project")

	if err := g.WriteAsTemplateFileName(); err != nil {
		t.Errorf("err=%s", err.Error())
	}
}
