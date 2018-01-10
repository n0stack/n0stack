package resource

import (
	"net"
	"reflect"
	"testing"

	"github.com/n0stack/n0core/model"
	"github.com/satori/go.uuid"
)

func TestNICGetModel(t *testing.T) {
	id := uuid.NewV4()
	m := &model.Model{
		ID:           id,
		Type:         "test/test",
		State:        "testing",
		Name:         "test_model",
		Meta:         map[string]string{"hoge": "hoge"},
		Dependencies: model.Dependencies{},
	}

	h, _ := net.ParseMAC("01:23:45:67:89:ab")

	n := &NIC{
		Model:   *m,
		HWAddr:  h,
		IPAddrs: []net.IP{net.ParseIP("192.168.0.1")},
	}

	f := n.GetModel()

	if !reflect.DeepEqual(f, m) {
		t.Errorf("Got another model on GetModel:\ngot  %v\nwant %v", f, m)
	}
}
