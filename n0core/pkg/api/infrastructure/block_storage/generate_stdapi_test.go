package ablockstorage

import (
	"testing"

	stdapi "n0st.ac/n0stack/n0core/pkg/api/standard_api"
	"n0st.ac/n0stack/n0core/pkg/util/generator"
)

func TestGenerate(t *testing.T) {
	g := generator.NewGoCodeGenerator("standard_api", "ablockstorage")
	stdapi.GenerateTemplateAPI(g, "infrastructure", "BlockStorage")

	if err := g.WriteAsTemplateFileName(); err != nil {
		t.Errorf("err=%s", err.Error())
	}
}