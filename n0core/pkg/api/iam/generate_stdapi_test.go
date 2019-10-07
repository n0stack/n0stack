package iam

import (
	"testing"

	stdapi "n0st.ac/n0stack/n0core/pkg/api/stdapi"
	"n0st.ac/n0stack/n0core/pkg/util/generator"
)

func TestGenerateProject(t *testing.T) {
	g := generator.NewGoCodeGenerator("stdapi", "iam")
	stdapi.GenerateTemplateAPI(g, "iam", "Project", "v1alpha")

	if err := g.WriteFile("project.generated.go"); err != nil {
		t.Fatalf("err=%+v", err)
	}
}

func TestGenerateServiceAccount(t *testing.T) {
	g := generator.NewGoCodeGenerator("stdapi", "iam")
	stdapi.GenerateTemplateAPI(g, "iam", "ServiceAccount", "v1alpha")

	if err := g.WriteFile("service_account.generated.go"); err != nil {
		t.Fatalf("err=%+v", err)
	}
}
