package target

import "github.com/n0stack/n0core/model"

// Target is application service to apply resources with some framework like KVM and iproute2.
//
// A target manage only one type `*/*/*` of resource like `resource/network/flat`.
// Directory structure and class name is ruled by resource type.
// For example, `resource/network/flat` define `class Flat` which is placed on `n0core.resource.network.flat`.
//
// Do not implement that resource is kill when target is killed.
//
// Args:
// 	support_model: Model type which is supported on each target.
//
// Example:
//
// 	in `n0core.target.vm.example`
// 	>>> class Exapmle(Target):
// 	>>>     def __init__(self):
// 	>>>         super().__init__("resource/vm/example")
type Target interface {

	// Apply resource with some framework.
	//
	// Args:
	// 	model: Model which you want to apply.
	//
	// Return:
	// 	Tuple of processed is_succeeded and description.
	Apply(m *model.AbstractModel) (bool, string)
}
