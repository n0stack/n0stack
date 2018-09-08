package datastore

import (
	"testing"

	n0stack "github.com/n0stack/proto.go/v0"
)

func TestCheckVersion(t *testing.T) {
	Cases := []struct {
		name    string
		prev    *n0stack.Metadata
		new     *n0stack.Metadata
		version uint64
		err     string
	}{
		{
			name: "[Valid] prev:nil, new:0 -> 1",
			prev: nil,
			new: &n0stack.Metadata{
				Version: 0,
			},
			version: 1,
			err:     "",
		},
		{
			name: "[Valid] prev:1, new:1 -> 2",
			prev: &n0stack.Metadata{
				Version: 1,
			},
			new: &n0stack.Metadata{
				Version: 1,
			},
			version: 2,
			err:     "",
		},
		{
			name: "[Invalid] prev:nil, new:1 -> err",
			prev: nil,
			new: &n0stack.Metadata{
				Version: 1,
			},
			version: 0,
			err:     "Set 0 when create new object, have:1, want:0",
		},
		{
			name: "[Invalid] prev:1, new:2 -> err",
			prev: &n0stack.Metadata{
				Version: 1,
			},
			new: &n0stack.Metadata{
				Version: 2,
			},
			version: 0,
			err:     "Set the same version as stored in database, have:2, want:1",
		},
	}

	for _, c := range Cases {
		v, err := CheckVersion(c.prev, c.new)

		if v != c.version {
			t.Errorf("[%s] Wrong version\n\thave:%d\n\twant:%d", c.name, v, c.version)
		}

		if err == nil && c.err != "" {
			t.Errorf("[%s] Need error message\n\thave:nil\n\twant:%s", c.name, c.err)
		}

		if err != nil && err.Error() != c.err {
			t.Errorf("[%s] Wrong error message\n\thave:%s\n\twant:%s", c.name, err.Error(), c.err)
		}
	}
}
