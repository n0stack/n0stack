package virtualmachine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	"n0st.ac/n0stack/n0core/pkg/datastore/memory"
	"n0st.ac/n0stack/n0proto.go/provisioning/v0"
)

func TestCreateVirtualMachineOnN0test(t *testing.T) {
	raw := os.Getenv("N0TEST_JSON_CreateVirtualMachine_REQUESTS")
	if raw == "" {
		// b, err := ioutil.ReadFile("CreateVirtualMachine.n0test.json")
		b, err := ioutil.ReadFile("n0test.CreateVirtualMachine.json")
		if err != nil {
			t.Fatalf("Failed to read CreateVirtualMachine.n0test.json: err=%s", err.Error())
		}

		raw = string(b)
	}

	var rawList []interface{}
	if err := json.Unmarshal([]byte(raw), &rawList); err != nil {
		t.Fatalf("Failed to parse N0TEST_JSON_CreateVirtualMachine_REQUESTS by JSON: err=%s", err.Error())
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	m := memory.NewMemoryDatastore()
	vma := NewMockVirtualMachineAPI(m)

	mnode, err := vma.NodeAPI.SetupMockNode(ctx)
	if err != nil {
		t.Fatalf("Failed to set up mocked node: err=%s", err.Error())
	}

	_, err = vma.NetworkAPI.FactoryNetwork(ctx)
	if err != nil {
		t.Fatalf("Failed to factory network: err='%s'", err.Error())
	}

	_, err = vma.BlockStorageAPI.FactoryBlockStorage(ctx, mnode.Name)
	if err != nil {
		t.Fatalf("Failed to factory blockstorage: err='%s'", err.Error())
	}

	baseDatastore := vma.api.dataStore
	for i, r := range rawList {
		vma.api.dataStore = baseDatastore.AddPrefix(fmt.Sprintf("test%d", i))

		b, _ := json.Marshal(r)
		reader := bytes.NewReader(b)
		req := &pprovisioning.CreateVirtualMachineRequest{}

		if err := jsonpb.Unmarshal(reader, req); err != nil {
			t.Fatalf("[N0TEST_OMIT] Failed to parse N0TEST_JSON_CreateVirtualMachine_REQUESTS from JSON to pb")
		}

		vma.CreateVirtualMachine(ctx, req)
	}
}
