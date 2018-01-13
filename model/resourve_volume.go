package model

import (
	"net/url"
	"path/filepath"

	uuid "github.com/satori/go.uuid"
)

const VolumeType = "resource/volume"

// Volume manage persistent volume resource.
//
// Example:
//
// 	.. code-block:: yaml
//
// 	type: resource/volume/file
// 	id: 486274b2-49e4-4bcd-a60d-4f627ce8c041
// 	state: allocated
// 	name: hogehoge
// 	size: 10 * 1024 * 1024 * 1024
// 	url: file:///data/hoge
//
// STATES:
// 	ALLOCATED: Allocate volume size and share volume.
// 	DEALLOCATED: Delete volume resource, but not delete data in volume.
// 	DELETED: Delete data in volume.
//
// Meta:
//
// Labels:
//
// Property:
//
// Args:
// 	id: UUID
// 	type:
// 	state:
// 	name: Name of volume.
// 	size: Size of volume.
// 	url: URL which is sharing like file:///data/hoge and nfs://hobge/data/hoge
// 	meta:
// 	dependencies: List of dependency to
type Volume struct {
	Model `yaml:",inline"`

	Size uint64
	URL  *url.URL
}

func (v Volume) ToModel() *Model {
	return &v.Model
}

func NewVolume(id uuid.UUID, specificType, state, name string, meta map[string]string, dependencies Dependencies, size uint64, u *url.URL) *Volume {
	return &Volume{
		Model: Model{
			ID:           id,
			Type:         filepath.Join(VolumeType, specificType),
			State:        state,
			Name:         name,
			Meta:         meta,
			Dependencies: Dependencies{},
		},
		Size: size,
		URL:  u,
	}
}
