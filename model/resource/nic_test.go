package resource

import (
	"net"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/n0stack/n0core/model"
	"github.com/satori/go.uuid"
)

func TestNICToModel(t *testing.T) {
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

	f := n.ToModel()

	if !reflect.DeepEqual(f, m) {
		t.Errorf("Got another model on ToModel:\ngot  %v\nwant %v", f, m)
	}
}

func TestNewNIC(t *testing.T) {
	id := uuid.NewV4()
	specificType := "hoge"
	m := &model.Model{
		ID:           id,
		Type:         filepath.Join(NICType, specificType),
		State:        "testing",
		Name:         "test_model",
		Meta:         map[string]string{"hoge": "hoge"},
		Dependencies: model.Dependencies{},
	}

	h, _ := net.ParseMAC("01:23:45:67:89:ab")

	v := &NIC{
		Model:   *m,
		HWAddr:  h,
		IPAddrs: []net.IP{net.ParseIP("192.168.0.1")},
	}

	nv := NewNIC(v.ID, specificType, v.State, v.Name, v.Meta, v.Dependencies, v.HWAddr, v.IPAddrs)

	if !reflect.DeepEqual(v, nv) {
		t.Errorf("Got another model on NewVM:\ngot  %v\nwant %v", v, nv)
	}
}
