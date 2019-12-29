package grpccmd

import (
	"bytes"
	"testing"

	"n0st.ac/n0stack/n0core/pkg/datastore"
)

func TestOutputTable(t *testing.T) {
	test := &datastore.Test{
		Name: "testing",
	}

	buf := &bytes.Buffer{}
	outputter := Outputter{
		Out: buf,
	}
	if err := outputter.OutputTable(test, []string{"name"}); err != nil {
		t.Errorf("Failed to OutputTable: err=%s", err.Error())
	}

	// TODO: compare output
}
