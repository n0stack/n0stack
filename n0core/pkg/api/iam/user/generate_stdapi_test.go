package user

import (
	"testing"

	stdapi "github.com/n0stack/n0stack/n0core/pkg/api/standard_api"
	"github.com/n0stack/n0stack/n0core/pkg/util/generator"
)

func TestGenerate(t *testing.T) {
	g := generator.NewGoCodeGenerator("standard_api", "user")
	stdapi.GenerateTemplateAPI(g, "iam", "User", "v1alpha")

	if err := g.WriteAsTemplateFileName(); err != nil {
		t.Errorf("err=%s", err.Error())
	}
}
