package datastore

import (
	"testing"

	"github.com/n0stack/proto.go/v0"

	n0stack "github.com/n0stack/proto.go/v0"
)

type stubHavingMetadata struct {
	metadata *pn0stack.Metadata
}

func (s stubHavingMetadata) GetMetadata() *pn0stack.Metadata {
	return s.metadata
}

func TestCheckVersion(t *testing.T) {
	Cases := []struct {
		name    string
		prev    *stubHavingMetadata
		new     *stubHavingMetadata
		version uint64
		err     string
	}{
		{
			name: "Valid: previous=nil, new=0, version=1",
			prev: nil,
			new: &stubHavingMetadata{
				metadata: &n0stack.Metadata{
					Version: 0,
				},
			},
			version: 1,
			err:     "",
		},
		{
			name: "Valid: prev=1, new=1, version=2",
			prev: &stubHavingMetadata{
				metadata: &n0stack.Metadata{
					Version: 1,
				},
			},
			new: &stubHavingMetadata{
				metadata: &n0stack.Metadata{
					Version: 1,
				},
			},
			version: 2,
			err:     "",
		},
		{
			name: "Invalid: previous=nil, new=1, err",
			prev: nil,
			new: &stubHavingMetadata{
				metadata: &n0stack.Metadata{
					Version: 1,
				},
			},
			version: 0,
			err:     "Set 0 when create new object: have=1, want=0",
		},
		{
			name: "Invalid: prev=1, new=2, err",
			prev: &stubHavingMetadata{
				metadata: &n0stack.Metadata{
					Version: 1,
				},
			},
			new: &stubHavingMetadata{
				metadata: &n0stack.Metadata{
					Version: 2,
				},
			},
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
