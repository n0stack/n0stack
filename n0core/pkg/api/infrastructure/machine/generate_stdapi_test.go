package amachine

import (
	"testing"

	stdapi "n0st.ac/n0stack/n0core/pkg/api/standard_api"
	"n0st.ac/n0stack/n0core/pkg/util/generator"
)

func TestGenerate(t *testing.T) {
	g := generator.NewGoCodeGenerator("standard_api", "amachine")
	stdapi.GenerateTemplateAPI(g, "infrastructure", "Machine")

	if err := g.WriteAsTemplateFileName(); err != nil {
		t.Errorf("err=%s", err.Error())
	}
}
