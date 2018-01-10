package message

import (
	"github.com/n0stack/n0core/model"
	"github.com/satori/go.uuid"
)

// Spec is sent to distributor to propagate Models.
//
// Args:
// 	spec_id: ID to distinguish spec as a user request.
// 	models: Models that the top of Model will be created.
// 	annotations: Options as scheduling hint and etc.
//
// Example:
// 	>>> from n0core.model import Model
// 	>>> m1 = Model(...)
// 	>>> m2 = Model(...)
// 	>>> Spec(spec_id="ba6f8ced-c8c2-41e9-98d0-5c961dff6c9cf",
// 			 models=[m1, m2])
type Spec struct {
	SpecID     uuid.UUID
	Models     []model.AbstractModel
	Annotaions map[string]string
}
