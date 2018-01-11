package resource

import (
	"net/url"

	"github.com/n0stack/n0core/model"
)

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
	model.Model

	Size uint64
	URL  *url.URL
}

func (v Volume) ToModel() *model.Model {
	return &v.Model
}
