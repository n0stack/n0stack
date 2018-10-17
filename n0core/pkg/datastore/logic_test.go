package datastore

import (
	"testing"
)

func TestCheckVersion(t *testing.T) {
	Cases := []struct {
		name    string
		prev    uint64
		new     uint64
		version uint64
		err     string
	}{
		{
			name:    "Valid: previous=0, new=0, version=1",
			prev:    0,
			new:     0,
			version: 1,
			err:     "",
		},
		{
			name:    "Valid: prev=1, new=1, version=2",
			prev:    1,
			new:     1,
			version: 2,
			err:     "",
		},
		{
			name:    "Invalid: previous=0, new=1, err",
			prev:    0,
			new:     1,
			version: 0,
			err:     "Set the same version as stored in database: have=1, want=0",
		},
		{
			name:    "Invalid: prev=1, new=2, err",
			prev:    1,
			new:     2,
			version: 0,
			err:     "Set the same version as stored in database: have=2, want=1",
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
