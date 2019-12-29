package blockstorage

import (
	"testing"

	stdapi "n0st.ac/n0stack/n0core/pkg/api/standard_api"
	"n0st.ac/n0stack/n0core/pkg/util/generator"
)

func TestGenerate(t *testing.T) {
	g := generator.NewGoCodeGenerator("template_api", "blockstorage")
	stdapi.GenerateTemplateAPI(g, "provisioning", "BlockStorage")

	if err := g.WriteAsTemplateFileName(); err != nil {
		t.Errorf("err=%s", err.Error())
	}
}
